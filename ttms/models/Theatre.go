package models

import (
	utils "TTMS_go/ttms/util"
	"encoding/json"
	"gorm.io/gorm"
)

type Theatre struct {
	gorm.Model
	Name  string
	Seat  []byte
	N     int
	M     int
	Info  string
	Inuse int //影院状态 0 1 2
	Num   int
	Plays []Play `gorm:"type:json"`
}

func (theatre Theatre) TableName() string {
	return "theatre_basic"
}

func CreateTheatre(theatre Theatre) {
	utils.DB.Create(&theatre)
}

func FindTheatreByid(id string) Theatre {
	theatre := Theatre{}
	utils.DB.Where("id = ?", id).First(&theatre)
	var playJson string
	utils.DB.Table("theatre_basic").Where("id = ?", id).Select("plays").Scan(&playJson)
	json.Unmarshal([]byte(playJson), &theatre.Plays)
	return theatre
}

func UpdateTheatre(theatre *Theatre) {
	//序列化：
	playJson, _ := json.Marshal(theatre.Plays)
	utils.DB.Model(theatre).Updates(map[string]interface{}{
		"plays": string(playJson),
		"info":  theatre.Info,
	})
}
