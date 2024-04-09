package models

import (
	utils "TTMS_go/ttms/util"
	"gorm.io/gorm"
)

type UserInfo struct {
	gorm.Model
	Wallet float64
	//Ticket []uint
	Flag     int      `gorm:"default:0"` // 0 用户 1 管理员
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
func (u UserInfo) RefleshUserInfo_() (err error) {
	uu := []UserInfo{}
	uu = append(uu, u)
	err = utils.DB.Where("id = ?", u.ID).Updates(uu).Error
	return
}
func (u UserInfo) FindUserinfoByid(id string) (user UserInfo, err error) {
	err = utils.DB.Where("id = ?", id).First(&user).Error
	return
}
