package storage

import (
	"TelegramShop/models"
	"fmt"
	"time"

	"gorm.io/gorm"
)

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
