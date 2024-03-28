package models

import (
	"gorm.io/gorm"
	"time"
)

type Place struct {
	gorm.Model
	Name      string
	Seat      [][]bool `gorm:"type:json"`
	N         int
	M         int
	Movie     int
	Inuse     bool
	Num       int //已定座位数量
	Begintime time.Time
	Aftertime time.Time
}

func (p *Place) Init() {
	p.Seat = make([][]bool, p.N)
	for i := 0; i < p.N; i++ {
		p.Seat[i] = make([]bool, p.M)
	}
}
func (place Place) TableName() string {
	return "place_basic"
}
