package service

import (
	"TTMS_go/ttms/models"
	utils "TTMS_go/ttms/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
	"time"
)

type PlayWithTheatre struct {
	Play        models.Play
	TheatreName string
	MovieName   string
}

func ShowPlaysByMovieId(c *gin.Context) {
	id := c.Query("movie_id")

	p := models.ShowPlaysByMovieId(id)
	if len(p) == 0 {
		utils.RespOk(c.Writer, "", "该电影目前没有演出～")
		return
	}

	var plays []PlayWithTheatre
	for _, value := range p {
		var pp PlayWithTheatre
		pp.Play = value
		pp.TheatreName = models.FindTheatreByid(value.TheatreId).Name
		plays = append(plays, pp)
	}
	utils.RespOk(c.Writer, plays, "返回电影的放映安排。")
}
func ShowPlaysByTheatreId(c *gin.Context) {
	id := c.Query("theatre_id")
	p := models.ShowPlaysByTheatreId(id)
	if len(p) == 0 {
		utils.RespOk(c.Writer, "", "该放映厅目前没有演出～")
		return
	}
	plays := make([]*PlayWithTheatre, 0)
	for _, value := range p {
		pp := &PlayWithTheatre{
			Play: value,
		}
		movieId, err := strconv.Atoi(value.MovieId)
		if err != nil {
			utils.RespFail(c.Writer, "注意检查数据！")
			return
		}
		pp.MovieName = models.FindMovieById(movieId).Name
		plays = append(plays, pp)
	}
	utils.RespOk(c.Writer, plays, "返回影院的放映安排。")
}

func ShowPlayDetails(c *gin.Context) {
	id := c.Query("play_id")
	p := models.ShowPlayById(id)
	if len(p.Seat) == 0 {
		utils.RespFail(c.Writer, "没有查询到对应的剧目信息。")
		return
	}
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
	fmt.Println(user.Birthday)
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
	fmt.Println("a:", treatre)
	p.Seat = treatre.Seat
	p.Num = treatre.Num

	p = models.CreatePlay(&p)
	fmt.Println("pid:", p)
	if err := isTimeable(&treatre, p); err != nil {
		utils.DB.Exec("DELETE FROM play_basic where id = ?", p.ID)
		utils.RespFail(c.Writer, "时间冲突："+err.Error()+"请检查输入")
		return
	}
	fmt.Println(treatre.Plays)
	utils.DB.Save(treatre)
	utils.RespOk(c.Writer, p, "演出安排成功。")
}
