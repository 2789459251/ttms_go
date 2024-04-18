package service

import (
	"TTMS_go/ttms/models"
	utils "TTMS_go/ttms/util"
	"context"
	"github.com/gin-gonic/gin"
	"strconv"
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
	var e error
	movie.Picture, e = upload(c.Request, c.Writer, c)
	if e != nil {
		utils.RespFail(c.Writer, "获取图片外链错误:"+e.Error())
		return
	}
	if aviliable(movie) != nil {
		utils.RespFail(c.Writer, "上传电影数据不可用，请重新上传")
		return
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

func UpdateMoviedetail(c *gin.Context) {
	if !isLimited(c) {
		return
	}
	n := c.Params.ByName("Num")
	nums := strings.Split(n, " ")
	movieId := c.Params.ByName("movie_id")
	movie := models.FindMovieByid(movieId)
	for _, num := range nums {
		switch num {
		case "1":
			movie.Name = c.Request.FormValue("name") //名字
		case "2":
			movie.Director = c.Request.FormValue("director") //导演
		case "3":
			movie.Money = float64(c.GetFloat64("money")) //单价
		case "4":
			movie.Info = c.Request.FormValue("info") //简述
		case "5":
			movie.Duration = c.GetDuration("duration") //时长
		case "6":
			movie.ReleaseTime = c.GetTime("release_time") //发映时间
		case "7":
			movie.Online = c.GetBool("online") //是否在院线上映
		default:
			utils.RespFail(c.Writer, "注意规范num输入~")
			return
		}
	}

	models.Update(movie)
	utils.RespOk(c.Writer, movie, "修改数据成功")
}

func FavoriteMovieList(c *gin.Context) {
	user := User(c)
	id_ := strconv.Itoa(int(user.ID))
	key := utils.User_Movie_favorite_set + id_
	//1) "A"
	//2) "B"
	//3) "C"

	str, err := utils.Red.SMembers(context.Background(), key).Result()
	if err != nil {
		utils.RespFail(c.Writer, "从redis获取缓存失败："+err.Error())
		return
	}

	movies := models.FindMovieByIds(str)
	utils.RespOk(c.Writer, movies, "获得电影收藏列表")
}
