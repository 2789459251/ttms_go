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
	userGroup.POST("/profile", service.Profile)
	userGroup.GET("/detail", service.UserDetail)
	//登出
	//userGroup.GET("/logout",service.Logout)
	snackGroup := r.Group("/snack/api") //9
	snackGroup.Use(jwtMiddleware.JWTAuthMiddleware())
	//零食操作
	snackGroup.POST("/buy", service.BuySnack)                  //购买 ok
	snackGroup.POST("/putaway", service.Putaway)               //上架 ok
	snackGroup.GET("/getsnackList", service.ShowSnacks)        //查询零食列表 ok
	snackGroup.GET("/search", service.SearchSnack)             //搜索特定零食 ok
	snackGroup.GET("/query", service.Getdetail)                //根据id查询
	snackGroup.DELETE("/removeByid", service.Remove)           //下架by——id
	snackGroup.DELETE("/removeByNamekey", service.Removes)     //下架by--name关键字
	snackGroup.PUT("/uploadFavorite", service.UploadFavorite)  //零食加入收藏 ok
	snackGroup.PUT("/updeteSnack", service.UpdateSnack)        //修改零食信息
	snackGroup.GET("/favoriteList", service.FavoriteSnackList) //收藏零食列表 ok
	snackGroup.PUT("/recover", service.Recover)                //一键修复删除信息

	r.GET("/movie/api/movieList", service.MovieList)                       //查询所有电影 ok-
	r.GET("/movie/api/upcoming/movieList", service.Upcoming)               //查询待映电影 ok-
	r.GET("/movie/api/hit/movieList", service.Hit)                         //查询热映电影 ok-
	r.GET("/movie/api/favoriteMovieRanking", service.FavoriteMovieRanking) //收藏排行榜 ok
	r.GET("/movie/api/averageMovieRanking", service.AverageMovieRanking)   //评分排行榜 ok
	r.GET("/movie/api/ticketNumRanking", service.TicketNumRanking)         //票房排行榜 ok
	//电影操作
	movieGroup := r.Group("/movie/api") //5
	movieGroup.Use(jwtMiddleware.JWTAuthMiddleware())
	movieGroup.POST("/addMovie", service.AddMovie)                 //上架电影 ok
	movieGroup.DELETE("/deletemoviesByid", service.DeleteMovies)   //下架电影 ok
	movieGroup.PUT("/markMovie", service.MarkMovie)                //评分 ok
	movieGroup.PUT("/uploadFavorite", service.UploadFavoriteMovie) //电影收藏 ok
	movieGroup.GET("/favoriteList", service.FavoriteMovieList)     //用户的收藏 ok
	movieGroup.PUT("/reputaway", service.Reputaway)
	movieGroup.PUT("/updateMoviedetail", service.UpdateMoviedetail) //修改电影信息 ok
	//theatre
	r.GET("/theatre/api/showPlayDetails", service.ShowPlayDetails) // 查询电影细节	ok

	theatreGroup := r.Group("/theatre/api") //
	theatreGroup.Use(jwtMiddleware.JWTAuthMiddleware())
	theatreGroup.POST("/addTheatre", service.AddTheatre)                    //添加放映厅 ok
	theatreGroup.GET("/showPlaysByMovieId", service.ShowPlaysByMovieId)     //查询某电影的放映安排 （用户）
	theatreGroup.GET("/showPlaysByTheatreId", service.ShowPlaysByTheatreId) //查询某影厅的放映安排

	//play
	playGroup := r.Group("/play/api")
	playGroup.Use(jwtMiddleware.JWTAuthMiddleware())
	playGroup.POST("/addPlay", service.AddPlay)    //添加剧目放映 ok
	playGroup.PUT("/buyTicket", service.BuyTicket) //买票 ok
	//other
	othersGroup := r.Group("/others/api")
	othersGroup.Use(jwtMiddleware.JWTAuthMiddleware())
	othersGroup.PUT("/recharge", service.Recharge) // 充值 ok
	//充值操作
	return r
}
