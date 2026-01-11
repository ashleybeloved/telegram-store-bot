package storage

import (
	"TelegramShop/models"
	"fmt"
	"log"
	"os"
	"time"

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

func GetUser(userid int64) (*models.User, error) {
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

	return GetUser(userid)
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

func GetPurchaseHistory(userid int64) ([]models.PurchasesHistory, error) {
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

func SetUserState(userid int64, state string) error {
	return DB.Model(&models.User{}).Where("user_id = ?", userid).Update("state", state).Error
}

func AddBalance(userid int64, amount int64) error {
	return DB.Model(&models.User{}).Where("user_id = ?", userid).Update("balance", gorm.Expr("balance + ?", amount)).Error
}

func GetPromocodes(page int) ([]models.Promocode, error) {
	var promocodes []models.Promocode
	pageSize := 5
	offset := (page - 1) * pageSize

	err := DB.Model(&models.Promocode{}).
		Limit(pageSize).
		Offset(offset).
		Order("id ASC").
		Find(&promocodes).Error

	return promocodes, err
}

func NewPromocode(code string, reward int64, maxUses int, expiresAt time.Time) error {
	promocode := models.Promocode{
		Code:      code,
		Reward:    reward,
		MaxUses:   maxUses,
		UsesLeft:  maxUses,
		ExpiresAt: expiresAt,
	}

	return DB.Where(models.Promocode{Code: code}).FirstOrCreate(&promocode).Error
}

func RedeemPromocode(userid int64, code string) (int64, error) {
	var promocode models.Promocode

	err := DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("code = ?", code).First(&promocode).Error; err != nil {
			return fmt.Errorf("промокод не найден")
		}

		if promocode.UsesLeft <= 0 {
			return fmt.Errorf("промокод больше не действует (использован максимально допустимое количество раз)")
		}

		if time.Now().After(promocode.ExpiresAt) {
			return fmt.Errorf("промокод истёк")
		}

		if err := tx.Model(&models.PromocodeUsage{}).Where("user_id = ? AND promocode_id = ?", userid, promocode.ID).First(&models.PromocodeUsage{}).Error; err == nil {
			return fmt.Errorf("вы уже активировали этот промокод ранее")
		}

		if err := tx.Model(&models.User{}).Where("user_id = ?", userid).Update("balance", gorm.Expr("balance + ?", promocode.Reward)).Error; err != nil {
			return err
		}

		usage := models.PromocodeUsage{
			UserID:      userid,
			PromocodeID: promocode.ID,
		}

		if err := tx.Create(&usage).Error; err != nil {
			return err
		}

		promocode.UsesLeft -= 1

		if err := tx.Save(&promocode).Error; err != nil {
			return err
		}

		return nil
	})

	return promocode.Reward, err
}

func GetPagesForPromocodes() (int, error) {
	var count int64

	err := DB.Model(&models.Promocode{}).
		Count(&count).Error

	if err != nil {
		return 0, err
	}

	return int((count + 5 - 1) / 5), nil
}

func GetPromocode(promocode_id int) (models.Promocode, error) {
	var promocode models.Promocode
	err := DB.First(&promocode, promocode_id).Error
	if err != nil {
		return models.Promocode{}, err
	}

	return promocode, nil
}

func DeletePromocode(promocode_id int) error {
	return DB.Delete(&models.Promocode{}, promocode_id).Error
}
