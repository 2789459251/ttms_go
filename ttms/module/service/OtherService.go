package service

import (
	"TTMS_go/ttms/models"
	utils "TTMS_go/ttms/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

func Recharge(c *gin.Context) {
	user := User(c)
	num_ := c.PostForm("num")
	num, _ := strconv.ParseFloat(num_, 64)
	fmt.Println(num)
	user.Wallet += num
	user = models.Recharge(user.Wallet, user)
	utils.RespOk(c.Writer, user, "充值成功")
}
