package models

import (
	utils "TTMS_go/ttms/util"
	"gorm.io/gorm"
	"sync"
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
	mu          sync.RWMutex
}

func (movie Movie) TableName() string {
	return "movie_basic"
}

func Update(m Movie) {
	m.mu.Lock()
	defer m.mu.Unlock()
	utils.DB.Where("name = ?", m.Name).Find(&m)
	utils.DB.Save(&m)
}

func MovieList() []Movie {
	m := []Movie{}
	utils.DB.Find(m)
	return m
}

func UpcommingList() []Movie {
	m := []Movie{}
	utils.DB.Order("release_time ASC").Where("release_time  > ?", time.Now()).Find(&m)
	return m
}

func HitList() []Movie {
	m := []Movie{}
	utils.DB.Order("score ASC").Where("release_time < ?", time.Now()).Find(&m)
	return m
}

func DeleteMovieById(ids []string) {
	utils.DB.Where("id in (?)", ids).Delete(&Movie{})
}
