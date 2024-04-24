package models

import (
	utils "TTMS_go/ttms/util"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type UserInfo struct {
	gorm.Model
	Wallet float64
	//Ticket []uint
	Flag          int       `gorm:"default:0"` // 0 用户 1 管理员
	Ticket        []Ticket  `gorm:"type:json"`
	Snack         []Snack   `gorm:"type:json"`
	FavoriteMovie []int     `gorm:"type:json"`
	FavoriteSnack []int     `gorm:"type:json"`
	Name          string    //昵称
	ProfilePhoto  string    //头像
	Birthday      time.Time `gorm:"type:date;DEFAULT:NULL"` // 生日
	Interest      []string  `gorm:"type:json"`              //兴趣
	Sign          string    //签名
}

func (user UserInfo) TableName() string {
	return "user_info"
}

func FindUserInfo(id string) UserInfo {
	u := UserInfo{}
	utils.DB.Where("id = ? ", id).First(&u)
	var interestJson, ticketJson, sign, snackJson, favoriteMovieJson, favoriteSnackJson string
	utils.DB.Table("user_info").Where("id = ? ", id).Select("interest").Scan(&interestJson)
	utils.DB.Table("user_info").Where("id = ? ", id).Select("ticket").Scan(&ticketJson)
	utils.DB.Table("user_info").Where("id = ? ", id).Select("sign").Scan(&sign)
	utils.DB.Table("user_info").Where("id = ? ", id).Select("snack").Scan(&snackJson)
	utils.DB.Table("user_info").Where("id = ? ", id).Select("favorite_movie").Scan(&favoriteMovieJson)
	utils.DB.Table("user_info").Where("id = ? ", id).Select("favorite_snack").Scan(&favoriteSnackJson)
	json.Unmarshal([]byte(interestJson), &u.Interest)
	json.Unmarshal([]byte(ticketJson), &u.Ticket)
	json.Unmarshal([]byte(snackJson), &u.Snack)
	json.Unmarshal([]byte(favoriteMovieJson), &u.FavoriteMovie)
	json.Unmarshal([]byte(favoriteSnackJson), &u.FavoriteSnack)
	u.Sign = sign
	fmt.Println(u)
	return u
}
func (u UserInfo) RefleshUserInfo_() (err error) {
	interestJson, _ := json.Marshal(u.Interest)
	ticketJson, _ := json.Marshal(u.Ticket)
	snackJson, _ := json.Marshal(u.Snack)
	favoriteMovieJson, _ := json.Marshal(u.FavoriteMovie)
	favoriteSnackJson, _ := json.Marshal(u.FavoriteSnack)

	err = utils.DB.Model(u).Updates(map[string]interface{}{
		"interest":       interestJson,
		"snack":          snackJson,
		"ticket":         ticketJson,
		"favorite_movie": favoriteMovieJson,
		"favorite_snack": favoriteSnackJson,
		"flag":           u.Flag,
		"name":           u.Name,         //昵称
		"profile_photo":  u.ProfilePhoto, //头像
		//"birthday":       u.Birthday,     // 生日         //兴趣
		"sign":   u.Sign, //签名
		"wallet": u.Wallet,
	}).Error
	//if time.Time.Unix(u.Birthday) != 0 {
	//	utils.DB.Model(u).Update("birthday", u.Birthday)
	//}
	return
}
func (u UserInfo) tx_RefleshUserInfo(tx *gorm.DB) (err error) {
	interestJson, _ := json.Marshal(u.Interest)
	ticketJson, _ := json.Marshal(u.Ticket)
	snackJson, _ := json.Marshal(u.Snack)
	favoriteMovieJson, _ := json.Marshal(u.FavoriteMovie)
	favoriteSnackJson, _ := json.Marshal(u.FavoriteSnack)

	err = tx.Model(u).Updates(map[string]interface{}{
		"interest":       interestJson,
		"snack":          snackJson,
		"ticket":         ticketJson,
		"favorite_movie": favoriteMovieJson,
		"favorite_snack": favoriteSnackJson,
		"flag":           u.Flag,
		"name":           u.Name,         //昵称
		"profile_photo":  u.ProfilePhoto, //头像
		//"birthday":       u.Birthday,     // 生日         //兴趣
		"sign":   u.Sign, //签名
		"wallet": u.Wallet,
	}).Error
	//if time.Time.Unix(u.Birthday) != 0 {
	//	utils.DB.Model(u).Update("birthday", u.Birthday)
	//}
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
