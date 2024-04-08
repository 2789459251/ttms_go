package router

import (
	"TTMS_go/ttms/module/service"
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
	//snackGroup.Use(utils.JWTAuth())
	//零食操作
	snackGroup.POST("/buy", service.BuySnack)                 //购买
	snackGroup.POST("/putaway", service.Putaway)              //上架
	snackGroup.GET("/getinfos", service.ShowSnacks)           //查询零食列表
	snackGroup.GET("/search", service.SearchSnack)            //搜索特定零食
	snackGroup.GET("/query", service.Getdetail)               //根据id查询
	snackGroup.DELETE("/removeByid", service.Remove)          //下架by——id
	snackGroup.DELETE("/removeByNamekey", service.Removes)    //下架by--name
	snackGroup.PUT("/uploadFavorite", service.UploadFavorite) //零食加入收藏
	snackGroup.PUT("/updeteSnack", service.UpdateSnack)       //修改零食信息
	//票务操作

	//充值操作
	return r
}
