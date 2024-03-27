package models

import "gorm.io/gorm"

type Movie struct {
	gorm.Model
	Info  string
	Name  string
	Money float64
	place Place //申请演出厅
}

func (movie Movie) TableName() string {
	return "movie_basic"
}
