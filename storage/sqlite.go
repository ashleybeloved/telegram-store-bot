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

	err = DB.AutoMigrate(&models.User{}, &models.Category{}, &models.Product{}, &models.Item{}, &models.PurchasesHistory{}, &models.Promocode{}, &models.PromocodeUsage{})
	if err != nil {
		return fmt.Errorf("ошибка миграции: %v", err)
	}

	log.Println("Подключение к SQLite через GORM успешно!")
	return nil
}
