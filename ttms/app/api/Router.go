package router

import (
	"TTMS_go/ttms/module/service"
	utils "TTMS_go/ttms/util"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	jwtMiddleware, _ := utils.InitAuth()
	r := gin.Default()
	r.Static("/Asset", "Asset/")

	userGroup := r.Group("/user/api") //7
	//用户登录注册
	userGroup.POST("/createUser", service.CreateUser)
	userGroup.POST("/loginByPassword", service.LoginByPassword)
	userGroup.POST("/sendCode", service.SendCode)
	userGroup.POST("/loginByCode", service.LoginByCode)
	userGroup.POST("/resetPassword", service.ResetPassword)
	userGroup.GET("/refreshToken", jwtMiddleware.RefreshHandler)
	userGroup.PUT("/admin", service.Admin)
	//登出
	//userGroup.GET("/logout",service.Logout)
	snackGroup := r.Group("/snack/api") //9
	snackGroup.Use(jwtMiddleware.JWTAuthMiddleware())
	//零食操作
	snackGroup.POST("/buy", service.BuySnack)                 //购买
	snackGroup.POST("/putaway", service.Putaway)              //上架 1
	snackGroup.GET("/getsnackList", service.ShowSnacks)       //查询零食列表
	snackGroup.GET("/search", service.SearchSnack)            //搜索特定零食
	snackGroup.GET("/query", service.Getdetail)               //根据id查询
	snackGroup.DELETE("/removeByid", service.Remove)          //下架by——id 1
	snackGroup.DELETE("/removeByNamekey", service.Removes)    //下架by--name关键字 1
	snackGroup.PUT("/uploadFavorite", service.UploadFavorite) //零食加入收藏
	//snackGroup.PUT("/updeteSnack", service.UpdateSnack)       //修改零食信息 1
	snackGroup.PUT("/recover", service.Recover) //一键修复删除信息 1

	//电影操作
	movieGroup := r.Group("/movie/api") //5
	movieGroup.Use(jwtMiddleware.JWTAuthMiddleware())
	movieGroup.POST("/addMovie", service.AddMovie)                        //上架电影
	movieGroup.GET("/movieList", service.MovieList)                       //查询所有电影
	movieGroup.GET("/upcoming/movieList", service.Upcoming)               //查询待映电影
	movieGroup.GET("/hit/movieList", service.Hit)                         //查询热映电影
	movieGroup.DELETE("/deletemoviesByid", service.DeleteMovies)          //下架电影
	movieGroup.PUT("/markMovie", service.MarkMovie)                       //评分
	movieGroup.PUT("/uploadFavorite", service.UploadFavoriteMovie)        //电影收藏
	movieGroup.GET("/favoriteList", service.FavoriteMovieList)            //用户的收藏
	movieGroup.GET("/favoriteMovieRanking", service.FavoriteMovieRanking) //收藏排行榜
	movieGroup.GET("/averageMovieRanking", service.AverageMovieRanking)   //评分排行榜
	movieGroup.GET("/ticketNumRanking", service.TicketNumRanking)         //票房排行榜
	movieGroup.PUT("/updateMoviedetail", service.UpdateMoviedetail)
	//theatre
	theatreGroup := r.Group("/theatre/api") //9
	theatreGroup.Use(jwtMiddleware.JWTAuthMiddleware())
	theatreGroup.POST("/addTheatre", service.AddTheatre)                //添加放映厅
	theatreGroup.POST("/addPlay", service.AddPlay)                      //安排放映
	theatreGroup.GET("/showPlaysByMovieId", service.ShowPlaysByMovieId) //查询某电影的放映安排
	theatreGroup.GET("/showPlaysByTheatreId", service.ShowPlaysByTheatreId)
	theatreGroup.GET("/showPlayDetails", service.ShowPlayDetails) // 查询电影细节
	theatreGroup.PUT("/buyTicket", service.BuyTicket)             //买票

	//充值操作
	return r
}
