package service

import (
	models2 "TTMS_go/ttms/models"
	utils "TTMS_go/ttms/util"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

func BuySnack(c *gin.Context) {
	user := User(c)
	id_ := c.Request.FormValue("id")
	num_ := c.Request.FormValue("num")
	num, _ := strconv.Atoi(num_)

	s := models2.QuerysnackByid(id_)
	// 读锁
	stock := s.GetStock()

	if user.Wallet < s.Price*float64(num) {
		utils.RespFail(c.Writer, "您的账户余额不足，请充值")
		return
	}
	if num > stock {
		utils.RespFail(c.Writer, "库存不足"+num_)
		return
	}

	//进入写锁
	s.UpdateStock(func() (err error) {
		s.Stock -= num
		user.Wallet -= s.Price * float64(num)
		if user.Snack == "" {
			user.Snack = strconv.Itoa(int(s.ID))
		} else {
			user.Snack = user.Snack + " " + strconv.Itoa(int(s.ID))
		}
		//Todo 开启事务
		err = utils.DB.Transaction(
			func(tx *gorm.DB) (err error) {
				// 进行数据库操作
				if err := user.Tx_RefleshUserInfo(tx); err != nil {
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
		return
	})

}

func ShowSnacks(c *gin.Context) {
	snack := models2.Showsnacks()
	utils.RespOk(c.Writer, snack, "返回所有零食")
}

// 查询特定名称零食
func SearchSnack(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		utils.RespFail(c.Writer, "名字不能为空")
		return
	}
	snack := models2.SearchSnack(name)
	utils.RespOk(c.Writer, snack, "返回相关零食")
}

// 上架零食 + 更新信息
func Putaway(c *gin.Context) {
	if !isLimited(c) {
		return
	}
	if len(models2.SearchSnack(c.Request.FormValue("name"))) != 0 {
		utils.RespFail(c.Writer, "您上架的零食已存在，请重新上传")
		return
	}
	r := c.Request
	w := c.Writer
	url, err := upload(r, w, c)
	if err != nil {
		utils.RespFail(c.Writer, err.Error())
		return
	}

	stock, _ := strconv.Atoi(c.Request.FormValue("stock"))
	price, _ := strconv.ParseFloat(c.Request.FormValue("price"), 64)
	snack := models2.Snack{
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
	models2.Insertsnack(snack)
	utils.RespOk(c.Writer, snack, snack.Name+"已上架")
}

func Getdetail(c *gin.Context) {
	id_ := c.Query("id")
	s := models2.QuerysnackByid(id_)
	utils.RespOk(c.Writer, s, "返回指定id零食")
}

// 下架按照id
func Remove(c *gin.Context) {
	if !isLimited(c) {
		return
	}
	id_ := c.Request.FormValue("id")
	id, _ := strconv.Atoi(id_)
	if id <= 0 {
		utils.RespFail(c.Writer, fmt.Sprintf("输入id:%v无效", id_))
		return
	}
	if err := models2.DeleteSnackByid(id); err != nil {
		utils.RespFail(c.Writer, "下架商品出错，请联系系统维护人员")
		return
	}
	utils.RespOk(c.Writer, id, "下架成功")
}

// 按照姓名关键字模糊删除
func Removes(c *gin.Context) {
	if !isLimited(c) {
		return
	}
	namekey := c.Request.FormValue("namekey")
	if err := models2.DeleteSnackByNamekey(namekey); err != nil {
		utils.RespFail(c.Writer, "下架商品出错，请联系系统维护人员")
		return
	}

	utils.RespOk(c.Writer, namekey, "下架成功")
}

func UploadFavorite(c *gin.Context) {
	userInfo := User(c)

	var flag bool
	snack_id := c.Request.FormValue("snack_id")
	key1 := utils.Snack_user_favorite_set + ":" + snack_id //一个零食受收藏人群
	userid := strconv.Itoa(int(userInfo.ID))
	key2 := utils.User_snack_favorite_set + ":" + userid //一个用户收藏的零食

	// 开始 Redis 事务
	err := utils.Red.Watch(context.Background(), func(tx *redis.Tx) error {
		flagCmd := tx.SIsMember(context.Background(), key1, userInfo.ID)

		// 根据用户是否已收藏决定添加或删除操作
		flag, _ = flagCmd.Result()

		if flag {
			tx.SRem(context.Background(), key1, userInfo.ID)
			tx.SRem(context.Background(), key2, snack_id)
		} else {
			tx.SAdd(context.Background(), key1, userInfo.ID)
			tx.SAdd(context.Background(), key2, snack_id)
		}

		return nil
	})

	if err != nil {
		utils.RespFail(c.Writer, "收藏零食："+err.Error())
		return
	}

	// 获取收藏数量
	num, _ := utils.Red.SCard(context.Background(), key1).Result()

	if flag {
		utils.RespOk(c.Writer, num, "您已取消对"+snack_id+"的收藏,data显示零食的收藏量")
	} else {
		utils.RespOk(c.Writer, num, "您收藏了"+snack_id+"的零食，可以在收藏夹中查看,data显示零食的收藏量")
	}
}

func Recover(c *gin.Context) {

	if !isLimited(c) {
		return
	}
	utils.DB.Exec("UPDATE `snack_basic` SET `deleted_at`= NULL WHERE `deleted_at` IS NOT NULL")
	utils.RespOk(c.Writer, nil, "ok")
}

//	type Snack struct {
//		gorm.Model
//		mu      sync.RWMutex
//		Name    string
//		Picture string
//		Info    string
//		Stock   int     //库存量
//		Price   float64 //价格
//	}
func UpdateSnack(c *gin.Context) {
	//todo 文件待开发，加入ES存储功能
	if !isLimited(c) {
		return
	}
	snackId := c.Request.FormValue("snack_id")
	snack := models2.QuerysnackByid(snackId)
	if snack.Name == "" {
		utils.RespFail(c.Writer, "注意输入snack_id,数据库无相关信息")
		return
	}
	var err error
	num := c.Request.FormValue("num")
	nums := strings.Split(num, " ")
	for _, num_ := range nums {
		switch num_ {
		case "1":
			snack.Name = c.Request.FormValue("name")
		case "2":
			snack.Info = c.Request.FormValue("info")
		case "3":
			snack.Picture, _ = upload(c.Request, c.Writer, c)
		case "4":
			snack.Stock, err = strconv.Atoi(c.Request.FormValue("stock"))
		case "5":
			snack.Price, err = strconv.ParseFloat(c.Request.FormValue("price"), 64)
		default:
			utils.RespFail(c.Writer, "输入不规范，请规范输入")
			return
		}
		if err != nil {
			break
		}
	}
	if err != nil {
		utils.RespFail(c.Writer, "零食信息更新出错，原因："+err.Error())
		return
	}
	snack.RefreshSnack()
	utils.RespOk(c.Writer, snack, "成功修改零食的信息")
}

func FavoriteSnackList(c *gin.Context) {
	user := User(c)
	userId := strconv.Itoa(int(user.ID))
	key := utils.User_snack_favorite_set + ":" + userId
	str, err := utils.Red.SMembers(context.Background(), key).Result()
	if err != nil {
		utils.RespFail(c.Writer, "从redis获取缓存失败："+err.Error())
		return
	}
	snacks := models2.FindSnackByIds(str)
	utils.RespOk(c.Writer, snacks, "获取到零食收藏的信息如下")
}
