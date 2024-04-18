package models

import (
	utils "TTMS_go/ttms/util"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"gorm.io/gorm"
	"time"
)

type UserInfo struct {
	gorm.Model
	Wallet float64
	//Ticket []uint
	Flag          int       `gorm:"default:0"` // 0 用户 1 管理员
	Ticket        []Ticket  `gorm:"type:json"`
	Snack         []Snack_  `gorm:"type:json"`
	FavoriteMovie []int     `gorm:"type:json"`
	FavoriteSnack []int     `gorm:"type:json"`
	Name          string    //昵称
	ProfilePhoto  string    //头像
	Birthday      time.Time `gorm:"type:date;DEFAULT:NULL"` // 生日
	Interest      []string  `gorm:"type:json"`              //兴趣
	Sign          string    //签名
}

type stringSlice []string

func (user UserInfo) TableName() string {
	return "user_info"
}
func (ri stringSlice) Value() (driver.Value, error) {
	// 将 ResultInfo 结构体切片转换为 JSON 格式的字符串
	value, err := json.Marshal(ri)
	if err != nil {
		return nil, err
	}
	return string(value), nil
}

// Scan 将数据库中的值解析为 ResultInfoSlice 结构体切片
func (ri stringSlice) Scan(value interface{}) error {
	// 将数据库中的值解析为字符串
	stringValue, ok := value.(string)
	if !ok {
		return errors.New("不是 ResultInfo 切片类型")
	}

	// 将 JSON 格式的字符串解析为 ResultInfo 结构体切片
	var resultInfoSlice stringSlice
	if err := json.Unmarshal([]byte(stringValue), &resultInfoSlice); err != nil {
		return err
	}

	ri = resultInfoSlice

	return nil
}

func FindUserInfo(id string) UserInfo {
	u := &UserInfo{}
	utils.DB.Where("id = ?", id).First(&u)
	return *u
}
func (u UserInfo) RefleshUserInfo_() (err error) {
	//uu := []UserInfo{}
	//uu = append(uu, u)
	//fmt.Println(u)
	//MySQL支持JSON数据类型，您可以将数组转换成JSON格式，
	//并将其存储到JSON类型的字段中。在MySQL中，
	//您可以使用JSON_ARRAY函数创建一个JSON数组。
	//INSERT INTO table_name (json_array_column) VALUES (JSON_ARRAY(1, 2, 3, 4, 5));

	err = utils.DB.Save(u).Error
	return
}
func (u UserInfo) FindUserinfoByid(id string) (user UserInfo, err error) {
	err = utils.DB.Where("id = ?", id).First(&user).Error
	return
}
