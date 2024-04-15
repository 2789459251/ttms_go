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

	userGroup := r.Group("/user/api")
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
	snackGroup := r.Group("/snack/api")
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

	//票务操作
	movieGroup := r.Group("/movie/api")
	movieGroup.POST("/addMovie", service.AddMovie)
	movieGroup.GET("/movieList", service.MovieList)
	movieGroup.GET("/upcoming/movieList", service.Upcoming)
	movieGroup.GET("/hit/movieList", service.Hit)
	movieGroup.DELETE("/deletemoviesByid", service.DeleteMovies)
	//充值操作
	return r
}
