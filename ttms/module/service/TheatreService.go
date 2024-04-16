package service

import (
	"TTMS_go/ttms/models"
	utils "TTMS_go/ttms/util"
	"github.com/gin-gonic/gin"
	"strconv"
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
	t := models.Theatre{
		Name:  c.Request.FormValue("name"),
		N:     c.GetInt("columns"),
		M:     c.GetInt("rows"),
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
	p := models.Play{
		MovieId:   c.Request.FormValue("movie_id"),
		TheatreId: c.Request.FormValue("theatre_id"),
		BeginTime: c.GetTime("begin_time"),
	}
	m := models.FindMovieByid(p.MovieId)
	if m.Name == "" {
		utils.RespFail(c.Writer, "电影不存在，请再次确认")
		return
	}

	//结束时间对不对
	p.EndTime = p.BeginTime.Add(m.Duration)
	t := models.FindTheatreByid(p.TheatreId)
	if t.Name == "" {
		utils.RespFail(c.Writer, "影厅不存在，请再次确认")
		return
	}
	p.Seat = make([][]int, t.N*t.M)
	p.Num = strconv.Itoa(t.Num)

	if err := isTimeable(&t, p); err != nil {
		utils.RespFail(c.Writer, "时间冲突："+err.Error()+"请检查输入")
		return
	}
	models.UpdateTheatre(&t)
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
