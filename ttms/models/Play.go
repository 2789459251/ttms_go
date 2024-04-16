package models

import (
	utils "TTMS_go/ttms/util"
	"gorm.io/gorm"
	"time"
)

type Play struct {
	gorm.Model
	MovieId   string
	TheatreId string
	Seat      [][]int //0 1 2
	Num       string  //剩余座位数量
	BeginTime time.Time
	EndTime   time.Time
}

func CreatePlay(play *Play) {
	utils.DB.Create(play)
}
func ShowPlaysByMovieId(id string) []Play {
	plays := []Play{}
	utils.DB.Where("movie_id = ?", id).Find(plays)
	return plays
}
func ShowPlaysByTheatreId(id string) []Play {
	plays := []Play{}
	utils.DB.Where("theatre_id = ?", id).Find(plays)
	return plays
}
func ShowPlayById(id string) *Play {
	p := &Play{}
	utils.DB.Where("id = ?", id).Find(p)
	return p
}
