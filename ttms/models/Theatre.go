package models

import (
	utils "TTMS_go/ttms/util"
	"errors"
	"gorm.io/gorm"
)

type Theatre struct {
	gorm.Model
	Name  string
	Seat  string
	N     int
	M     int
	Info  string
	Inuse int //影院状态 0 1 2
	Num   int
	Plays string
}

func (theatre Theatre) TableName() string {
	return "theatre_basic"
}

func CreateTheatre(theatre Theatre) {
	utils.DB.Create(&theatre)
}

func FindAllTheatre() []Theatre {
	var theatre []Theatre
	utils.DB.Find(&theatre)
	return theatre
}
func FindTheatreByid(id string) Theatre {
	theatre := Theatre{}
	utils.DB.Where("id = ?", id).First(&theatre)
	//utils.DB.Raw("SELECT * FROM theatre_basic WHERE id = ?", id).Scan(&theatre)
	return theatre
}

func UpdateTheatre(theatre *Theatre) {
	utils.DB.Model(theatre).Save(&theatre)
}
func DeleteTheatre(id string) error {
	t := Theatre{}
	utils.DB.Where("id = ?", id).Delete(&t)
	utils.DB.Exec("DELETE FROM theatre_basic WHERE id = ?", id)
	if t.Name == "" {
		return errors.New("没有id为" + id + "的影院")
	}
	return nil
}
