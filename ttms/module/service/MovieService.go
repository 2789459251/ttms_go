package service

import (
	"TTMS_go/ttms/models"
	utils "TTMS_go/ttms/util"
	"github.com/gin-gonic/gin"
	"strings"
)

//	type Movie struct {
//		gorm.Model
//		Info     string
//		Name     string
//		Director string
//		Actor    string
//		Score    float64	//todo 评分
//		Duration time.Duration
//		ReleaseTime time.Time
//		Money    float64
//	}
func AddMovie(c *gin.Context) {
	if !isLimited(c) {
		return
	}
	movie := models.Movie{
		Name:        c.Request.FormValue("name"),
		Director:    c.Request.FormValue("director"),
		Money:       float64(c.GetFloat64("money")),
		Info:        c.Request.FormValue("info"),
		Duration:    c.GetDuration("duration"),
		ReleaseTime: c.GetTime("release_time"),
	}
	if aviliable(movie) != nil {
		utils.RespFail(c.Writer, "上传电影数据不可用，请重新上传")
	}
	models.Update(movie)
	utils.RespOk(c.Writer, movie, "电影上架成功")
}

func MovieList(c *gin.Context) {
	m := models.MovieList()
	utils.RespOk(c.Writer, m, "返回所有电影")
}

func Upcoming(c *gin.Context) {
	m := models.UpcommingList()
	utils.RespOk(c.Writer, m, "返回待映电影")
}

func Hit(c *gin.Context) {
	m := models.HitList()
	utils.RespOk(c.Writer, m, "返回热映电影")
}

func DeleteMovies(c *gin.Context) {
	if !isLimited(c) {
		return
	}
	id := c.Query("id")
	ids := strings.Split(id, " ")
	models.DeleteMovieById(ids)
}
