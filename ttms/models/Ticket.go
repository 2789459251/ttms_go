package models

import (
	utils "TTMS_go/ttms/util"
	"gorm.io/gorm"
	"time"
)

type Ticket struct {
	gorm.Model
	//影片名
	Name string
	Num  int
	//演出厅
	Place int
	//座位
	Seat []Seat `gorm:"type:json"`
	//状态
	Issold bool
	//影片开始结束时间
	Begintime time.Time
	Endtime   time.Time
}

type Seat struct {
	Column int
	Row    int
}

func (ticket Ticket) TableName() string {
	return "ticket_basic"
}

func CreateTicket(play Play, movie Movie, seat []Seat) Ticket {
	t_ := Ticket{
		Name:      movie.Name,
		Num:       len(seat),
		Begintime: play.BeginTime,
		Endtime:   play.EndTime,
	}
	t_.Seat = seat
	t_.Issold = true
	utils.DB.Create(t_)
	return t_
}
