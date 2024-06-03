package service

import (
	"TTMS_go/ttms/models"
	utils "TTMS_go/ttms/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/spf13/viper"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func CreateUser(c *gin.Context) {
	user := models.User{}
	phone := c.Request.FormValue("phone")
	user2 := models.FindUserByPhone(phone)
	if user2.Password != "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1, //失败
			"message": "该号码已被使用",
			"data":    nil,
		})
		return
	}
	user.Phone = phone

	if !isMatchPhone(user.Phone) {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1, //失败
			"message": "电话号码无效",
			"data":    nil,
		})
		return
	}
	password := c.Request.FormValue("password")
	repassword := c.Request.FormValue("repassword")
	if !isStrongPassword(password) {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1, //失败
			"message": "密码无效,请输入长度在8-16位的字母数字或特殊字符",
			"data":    nil,
		})
		return
	}
	if password != repassword {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1, //失败
			"message": "密码不一致",
			"data":    nil,
		})
		return
	} else {
		user.Password, _ = utils.GetPwd(password)
		models.CreateUser(user)

		//c.Redirect(http.StatusMovedPermanently, "/user/api/loginByPassword")
		c.JSON(http.StatusOK, gin.H{
			"code":    0, //成功
			"message": "登陆成功",
			"data":    user,
		})
		return
	}

}

func LoginByPassword(c *gin.Context) {
	//不要明文存储密码=
	phone := c.Request.FormValue("phone")
	user := models.FindUserByPhone(phone)
	fmt.Println(user)
	if user.Password == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1, //失败
			"message": "用户尚未注册",
			"data":    nil,
		})
		return
	}
	password := c.Request.FormValue("password")
	if utils.ComparePwd(user.Password, password) {
		if !signed(user, c) {
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":    0, //成功
			"message": "欢迎回来",
			"data":    user,
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1, //失败
			"message": "密码错误",
			"data":    nil,
		})
		return
	}
}

func SendCode(c *gin.Context) {
	//post请求->phone
	phone := c.Request.FormValue("phone")
	if !isMatchPhone(phone) {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": "手机号码无效",
			"data":    nil,
		})
		return
	}
	code := utils.GenerateSMSCode()
	fmt.Println("验证码：", code)
	//将验证码存入redis
	utils.Red.Set(c, phone, code, 5*time.Minute)
	c.JSON(http.StatusOK, gin.H{
		"code":    0, //成功
		"message": strconv.Itoa(code),
		"data":    nil,
	})
	return
}

func LoginByCode(c *gin.Context) {
	//post请求->phone
	phone := c.Request.FormValue("phone")
	if !isMatchPhone(phone) {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": "手机号码无效",
			"data":    nil,
		})
		return
	}
	code := c.Request.FormValue("code")
	if code == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": "请输入验证码！",
			"data":    nil,
		})
		return
	}
	//查询redis
	cacheCode, _ := utils.Red.Get(c, phone).Result()
	//不一致就不放行
	if code != cacheCode {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": "验证码错误",
			"data":    nil,
		})
		return
	}
	//一致就放行->如果用户尚且未注册，直接可以注册并告知默认密码
	user := models.FindUserByPhone(phone)
	if user.Password == "" {
		user.Phone = phone
		user.Password, _ = utils.GetPwd("111111Az*")
		models.CreateUser(user)
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "已自动帮您注册，默认密码为111111Az*",
			"data":    user,
		})
		return
	}
	if !signed(user, c) {
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "欢迎回来",
		"data":    user,
	})
	return
}

func ResetPassword(c *gin.Context) {
	phone := c.Request.FormValue("phone")
	password := c.Request.FormValue("password")
	user := models.FindUserByPhone(phone)
	if !isMatchPhone(phone) {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": "手机号码无效",
			"data":    nil,
		})
		return
	}
	if user.Password == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": "用户尚未注册",
			"data":    nil,
		})
		return
	}
	if !isStrongPassword(password) {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": "密码无效",
			"data":    nil,
		})
		return
	}
	models.EditUserPassword(password, phone)
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "重置成功",
		"data":    nil,
	})
}

