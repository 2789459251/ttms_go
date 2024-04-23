package models

import (
	utils "TTMS_go/ttms/util"
	"context"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"strconv"
	"sync"
	"time"
)

var play_mu sync.Mutex

type Play struct {
	gorm.Model
	MovieId   string
	TheatreId string
	Seat      string ` gorm:"type:json"` //0 1 2//数组
	Num       int    //剩余座位数量
	BeginTime time.Time
	EndTime   time.Time
}

func (play Play) TableName() string {
	return "play_basic"
}

// todo 演出返回需要注意当前时间
func CreatePlay(play *Play) {
	utils.DB.Create(&play)
}
func ShowPlaysByMovieId(id string) []Play {
	plays := []Play{}
	//utils.DB.Where("movie_id = ? AND  begin_time > ?", id, time.Now()).Find(plays)
	utils.DB.Where("movie_id = ?", id).Find(plays)
	return plays
}
func ShowPlaysByTheatreId(id string) []Play {
	plays := []Play{}
	utils.DB.Where("theatre_id = ? AND  begin_time > ?", id, time.Now()).Find(plays)
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
	var playSeat [][]int
	json.Unmarshal([]byte(play.Seat), &playSeat)

	for _, seat := range seats {
		if playSeat[seat.Row-1][seat.Column-1] != 0 {
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
	play_mu.Lock()
	defer play_mu.Unlock()
	for _, seat := range seats {
		playSeat[seat.Row-1][seat.Column-1] = 1
	}
	user.Wallet -= (movie.Money * float64(len(seats)))
	movie.TicketNum += len(seats)
	utils.DB.Save(movie)
	//保存剧目信息
	s_, _ := json.Marshal(playSeat)
	play.Seat = string(s_)
	utils.DB.Save(play)

	key := utils.Movie_Ticket_Num_set
	utils.Red.ZAdd(context.Background(), key, &redis.Z{
		Member: strconv.Itoa(int(movie.ID)),
		Score:  float64(movie.TicketNum),
	})
	//生成票，保存到user.tacket
	ticket := CreateTicket(play, movie, seats)
	user.Ticket = append(user.Ticket, ticket)
	user.RefleshUserInfo_()
	return tx.Commit().Error
}
