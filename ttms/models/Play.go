package models

import (
	utils "TTMS_go/ttms/util"
	"errors"
	"gorm.io/gorm"
	"sync"
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
	mu        sync.Mutex
}

// todo 演出返回需要注意当前时间
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

// todo kafka,websocket
func Reserve(user UserInfo, id string, seats []Seat) error {
	play := Play{}
	//查座位状态//多个座位
	utils.DB.Where("id = ?", id).Find(&play)
	for _, seat := range seats {
		if play.Seat[seat.Row-1][seat.Column-1] != 0 {
			return errors.New("座位不可用或已被预定！")
		}
	}

	//查余额
	movie := FindMovieByid(play.MovieId)
	if movie.Money*float64(len(seats)) > user.Wallet {
		return errors.New("余额不足！")
	}
	//开启事务
	tx := utils.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	//修改//扣除
	play.mu.Lock()
	defer play.mu.Unlock()
	for _, seat := range seats {
		play.Seat[seat.Row-1][seat.Column-1] = 1
	}
	user.Wallet -= (movie.Money * float64(len(seats)))

	//保存剧目信息
	utils.DB.Save(play)

	//生成票，保存到user.tacket
	ticket := CreateTicket(play, movie, seats)
	user.Ticket = append(user.Ticket, ticket)
	user.RefleshUserInfo_()
	return tx.Commit().Error
}
