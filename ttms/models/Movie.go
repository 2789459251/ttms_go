package models

import (
	utils "TTMS_go/ttms/util"
	"gorm.io/gorm"
	"time"
)

type Movie struct {
	gorm.Model
	Info        string
	Name        string
	Director    string
	Actor       string
	Score       float64
	Duration    time.Duration
	ReleaseTime time.Time
	Money       float64
}

func (movie Movie) TableName() string {
	return "movie_basic"
}

func Update(m Movie) {
	utils.DB.Where("name = ?", m.Name).Find(&m)
	utils.DB.Save(&m)
}
