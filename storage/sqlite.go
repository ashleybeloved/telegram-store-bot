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