// 升级-->管理端
func Admin(c *gin.Context) {
	key := c.Request.FormValue("key")
	if key != viper.GetString("root.key") {
		utils.RespFail(c.Writer, "密钥错误")
		return
	}
	id_ := c.Request.FormValue("id")
	user := models.FindUserByUserInfoId(id_).UserInfo
	user.Flag = 1
	utils.DB.Exec("Update user_info set flag=1 where id=?", id_)
	fmt.Println(user.Flag)
	utils.RespOk(c.Writer, user, "升级成功，您可以执行管理员的任务了")
	return
}

func Profile(c *gin.Context) {
	n := c.Request.FormValue("num")
	userId := c.Request.FormValue("user_id")
	userInfo := models.FindUserInfo(userId)
	nums := strings.Split(n, " ")
	for _, num := range nums {
		switch num {
		case "1":
			name := c.Request.FormValue("name")
			userInfo.Name = name
		case "2":
			userInfo.ProfilePhoto, _ = upload(c.Request, c.Writer, c)
		case "3":
			p := c.PostForm("birthday")
			p_, _ := strconv.Atoi(p)
			time := time.Unix(int64(p_), 0)
			userInfo.Birthday = time
		case "4":
			interest := c.PostFormArray("interest")
			interest_, _ := json.Marshal(interest)
			userInfo.Interest = string(interest_)
		case "5":
			userInfo.Sign = c.Request.FormValue("sign")
		default:
			utils.RespFail(c.Writer, "输入不规范！")
			return
		}

	}
	userInfo.RefleshUserInfo_()
	//if err != nil {
	//	utils.RespFail(c.Writer, "修改失败："+err.Error())
	//	return
	//}
	utils.RespOk(c.Writer, userInfo, "修改成功")
	return
}

func UserDetail(c *gin.Context) {
	id := c.Query("user_id")
	fmt.Println(id)
	//todo 兴趣和生日和签名，没有
	user := models.FindUserInfo(id)
	utils.RespOk(c.Writer, user, user.Name+"个人信息图下")
}

type res struct {
	Name    string
	Picture string
	Info    string
	Price   float64 //价格
	Num     int
}
type resTicket struct {
	Name string
	Num  int
	//演出厅
	Place int
	//影片开始结束时间
	Begintime time.Time
	Endtime   time.Time
	//座位
	Seat []models.Seat `json:"seat"`
}

func MyOrder(c *gin.Context) {
	id := c.Query("user_id")
	user := models.FindUserInfo(id)
	ticker := user.Ticket
	snack := user.Snack

	snacks := strings.Split(snack, " ")
	tickers := strings.Split(ticker, " ")
	Tickets := make([]resTicket, 0)

	if ticker != "" {
		for _, ticketID := range tickers {
			t := models.GetTicketByID(ticketID)
			seat := []models.Seat{}
			err := json.Unmarshal(t.Seat, &seat)
			if err != nil {
				utils.RespFail(c.Writer, err.Error())
				return
			}
			tt := resTicket{
				Name:      t.Name,
				Num:       t.Num,
				Place:     t.Place,
				Begintime: t.Begintime,
				Endtime:   t.Endtime,
				Seat:      seat,
			}
			Tickets = append(Tickets, tt)
		}
	} else {
		Tickets = nil
	}

	Snacks := make([]*models.Snack, 0)
	ress := make([]res, 0)
	if snack != "" {
		for _, snackId := range snacks {
			s := models.GetsnackByid(snackId)
			Snacks = append(Snacks, &s)
		}
		SnackInfo := mergeSnacks(Snacks)
		for k, v := range SnackInfo {
			snack_ := models.GetsnackByid(k)
			ress_ := &res{
				Name:    snack_.Name,
				Picture: snack_.Picture,
				Info:    snack_.Info,
				Price:   snack_.Price,
				Num:     v,
			}
			ress = append(ress, *ress_)
		}
	} else {
		ress = nil
	}
	utils.RespOk(c.Writer, gin.H{
		"ticket": Tickets,
		"snack":  ress,
	}, "返回订单信息")
}
