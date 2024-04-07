package models

import (
	utils "TTMS_go/ttms/util"
	"gorm.io/gorm"
)

type UserInfo struct {
	gorm.Model
	Wallet float64
	//Ticket []uint
	Ticket   []Ticket `gorm:"type:json"`
	Snack    []Snack_ `gorm:"type:json"`
	Favorite []int    `gorm:"type:json"`
}

func (user UserInfo) TableName() string {
	return "user_info"
}
func FindUserInfo(id string) UserInfo {
	u := &UserInfo{}
	utils.DB.Where("id = ?", id).First(&u)
	return *u
}
func (u UserInfo) RefleshUserInfo() (err error) {
	err = utils.DB.Updates(u).Error
	return
}
