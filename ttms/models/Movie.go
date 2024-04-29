package models

import (
	"TTMS_go/ttms/models/docs"
	"TTMS_go/ttms/models/model"
	utils "TTMS_go/ttms/util"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"sort"
	"sync"
	"time"
)

var mu sync.RWMutex

// todo 在电影类可以加预告片的放映,图片组？？？。
type Movie struct {
	gorm.Model
	Picture     string
	Info        string
	Name        string
	Director    string
	Actor       string
	Duration    int64
	ReleaseTime time.Time `json:"release_time,omitempty" gorm:"type:datetime;default:null"`
	Money       float64   `gorm:"default:0.0"`
	Online      bool
	TicketNum   int     `json:"ticket_num"`
	Total       float64 `json:"total" gorm:"default:0.0"`
	Count       int     `json:"count"`
	Average     float64 `json:"average" gorm:"default:0.0"`
}

func (movie Movie) TableName() string {
	return "movie_basic"
}

func CreateMovie(movie Movie) {
	utils.DB.Create(&movie)
}
func FindMovieByid(id string) Movie {
	m := Movie{}
	utils.DB.Where("id = ?", id).Find(&m)
	return m
}

func Update(m Movie) {
	mu.Lock()
	defer mu.Unlock()
	//utils.DB.Where("id = ?", m.ID).Find(&m)
	utils.DB.Save(&m)
}

func MovieList() []Movie {
	m := []Movie{}
	utils.DB.Exec("select * from movie_basic").Find(&m)

	return m
}

//func (m *Movie) AfterFind(tx *gorm.DB) error {
//	// 将数据库中存储的 JSON 字符串解码到 []string 中
//	var actors []string
//	json.Marshal(m.Actor)
//	if err := json.Unmarshal(actorBytes, &actors); err != nil {
//		return err
//	}
//	// 将解码后的值赋给 Movie 结构体的 Actor 字段
//	m.Actor = actors
//	return nil
//}

func UpcommingList() []Movie {
	m := []Movie{}
	utils.DB.Order("release_time ASC").Where("release_time  > ?", time.Now()).Find(&m)
	return m
}

func HitList() []Movie {
	m := []Movie{}
	utils.DB.Order("average ASC").Where("release_time < ?", time.Now()).Find(&m)
	return m
}

func DeleteMovieById(ids []string) ([]Movie, bool) {
	IDS := []string{}
	m := []Movie{}
	utils.DB.Where("id in (?)", ids).Find(&m)
	if len(m) == 0 {
		return nil, false
	}
	utils.DB.Where("id in (?)", ids).Delete(&m)

	for _, movie := range m {
		IDS = append(IDS, movie.Info)
	}
	docs.DeleteDocs(model.MovieInfo{}, IDS)
	return m, true
}

type MovieWithScore struct {
	M     Movie
	Score float64
}

func RankingMovies(members []redis.Z) []MovieWithScore {

	result := []MovieWithScore{}
	for _, member := range members {
		res := Movie{}
		utils.DB.Where("id = (?)", member.Member.(string)).Find(&res)
		result = append(result, MovieWithScore{
			M:     res,
			Score: member.Score,
		})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Score > result[j].Score // 降序排列
	})
	return result
}

func UpdateMovieMark(m Movie, IMDbScore float64, key string, movieId string) Movie {
	//todo 事务
	mykey := utils.Movie_Average_set
	if err := utils.Red.ZScore(context.Background(), key, movieId).Err(); err != nil {
		//没评价过
		if err == redis.Nil {
			m.Total += IMDbScore
			m.Count++
			m.Average = m.Total / float64(m.Count)
		} else {
			fmt.Errorf("查询用户评分出现错误，err:%v", err.Error())
		}

	} else {
		//评价过
		score, _ := utils.Red.ZScore(context.Background(), key, movieId).Result()
		m.Total = m.Total - score + IMDbScore
		m.Average = m.Total / float64(m.Count)
	}
	utils.Red.ZAdd(context.Background(), key, &redis.Z{Member: movieId, Score: IMDbScore})

	utils.Red.ZAdd(context.Background(), mykey, &redis.Z{Member: movieId, Score: m.Average})
	fmt.Println(m.Average)
	utils.DB.Model(&m).Where("id = ?", m.ID).Updates(map[string]interface{}{
		"total":   m.Total,
		"count":   m.Count,
		"average": m.Average,
	})
	return m
}

func FindMovieByIds(ids []string) []Movie {
	movies := []Movie{}
	for _, id := range ids {
		fmt.Println(id)
		movie := Movie{}
		utils.DB.Where("id = ?", id).Find(&movie)
		fmt.Println(movie)
		movies = append(movies, movie)
	}
	return movies
}
