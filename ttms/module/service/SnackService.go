package service

import (
	"TTMS_go/ttms/domain/models"
	"TTMS_go/ttms/domain/models/dao"
	utils "TTMS_go/ttms/util"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strconv"
)

func BuySnack(c *gin.Context) {
	_, user := User(c)
	id_ := c.Request.FormValue("id")
	num_ := c.Request.FormValue("num")
	id, _ := strconv.Atoi(id_)
	num, _ := strconv.Atoi(num_)
	s := models.Querysnack(id)
	//Todo 设计读写锁
	if user.Wallet < s.Price*float64(num) {
		utils.RespFail(c.Writer, "您的账户余额不足，请充值")
		return
	}
	if num > s.Stock {
		utils.RespFail(c.Writer, "库存不足"+num_)
		return
	}
	s.Stock -= num
	user.Wallet -= s.Price * float64(num)

	s_ := dao.Snack_{
		Id:   s.ID,
		Name: s.Name,
		Num:  num,
	}
	user.Snack = append(user.Snack, s_)
	//Todo 开启事务
	err := utils.DB.Transaction(
		func(tx *gorm.DB) (err error) {
			// 进行数据库操作
			if err := user.RefleshUserInfo(); err != nil {
				// 发生错误，进行回滚
				tx.Rollback()
				return err
			}

			if err := s.Refleshsnack(); err != nil {
				// 发生错误，进行回滚
				tx.Rollback()
				return err
			}

			// 没有错误，提交事务
			return nil
		})
	if err != nil {
		utils.RespFail(c.Writer, "购买失败")
		return
	}

	utils.RespOk(c.Writer, user.Snack, "已购买"+num_+"份"+s.Name)
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
func Getdetail(c *gin.Context) {
	id_ := c.Param("id")
	id, _ := strconv.Atoi(id_)
	s := models.Querysnack(id)
	utils.RespOk(c.Writer, s, "返回指定id零食")
}
