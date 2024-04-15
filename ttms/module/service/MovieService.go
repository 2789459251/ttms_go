package service

import (
	"TTMS_go/ttms/models"
	utils "TTMS_go/ttms/util"
	"github.com/gin-gonic/gin"
)

//	type Movie struct {
//		gorm.Model
//		Info     string
//		Name     string
//		Director string
//		Actor    string
//		Score    float64
//		Duration time.Duration
//		ReleaseTime time.Time
//		Money    float64
//	}
func AddMovie(c *gin.Context) {
	movie := models.Movie{
		Name:        c.Request.FormValue("name"),
		Director:    c.Request.FormValue("director"),
		Money:       float64(c.GetFloat64("money")),
		Duration:    c.GetDuration("duration"),
		ReleaseTime: c.GetTime("release_time"),
	}
	if aviliable(movie) != nil {
		utils.RespFail(c.Writer, "上传电影数据不可用，请重新上传")
	}

	models.Update(movie)
	utils.RespOk(c.Writer, movie, "电影上架成功")
}
