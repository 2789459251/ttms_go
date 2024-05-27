package models

import (
	utils "TTMS_go/ttms/util"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"sync"
	"time"
)

var play_mu sync.Mutex

type Play struct {
	gorm.Model
	MovieId   string
	TheatreId string
	Seat      string //0 1 2//数组
	Num       int    //剩余座位数量
	BeginTime time.Time
	EndTime   time.Time
}

func (play Play) TableName() string {
	return "play_basic"
}

// todo 演出返回需要注意当前时间
func CreatePlay(play *Play) Play {
	utils.DB.Create(&play)
	return *play
}
func ShowPlaysByMovieId(id string) []Play {
	plays := []Play{}
	//utils.DB.Where("movie_id = ? AND  begin_time > ?", id, time.Now()).Find(plays)
	utils.DB.Model(Play{}).Where("movie_id = ?", id).Find(&plays)
	return plays
}
func ShowPlaysByTheatreId(id string) []Play {
	plays := []Play{}
	utils.DB.Where("theatre_id = ? AND  begin_time > ?", id, time.Now()).Find(&plays)
	return plays

}
func ShowPlayById(id string) *Play {
	p := &Play{}
	utils.DB.Where("id = ?", id).Find(p)
	return p
}
func FindPlayByIds(ids []string) []Play {
	var plays []Play
	id_ := strings.Join(ids, ",")
	fmt.Println("qq:", id_)
	query := "id IN (?) AND `play_basic`.`deleted_at` IS NULL"
	utils.DB.Where(query, id_).Find(&plays)
	return plays
}

// todo kafka,websocket
func Reserve(user UserInfo, id string, seats []Seat) error {
	play := Play{}
	//查座位状态//多个座位
	utils.DB.Where("id = ?", id).Find(&play)
	playSeat, err := ConvertTo2DIntSlice(play.Seat)

	//json.Unmarshal(play.Seat, &playSeat)

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

	//defer func() {
	//	//if r := recover(); r != nil {
	//	//	tx.Rollback()
	//	//}
	//}()
	//修改//扣除
	//play_mu.Lock()
	//defer play_mu.Unlock()
	for _, seat := range seats {
		playSeat[seat.Row-1][seat.Column-1] = 1
	}
	user.Wallet -= (movie.Money * float64(len(seats)))
	movie.TicketNum += len(seats)
	if err = tx.Save(movie).Error; err != nil {
		tx.Rollback()
		return err
	}
	//保存剧目信息
	//play.Seat, _ = json.Marshal(playSeat)
	play.Seat, _ = ConvertToString(playSeat)
	if err = tx.Save(play).Error; err != nil {
		tx.Rollback()
		return err
	}
	key := utils.Movie_Ticket_Num_set
	if err := utils.Red.ZAdd(context.Background(), key, &redis.Z{
		Member: strconv.Itoa(int(movie.ID)),
		Score:  float64(movie.TicketNum),
	}).Err(); err != nil {

		tx.Rollback()
		return err
	}
	//生成票，保存到user.tacket
	t_ := Ticket{
		Name:      movie.Name,
		Num:       len(seats),
		Begintime: play.BeginTime,
		Endtime:   play.EndTime,
	}
	t_.Seat, _ = json.Marshal(seats)
	t_.Issold = true
	err = tx.Model(t_).Create(&t_).Error
	fmt.Println(err)
	if err != nil {
		tx.Rollback()
		return err
	}
	if user.Ticket == "" {
		user.Ticket = strconv.Itoa(int(t_.ID))
	} else {

		user.Ticket = user.Ticket + " " + strconv.Itoa(int(t_.ID))
	}
	fmt.Println(user.Birthday)

	if err_ := user.Tx_RefleshUserInfo(tx); err_ != nil {
		fmt.Println(err_)
		tx.Rollback()
		return err_
	}
	fmt.Println(user.Birthday)
	return tx.Commit().Error
}
