package storage

import (
	"TelegramShop/models"

	"gorm.io/gorm"
)

type ItemsBrief struct {
	ID   int
	Data string
}

func GetItems(page int, productid int) ([]ItemsBrief, error) {
	var results []ItemsBrief
	pageSize := 5
	offset := (page - 1) * pageSize

	err := DB.Model(&models.Item{}).
		Select("id", "data").
		Where("product_id = ?", productid).
		Limit(pageSize).
		Offset(offset).
		Order("id ASC").
		Scan(&results).Error

	return results, err
}

func GetItem(itemid int) (models.Item, error) {
	var item models.Item
	err := DB.First(&item, itemid).Error

	return item, err
}

func GetPagesForItems(productid int) (int, error) {
	var count int64

	err := DB.Model(&models.Item{}).
		Where("product_id = ?", productid).
		Count(&count).Error

	if err != nil {
		return 0, err
	}

	return int((count + 5 - 1) / 5), nil
}

func AddItem(productID int, itemData string) error {
	item := models.Item{
		ProductID: uint(productID),
		Data:      itemData,
		IsSold:    false,
	}

	err := DB.Model(&models.Product{}).Where("id = ?", productID).Update("Stock", gorm.Expr("Stock + ?", 1)).Error
	if err != nil {
		return err
	}

	return DB.Create(&item).Error
}

func DelItem(itemid int) error {
	item, err := GetItem(itemid)
	if err != nil {
		return err
	}

	err = DB.Model(&models.Product{}).Where("id = ?", item.ProductID).Update("Stock", gorm.Expr("Stock - ?", 1)).Error
	if err != nil {
		return err
	}

	return DB.Delete(&models.Item{}, itemid).Error
}
