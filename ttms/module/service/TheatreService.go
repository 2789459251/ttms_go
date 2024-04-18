package service

import (
	"TTMS_go/ttms/models"
	utils "TTMS_go/ttms/util"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"strconv"
	"strings"
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

func UploadFavoriteMovie(c *gin.Context) {
	var flag bool
	user := User(c)
	movieId := c.Params.ByName("movie_id")
	key1 := utils.Movie_user_favorite_set + movieId
	id_ := strconv.Itoa(int(user.ID))
	key2 := utils.User_Movie_favorite_set + id_
	key3 := utils.Movie_ranking_sorted_set
	err := utils.Red.Watch(context.Background(), func(tx *redis.Tx) error { //乐观锁
		var err error
		flag, err = utils.Red.SIsMember(context.Background(), key1, user.ID).Result()
		if flag {
			_, err = tx.SRem(context.Background(), key1, user.ID).Result()
			if err != nil {
				return err
			}
			_, err = tx.SRem(context.Background(), key2, movieId).Result()
			if err != nil {
				return err
			}
		} else {
			_, err = tx.SAdd(context.Background(), key1, user.ID).Result()
			if err != nil {
				return err
			}
			_, err = tx.SAdd(context.Background(), key2, movieId).Result()
			if err != nil {
				return err
			}
		}
		score, _ := utils.Red.SCard(context.Background(), key1).Result()
		if errs := utils.Red.ZScore(context.Background(), key3, movieId).Err(); errs != nil {
			if errs == redis.Nil {
				utils.Red.ZAdd(context.Background(), key3, &redis.Z{Member: movieId, Score: float64(score)})
			} else {
				fmt.Errorf("查询redis缓存有误，请及时处理。err:%v", errs.Error())
				err = errs
			}
		} else {
			utils.Red.ZIncrBy(context.Background(), key3, float64(score), movieId)
		}
		return err
	})
	if err != nil {
		utils.RespFail(c.Writer, "收藏电影失败："+err.Error())
		return
	}

	if flag {
		utils.RespOk(c.Writer, "", "已经取消收藏")
		return
	} else {
		utils.RespOk(c.Writer, "", "已经添加到收藏")
		return
	}
}

func FavoriteMovieRanking(c *gin.Context) {
	key := utils.Movie_ranking_sorted_set
	members, _ := utils.Red.ZRevRangeByScoreWithScores(context.Background(), key, &redis.ZRangeBy{
		Min:    "-inf",
		Max:    "+inf",
		Offset: 0,
		Count:  10,
	}).Result()
	m := models.RankingMovies(members)
	utils.RespOk(c.Writer, string(m), "获取到收藏前十的电影，及其收藏数量。")
}

// Total       int     `json:"total"`   // 电影的总分
// Count       int     `json:"count"`   // 评分人数
// Average     float64 `json:"average"` // 平均分
// 打分	电影 评分 人数 每个人应该评论只有一次评分计算机会,不考虑错误输入
func MarkMovie(c *gin.Context) {
	user := User(c)
	userId := strconv.Itoa(int(user.ID))
	movieId := c.Params.ByName("movie_id")      //json(movieid + num + sum )= member1
	key := utils.User_Movie_marked_set + userId //key
	star, _ := strconv.Atoi(c.Params.ByName("star"))
	IMDbScore := (star * 2) - 1
	m := models.FindMovieByid(movieId)
	m = models.UpdateMovieMark(m, IMDbScore, key, movieId)
	utils.RespOk(c.Writer, m, "评价完成")
}
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
	utils.RespOk(c.Writer, string(result), "获取到评分前十条电影，及其评分")
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
	utils.RespOk(c.Writer, string(result), "获取到票房前十条电影，及其票房")
}
