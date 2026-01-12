package storage

import (
	"TelegramShop/models"
	"fmt"

	"gorm.io/gorm"
)

func GetPagesForProducts(catid int) (int, error) {
	var count int64

	err := DB.Model(&models.Product{}).
		Where("category_id = ?", catid).
		Count(&count).Error

	if err != nil {
		return 0, err
	}

	return int((count + 5 - 1) / 5), nil
}

type ProductsBrief struct {
	ID   int
	Name string
}

func AddProduct(categoryID int, name string, description string, price int64) error {
	product := models.Product{
		CategoryID:  uint(categoryID),
		Name:        name,
		Description: description,
		Price:       price,
	}

	return DB.Where(models.Product{CategoryID: uint(categoryID), Name: name}).FirstOrCreate(&product).Error
}

func DelProduct(productid int) error {
	return DB.Delete(&models.Product{}, productid).Error
}

func GetProducts(page int, cat_id int) ([]ProductsBrief, error) {
	var results []ProductsBrief
	pageSize := 5
	offset := (page - 1) * pageSize

	err := DB.Model(&models.Product{}).
		Select("id", "name").
		Where("category_id = ?", cat_id).
		Limit(pageSize).
		Offset(offset).
		Order("id ASC").
		Scan(&results).Error

	return results, err
}

func GetProduct(product_id int) (models.Product, error) {
	var product models.Product
	err := DB.First(&product, product_id).Error
	if err != nil {
		return models.Product{}, err
	}

	return product, nil
}

func BuyProduct(userid int64, productid uint) (string, error) {
	var soldData string

	err := DB.Transaction(func(tx *gorm.DB) error {
		var user models.User
		var product models.Product
		var item models.Item

		if err := tx.Where("user_id = ?", userid).First(&user).Error; err != nil {
			return fmt.Errorf("пользователь не найден")
		}

		if err := tx.First(&product, productid).Error; err != nil {
			return fmt.Errorf("товар не найден")
		}

		if user.Balance < product.Price {
			return fmt.Errorf("недостаточно средств")
		}

		if err := tx.Where("product_id = ? AND is_sold = ?", productid, false).First(&item).Error; err != nil {
			return fmt.Errorf("товара нет в наличии (все продано)")
		}

		user.Balance -= product.Price
		if err := tx.Save(&user).Error; err != nil {
			return err
		}

		item.IsSold = true
		if err := tx.Save(&item).Error; err != nil {
			return err
		}

		soldData = item.Data

		err := AddPurchaseToPurchasesHistory(tx, userid, product, item.Data, item.ID)
		if err != nil {
			return err
		}

		if err := tx.Model(&models.Product{}).Where("id = ?", productid).Update("Stock", gorm.Expr("Stock - ?", 1)).Error; err != nil {
			return err
		}

		return nil
	})

	return soldData, err
}

func AddPurchaseToPurchasesHistory(tx *gorm.DB, userid int64, product models.Product, itemData string, itemid uint) error {
	purchase := models.PurchasesHistory{
		UserID:      userid,
		ItemID:      itemid,
		ProductID:   product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Data:        itemData,
	}

	return tx.Create(&purchase).Error
}

func GetPurchase(purchase_id int) (models.PurchasesHistory, error) {
	var purchase models.PurchasesHistory
	err := DB.First(&purchase, purchase_id).Error
	if err != nil {
		return purchase, err
	}

	return purchase, nil
}

func GetPurchasesHistory(userid int64) ([]models.PurchasesHistory, error) {
	var history []models.PurchasesHistory

	err := DB.Where("user_id = ?", userid).Order("created_at DESC").Find(&history).Error
	if err != nil {
		return nil, err
	}

	return history, nil
}

func GetPagesForPurchasesHistory(userid int64) (int, error) {
	var count int64

	err := DB.Model(&models.PurchasesHistory{}).
		Where("user_id = ?", userid).
		Count(&count).Error

	if err != nil {
		return 0, err
	}

	return int((count + 5 - 1) / 5), nil
}
