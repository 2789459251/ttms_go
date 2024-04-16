package models

import (
	utils "TTMS_go/ttms/util"
	"context"
	"encoding/json"
	"fmt"
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
	Duration    time.Duration
	ReleaseTime time.Time
	Money       float64
	Total       int     `json:"total"`   // 电影的总分
	Count       int     `json:"count"`   // 评分人数
	Average     float64 `json:"average"` // 平均分
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

func RankingMovies(members []redis.Z) []byte {
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

func UpdateMovieMark(m Movie, IMDbScore int, key string, movieId string) Movie {
	m.mu.Lock()
	defer m.mu.Unlock()
	//todo 事务
	if err := utils.Red.ZScore(context.Background(), key, movieId).Err(); err != nil {
		//没评价过
		if err == redis.Nil {
			m.Total += IMDbScore
			m.Count++
			m.Average = float64(m.Total) / float64(m.Count)
		} else {
			fmt.Errorf("查询用户评分出现错误，err:%v", err.Error())
		}

	} else {
		//评价过
		score, _ := utils.Red.ZScore(context.Background(), key, movieId).Result()
		m.Total = m.Total - int(score) + IMDbScore
		m.Average = float64(m.Total) / float64(m.Count)
	}
	utils.Red.ZAdd(context.Background(), key, &redis.Z{Member: movieId, Score: float64(IMDbScore)})
	mykey := utils.Movie_Average_set
	utils.Red.ZAdd(context.Background(), mykey, &redis.Z{Member: movieId, Score: float64(m.Average)})
	Update(m)
	return m
}
