package models

import (
	utils "TTMS_go/ttms/util"
	"encoding/json"
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
	//ticketJson, _ := json.Marshal(u.Ticket)
	//snackJson, _ := json.Marshal(u.Snack)
	//favoriteMovieJson, _ := json.Marshal(u.FavoriteMovie)
	//favoriteSnackJson, _ := json.Marshal(u.FavoriteSnack)
	var interestJson, ticketJson, sign, snackJson, favoriteMovieJson, favoriteSnackJson string
	//"sign","snack","favorite_movie","favorite_snack").Rows()
	utils.DB.Table("user_info").Where("id = ? ", id).Select("interest").Scan(&interestJson)
	utils.DB.Table("user_info").Where("id = ? ", id).Select("ticket").Scan(&ticketJson)
	utils.DB.Table("user_info").Where("id = ? ", id).Select("sign").Scan(&sign)
	utils.DB.Table("user_info").Where("id = ? ", id).Select("snack").Scan(&snackJson)
	utils.DB.Table("user_info").Where("id = ? ", id).Select("favorite_movie").Scan(&favoriteMovieJson)
	utils.DB.Table("user_info").Where("id = ? ", id).Select("favorite_snack").Scan(&favoriteSnackJson)

	interest := []string{}
	tickets := []Ticket{}
	Snacks := []Snack{}
	favoriteMovies := []int{}
	favoriteSnacks := []int{}

	json.Unmarshal([]byte(interestJson), &interest)
	json.Unmarshal([]byte(ticketJson), &tickets)
	json.Unmarshal([]byte(snackJson), &Snacks)
	json.Unmarshal([]byte(favoriteMovieJson), &favoriteMovies)
	json.Unmarshal([]byte(favoriteSnackJson), &favoriteSnacks)
	u.Interest = interest
	u.Ticket = tickets
	u.Snack = Snacks
	u.FavoriteSnack = favoriteMovies
	u.FavoriteMovie = favoriteMovies
	u.Sign = sign
	return u
}
func (u UserInfo) RefleshUserInfo_() (err error) {
	//uu := []UserInfo{}
	//uu = append(uu, u)
	//fmt.Println(u)
	//MySQL支持JSON数据类型，您可以将数组转换成JSON格式，
	//并将其存储到JSON类型的字段中。在MySQL中，
	//您可以使用JSON_ARRAY函数创建一个JSON数组。
	//INSERT INTO table_name (json_array_column) VALUES (JSON_ARRAY(1, 2, 3, 4, 5));
	//Ticket        []Ticket  `gorm:"type:json"`
	//Snack         []Snack_  `gorm:"type:json"`
	//FavoriteMovie []int     `gorm:"type:json"`
	//FavoriteSnack []int     `gorm:"type:json"`
	interestJson, _ := json.Marshal(u.Interest)
	ticketJson, _ := json.Marshal(u.Ticket)
	snackJson, _ := json.Marshal(u.Snack)
	favoriteMovieJson, _ := json.Marshal(u.FavoriteMovie)
	favoriteSnackJson, _ := json.Marshal(u.FavoriteSnack)
	err = utils.DB.Save(u).Error
	utils.DB.Updates(map[string]interface{}{
		"interest":      interestJson,
		"snack":         snackJson,
		"ticket":        ticketJson,
		"favoriteMovie": favoriteMovieJson,
		"favoriteSnack": favoriteSnackJson,
	})

	return
}
func (u UserInfo) FindUserinfoByid(id string) (user UserInfo, err error) {
	err = utils.DB.Where("id = ?", id).First(&user).Error
	return
}
