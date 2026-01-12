package storage

import "TelegramShop/models"

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

func AddCategory(catname string) error {
	user := models.Category{
		Name: catname,
	}

	return DB.Where(models.Category{Name: catname}).FirstOrCreate(&user).Error
}

func DelCategory(catid int) error {
	return DB.Delete(&models.Category{}, catid).Error
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

func GetCategory(catid int) (models.Category, error) {
	var category models.Category

	err := DB.First(&category, catid).Error
	if err != nil {
		return models.Category{}, err
	}

	return category, nil
}
