package models

import (
	utils "TTMS_go/ttms/util"
	"gorm.io/gorm"
)

type Theatre struct {
	gorm.Model
	Name  string
	Seat  [][]bool `gorm:"type:json"`
	N     int
	M     int
	Info  string
	Inuse int //影院状态 0 1 2
	Num   int
	Plays *Node `gorm:"type:json"`
}
type Node struct {
	Play Play
	Next *Node
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
	return theatre
}

func UpdateTheatre(theatre *Theatre) {
	utils.DB.Where("id = ?", theatre.ID).Updates(theatre)
}
