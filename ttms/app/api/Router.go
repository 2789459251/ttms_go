package router

import (
	"TTMS_go/ttms/module/service"
	utils "TTMS_go/ttms/util"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	r := gin.Default()
	r.Static("/Asset", "Asset/")

	userGroup := r.Group("/user/api")
	//用户登录注册
	userGroup.POST("/createUser", service.CreateUser)
	userGroup.POST("/loginByPassword", service.LoginByPassword)
	userGroup.POST("/sendCode", service.SendCode)
	userGroup.POST("/loginByCode", service.LoginByCode)
	userGroup.POST("/resetPassword", service.ResetPassword)

	snackGroup := r.Group("/snack/api")
	snackGroup.Use(utils.JWTAuth())
	//零食操作
	snackGroup.POST("/buy", service.BuySnack)
	snackGroup.POST("/putaway", service.Putaway)
	snackGroup.GET("/getinfos", service.ShowSnacks) //查询零食列表
	snackGroup.GET("/search", service.SearchSnack)  //搜索特定零食
	snackGroup.GET("/query", service.Getdetail)
	snackGroup.DELETE("/removeByid", service.Remove)
	snackGroup.DELETE("/removeByNamekey", service.Removes)
	snackGroup.GET("/uploadFavorite", service.UploadFavorite) //零食加入收藏

	//票务操作

	//充值操作
	return r
}
