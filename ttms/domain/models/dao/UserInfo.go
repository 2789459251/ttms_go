package dao

import (
	"TTMS_go/ttms/domain/models"
	utils "TTMS_go/ttms/util"
	"gorm.io/gorm"
)

type UserInfo struct {
	gorm.Model
	Wallet int
	Ticket []models.Ticket
	Snack  []models.Snack
}

func (user UserInfo) TableName() string {
	return "user_info"
}
func FindUserInfo(id string) UserInfo {
	u := &UserInfo{}
	utils.DB.Where("id = ?", id).First(&u)
	return *u
}
