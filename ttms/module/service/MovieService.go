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
	"time"
)

func AddMovie(c *gin.Context) {
	if !isLimited(c) {
		return
	}

	movie := models.Movie{
		Name:     c.Request.FormValue("name"),
		Director: c.Request.FormValue("director"),
	}
	//models.CreateMovie(movie)
	actors := c.Request.FormValue("actor")
	movie.Actor = actors
	movie.Money, _ = strconv.ParseFloat(c.Request.FormValue("money"), 64)
	duration, _ := strconv.Atoi(c.Request.FormValue("duration"))
	release_time, _ := strconv.Atoi(c.Request.FormValue("release_time"))
	movie.ReleaseTime = time.Unix(int64(release_time), 0)
	movie.Duration = int64(duration)
	var e error
	movie.Picture, e = upload(c.Request, c.Writer, c)
	if e != nil {
		utils.RespFail(c.Writer, "获取图片外链错误:"+e.Error())
		return
	}

	if err := aviliable(movie); err != nil {
		utils.RespFail(c.Writer, "上传电影数据不可用，请重新上传:"+err.Error())
		return
	}
	//if !Index.Isexist(model.MovieInfo{}.Index()) {
	//	Index.CreateIndex(model.MovieInfo{}.Index(), model.MovieInfo{}.Mapping())
	//	fmt.Println("创建了movie_index的索引！")
	//}
	//movie.Info, e = docs.CreateDoc(model.MovieInfo{Info: c.Request.FormValue("info")})
	//if e != nil {
	//	utils.RespFail(c.Writer, "创建info文档失败："+e.Error())
	//	return
	//}
	movie.Info = c.Request.FormValue("info")
	if movie.Name == "" || movie.Director == "" || movie.Actor == "" || movie.Info == "" {
		utils.RespFail(c.Writer, "添加电影时，请注意：电影的名称、导演、主演、简介不能为空！")
		return
	}
	models.CreateMovie(movie)
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
	movies, ok := models.DeleteMovieById(ids)
	if !ok {
		utils.RespFail(c.Writer, "您选择删除的电影一部也存在！请检查输入！")
		return
	}
	utils.RespOk(c.Writer, movies, "删除成功")
}

// todo 修改哦
func UpdateMoviedetail(c *gin.Context) {
	if !isLimited(c) {
		return
	}
	var err error
	n := c.Request.FormValue("num")
	nums := strings.Split(n, " ")
	movieId := c.Request.FormValue("movie_id")
	movie := models.FindMovieByid(movieId)
	if movie.Name == "" {
		utils.RespFail(c.Writer, "没有该电影，请确定电影id！")
		return
	}
	for _, num := range nums {
		switch num {
		case "1":
			movie.Name = c.Request.FormValue("name") //名字
		case "2":
			movie.Director = c.Request.FormValue("director") //导演
		case "3":
			movie.Money, _ = strconv.ParseFloat(c.Request.FormValue("money"), 64) //单价
		case "4":
			info := c.Request.FormValue("info") //简述
			//if movie.Info == "" {
			//	movie.Info, err = docs.CreateDoc(model.MovieInfo{Info: info})
			//} else {
			//	movie.Info, err = docs.UpdateMovieDoc(model.MovieInfo{}, model.MovieInfo{Info: info}, movie.Info)
			//}
			//if err != nil {
			//	break
			//}
			movie.Info = info
		case "5":
			Duration, _ := strconv.Atoi(c.Request.FormValue("duration")) //时长
			movie.Duration = int64(Duration)
		case "6":
			t, _ := strconv.Atoi(c.Request.FormValue("release_time"))
			time_ := time.Unix(int64(t), 0)
			movie.ReleaseTime = time_ //发映时间
		case "7":
			movie.Online, _ = strconv.ParseBool(c.Request.FormValue("online")) //是否在院线上映
		case "8":

			movie.Picture, _ = upload(c.Request, c.Writer, c)

		default:
			utils.RespFail(c.Writer, "注意规范num输入~")
			return
		}
	}
	if err != nil {
		utils.RespFail(c.Writer, "数据修改发生错误："+err.Error())
		return
	}
	models.Update(movie)
	utils.RespOk(c.Writer, movie, "修改数据成功")
}

func Reputaway(c *gin.Context) {
	id := c.Query("id")
	if !isLimited(c) {
		return
	}
	var movie models.Movie
	utils.DB.Raw("select * from movie_basic where id = ?", id).Scan(&movie)

	if movie.Name == "" {
		utils.RespFail(c.Writer, "没有id为"+id+"的零食!")
		return
	}
	utils.DB.Exec("UPDATE `movie_basic` SET `deleted_at`= NULL WHERE `deleted_at` IS NOT NULL AND `id`=?", id)
	utils.RespOk(c.Writer, nil, "ok")
}
func FavoriteMovieList(c *gin.Context) {
	user := User(c)
	id_ := strconv.Itoa(int(user.ID))
	key := utils.User_Movie_favorite_set + id_
	str, err := utils.Red.SMembers(context.Background(), key).Result()
	if err != nil {
		utils.RespFail(c.Writer, "从redis获取缓存失败："+err.Error())
		return
	}

	movies := models.FindMovieByIds(str)
	utils.RespOk(c.Writer, movies, "获得电影收藏列表")
}

// todo redis没存上
func UploadFavoriteMovie(c *gin.Context) {
	var flag bool
	user := User(c)
	movieId := c.Query("movie_id")
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

// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!1
func FavoriteMovieRanking(c *gin.Context) {
	key := utils.Movie_ranking_sorted_set
	members, _ := utils.Red.ZRevRangeByScoreWithScores(context.Background(), key, &redis.ZRangeBy{
		Min:    "-inf",
		Max:    "+inf",
		Offset: 0,
		Count:  10,
	}).Result()
	m := models.RankingMovies(members)
	if len(m) > 3 {
		m = m[:3]
	}
	utils.RespOk(c.Writer, m, "获取到收藏前3的电影，及其收藏数量。")
}
func MarkMovie(c *gin.Context) {
	user := User(c)
	userId := strconv.Itoa(int(user.ID))
	movieId := c.Query("movie_id")                    //json(movieid + num + sum )= member1
	key := utils.User_Movie_marked_set + ":" + userId //key
	star, _ := strconv.Atoi(c.Query("star"))
	if star < 0 || star > 5 {
		utils.RespFail(c.Writer, "注意star规范")
		return
	}
	IMDbScore := float64((star * 2) - 1)
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
	if len(result) > 3 {
		result = result[:3]
	}
	utils.RespOk(c.Writer, result, "获取到评分前3条电影，及其评分")
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
	if len(result) > 3 {
		result = result[:3]
	}
	utils.RespOk(c.Writer, result, "获取到票房前3条电影，及其票房")
}
