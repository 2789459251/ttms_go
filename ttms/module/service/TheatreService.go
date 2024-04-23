package service

import (
	"TTMS_go/ttms/models"
	utils "TTMS_go/ttms/util"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"strconv"
	"strings"
	"time"
)

//type Theatre struct {
//	gorm.Model
//	Name  string
//	Seat  [][]bool `gorm:"type:json"`
//	N     int
//	M     int
//	Inuse int//影院状态 0 1 2
//	Num   int //座位数量
//}

func AddTheatre(c *gin.Context) {
	if !isLimited(c) {
		return
	}
	column, _ := strconv.Atoi(c.Request.FormValue("columns"))
	row, _ := strconv.Atoi(c.Request.FormValue("rows"))
	t := models.Theatre{
		Name:  c.Request.FormValue("name"),
		N:     column,
		M:     row,
		Inuse: 0,
	}
	t.Num = t.N * t.M
	if t.Num <= 0 || t.Num > 50 {
		utils.RespFail(c.Writer, "无效添加操作，请检查输入是否合法！")
		return
	}
	t.Info = c.Request.FormValue("info")
	if t.Num <= 25 {
		t.Info += "全景声MAX激光厅"
	} else {
		t.Info += "RGB基色激光厅"
	}
	models.CreateTheatre(t)
	utils.RespOk(c.Writer, t, "添加放映厅成功")
	return
}

//type Play struct {
//	gorm.Model
//	MovieId   int
//	TheatreId int
//	Seat      [][]int //0 1 2
//	Num       int     //剩余座位数量
//	BeginTime time.Time
//}

func AddPlay(c *gin.Context) {
	if !isLimited(c) {
		return
	}
	p := models.Play{
		MovieId:   c.Request.FormValue("movie_id"),
		TheatreId: c.Request.FormValue("theatre_id"),
	}
	t, _ := strconv.Atoi(c.Request.FormValue("begin_time"))
	p.BeginTime = time.Unix(int64(t), 64)
	m := models.FindMovieByid(p.MovieId)
	if m.Name == "" {
		utils.RespFail(c.Writer, "电影不存在，请再次确认")
		return
	}

	if p.BeginTime.Before(m.ReleaseTime) {
		utils.RespFail(c.Writer, "未到该电影的放映时间！")
		return
	}
	//结束时间对不对
	p.EndTime = p.BeginTime.Add(time.Duration(m.Duration) * time.Second)

	treatre := models.FindTheatreByid(p.TheatreId)
	if treatre.Name == "" {
		utils.RespFail(c.Writer, "影厅不存在，请再次确认")
		return
	}
	p.Seat = make([][]int, treatre.N*treatre.M)
	p.Num = strconv.Itoa(treatre.Num)

	if err := isTimeable(&treatre, p); err != nil {
		utils.RespFail(c.Writer, "时间冲突："+err.Error()+"请检查输入")
		return
	}
	models.UpdateTheatre(&treatre)
	models.CreatePlay(&p)
	utils.RespOk(c.Writer, p, "演出安排成功。")
}

func ShowPlaysByMovieId(c *gin.Context) {
	id := c.Params.ByName("movie_id")
	p := models.ShowPlaysByMovieId(id)
	if len(p) == 0 {
		utils.RespOk(c.Writer, "", "该电影目前没有演出～")
		return
	}
	utils.RespOk(c.Writer, p, "返回电影的放映安排。")
}
func ShowPlaysByTheatreId(c *gin.Context) {
	id := c.Params.ByName("theatre_id")
	p := models.ShowPlaysByTheatreId(id)
	if len(p) == 0 {
		utils.RespOk(c.Writer, "", "该电影目前没有演出～")
		return
	}
	utils.RespOk(c.Writer, p, "返回影院的放映安排。")
}

func ShowPlayDetails(c *gin.Context) {
	id := c.Params.ByName("play_id")
	p := models.ShowPlayById(id)
	m := models.FindMovieByid(p.MovieId)
	t := models.FindTheatreByid(p.TheatreId)
	var response []interface{}
	response = append(response, p)
	response = append(response, m)
	response = append(response, t)

	utils.RespOk(c.Writer, response, "获得放映场次具体数据")
}

// kafka
func BuyTicket(c *gin.Context) {
	playId := c.PostForm("play_id")
	column := c.Request.FormValue("column")
	columns := strings.Split(column, " ")
	raw := c.PostForm("row")
	raws := strings.Split(raw, " ")
	user := User(c)
	seats := []models.Seat{}
	for i, _ := range raws {
		seat := models.Seat{}
		seat.Column, _ = strconv.Atoi(columns[i])
		seat.Row, _ = strconv.Atoi(raws[i])
		seats = append(seats, seat)
	}
	err := models.Reserve(user, playId, seats)
	if err != nil {
		utils.RespFail(c.Writer, "发生err :"+err.Error())
		return
	}
	utils.RespOk(c.Writer, "", "购票成功")
}

// Total       int     `json:"total"`   // 电影的总分
// Count       int     `json:"count"`   // 评分人数
// Average     float64 `json:"average"` // 平均分
// 打分	电影 评分 人数 每个人应该评论只有一次评分计算机会,不考虑错误输入

func AverageMovieRanking(c *gin.Context) {
	key := utils.Movie_Average_set
	members, _ := utils.Red.ZRevRangeByScoreWithScores(context.Background(), key,
		&redis.ZRangeBy{
			Min:    "-inf",
			Max:    "+inf",
			Offset: 0,
			Count:  10,
		}).Result()
	result := models.RankingMovies(members)
	utils.RespOk(c.Writer, result, "获取到评分前十条电影，及其评分")
}

func TicketNumRanking(c *gin.Context) {
	key := utils.Movie_Ticket_Num_set
	members, _ := utils.Red.ZRevRangeByScoreWithScores(context.Background(), key,
		&redis.ZRangeBy{
			Min:    "-inf",
			Max:    "+inf",
			Offset: 0,
			Count:  10,
		}).Result()
	result := models.RankingMovies(members)
	utils.RespOk(c.Writer, result, "获取到票房前十条电影，及其票房")
}
