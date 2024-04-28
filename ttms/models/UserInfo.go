package models

import (
	utils "TTMS_go/ttms/util"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type UserInfo struct {
	gorm.Model
	Wallet        float64
	Flag          int `gorm:"default:0"` // 0 用户 1 管理员
	Ticket        string
	Snack         string
	FavoriteMovie string
	FavoriteSnack string
	Name          string    //昵称
	ProfilePhoto  string    //头像
	Birthday      time.Time `gorm:"type:date;DEFAULT:NULL"` // 生日
	Interest      string    //兴趣
	Sign          string    //签名
}

func (user UserInfo) TableName() string {
	return "user_info"
}

func FindUserInfo(id string) UserInfo {
	u := UserInfo{}
	utils.DB.Where("id = ? ", id).First(&u)
	fmt.Println(u)
	return u
}
func (u UserInfo) RefleshUserInfo_() {
	utils.DB.Save(&u)
	return
}
func (u UserInfo) Tx_RefleshUserInfo(tx *gorm.DB) (err error) {

	err = tx.Model(u).Updates(map[string]interface{}{
		"interest":       u.Interest,
		"snack":          u.Snack,
		"ticket":         u.Ticket,
		"favorite_movie": u.FavoriteMovie,
		"favorite_snack": u.FavoriteSnack,
		"flag":           u.Flag,
		"name":           u.Name,         //昵称
		"profile_photo":  u.ProfilePhoto, //头像
		//"birthday":       u.Birthday,     // 生日         //兴趣
		"sign":   u.Sign, //签名
		"wallet": u.Wallet,
	}).Error
	if time.Time.Unix(u.Birthday) != 0 {
		utils.DB.Model(u).Update("birthday", u.Birthday)
	}
	return
}
func (u UserInfo) FindUserinfoByid(id string) (user UserInfo, err error) {
	err = utils.DB.Where("id = ?", id).First(&user).Error
	return
}
func Recharge(num float64, user UserInfo) UserInfo {
	utils.DB.Model(&user).Update("wallet", num)
	return user
}
