package service

import (
	"TTMS_go/ttms/domain/models"
	utils "TTMS_go/ttms/util"
	"github.com/gin-gonic/gin"
	"strconv"
)

func BuySnack(c *gin.Context) {
	//_, user := User(c)

}

func ShowSnacks(c *gin.Context) {
	snack := models.Showsnacks()
	utils.RespOk(c.Writer, snack, "返回所有零食")
}

// 查询特定名称零食
func SearchSnack(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		utils.RespFail(c.Writer, "名字不能为空")
		return
	}
	snack := models.SearchSnack(name)
	utils.RespOk(c.Writer, snack, "返回相关零食")
}

// 上架零食
func Putaway(c *gin.Context) {
	r := c.Request
	w := c.Writer
	url := upload(r, w, c)
	stock, _ := strconv.Atoi(c.Request.FormValue("stock"))
	price, _ := strconv.ParseFloat(c.Request.FormValue("price"), 64)
	snack := models.Snack{
		Name:    c.Request.FormValue("name"),
		Picture: url,
		Info:    c.Request.FormValue("info"),
		Stock:   stock,
		Price:   price,
	}
	if snack.Name == "" {
		utils.RespFail(c.Writer, "名字不能为空")
		return
	}
	if price < 0.0 {
		utils.RespFail(c.Writer, "价格不能小于0")
		return
	}
	if snack.Info == "" {
		utils.RespFail(c.Writer, "描述不能为空")
		return
	}
	models.Insertsnack(snack)
	utils.RespOk(c.Writer, snack, snack.Name+"已上架")
}
