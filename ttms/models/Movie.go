package models

import (
	utils "TTMS_go/ttms/util"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"sync"
	"time"
)

// todo 在电影类可以加预告片的放映。
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

func FindMovieByid(id string) Movie {
	m := Movie{}
	utils.DB.Where("id = ?", id).Find(&m)
	return m
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

func FavoriteRankingMovies(members []redis.Z) []byte {
	str := []byte{}
	type movieInfo struct {
		m     Movie
		score float64
	}
	tmp := &movieInfo{}
	for _, member := range members {
		utils.DB.Where("id = ?", member.Member).First(&tmp.m)
		tmp.score = member.Score
		t, _ := json.Marshal(tmp)
		str = append(str, t...)
	}
	return str
}
