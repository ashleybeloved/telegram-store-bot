package storage

import (
	"TelegramShop/models"

	"gorm.io/gorm"
)

func SetUserState(userid int64, state string) error {
	return DB.Model(&models.User{}).Where("user_id = ?", userid).Update("state", state).Error
}

func AddBalance(userid int64, amount int64) error {
	return DB.Model(&models.User{}).Where("user_id = ?", userid).Update("balance", gorm.Expr("balance + ?", amount)).Error
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
