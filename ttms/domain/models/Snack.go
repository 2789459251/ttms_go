package models

import (
	utils "TTMS_go/ttms/util"
	"gorm.io/gorm"
)

type Snack struct {
	gorm.Model
	Name    string
	Picture string
	Info    string
	Stock   int     //库存量
	Price   float64 //价格
}

func (snack Snack) TableName() string {
	return "snack_basic"
}
func Showsnacks() (snacks []Snack) {
	utils.DB.Find(snacks)
	return
}
func SearchSnack(name string) (snacks []Snack) {
	str := ""
	for i, i2 := range name {
		c := string(i2)
		if i == 0 {
			str += "NameLIKE %" + c + "%"
		} else {
			str += "AND Name LIKE %" + c + "%"
		}
	}
	utils.DB.Where(str).Find(snacks)
	return
}
func Insertsnack(snack Snack) {
	utils.DB.Create(snack)
}
