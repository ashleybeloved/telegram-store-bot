package storage

import (
	"TelegramShop/models"
	"fmt"
	"log"
	"os"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func OpenSQLite() error {
	if err := os.MkdirAll("./data", 0755); err != nil {
		return err
	}

	var err error
	DB, err = gorm.Open(sqlite.Open("./data/database.db"), &gorm.Config{})
	if err != nil {
		return err
	}

	err = DB.AutoMigrate(&models.User{}, &models.Category{}, &models.Product{})
	if err != nil {
		return fmt.Errorf("ошибка миграции: %v", err)
	}

	log.Println("Подключение к SQLite через GORM успешно!")
	return nil
}

func AddUser(userID int64, username, firstname, lastname, langCode string) error {
	user := models.User{
		UserID:    userID,
		Username:  username,
		Firstname: firstname,
		Lastname:  lastname,
		LangCode:  langCode,
	}

	return DB.Where(models.User{UserID: userID}).FirstOrCreate(&user).Error
}

func FindUser(userid int64) (*models.User, error) {
	var user models.User
	err := DB.Where("user_id = ?", userid).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func RefreshUser(userid int64, username, firstname, lastname string, lang_code string) (*models.User, error) {
	var user models.User

	err := DB.Model(&user).Where("user_id = ?", userid).Updates(models.User{
		Username:  username,
		Firstname: firstname,
		Lastname:  lastname,
		LangCode:  lang_code,
	}).Error

	if err != nil {
		return nil, err
	}

	return FindUser(userid)
}

func GetPagesForCategories() (int, error) {
	var count int64
	err := DB.Model(&models.Category{}).Count(&count).Error
	if err != nil {
		return 0, err
	}

	return int((count + 5 - 1) / 5), nil
}

type CategoryBrief struct {
	ID   int
	Name string
}

func GetCategories(page int) ([]CategoryBrief, error) {
	var results []CategoryBrief
	pageSize := 5
	offset := (page - 1) * pageSize

	err := DB.Model(&models.Category{}).
		Select("id", "name").
		Limit(pageSize).
		Offset(offset).
		Order("id ASC").
		Scan(&results).Error

	return results, err
}

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
			return fmt.Errorf("недостаточно средств (нужно %d, у вас %d)", product.Price, user.Balance)
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

		err := AddPurchaseToPurchaseHistory(tx, userid, product, item.Data, item.ID)
		if err != nil {
			return err
		}

		return nil
	})

	return soldData, err
}

func AddPurchaseToPurchaseHistory(tx *gorm.DB, userid int64, product models.Product, itemData string, itemid uint) error {
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
