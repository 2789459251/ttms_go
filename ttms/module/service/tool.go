package service

import (
	models2 "TTMS_go/ttms/models"
	utils "TTMS_go/ttms/util"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/spf13/viper"
	"net/http"
	"regexp"
	"strconv"
)

// const url_ = "http://sb1cf9mjk.hb-bkt.clouddn.com/"
const url_ = "http://video.cdn.zy520.online/"

func isMatchPhone(phone string) bool {
	flag, _ := regexp.Match("^1[3-9]{1}\\d{9}", []byte(phone))
	if len(phone) != 11 {
		flag = false
	}
	return flag
}

func isStrongPassword(password string) bool {
	// 密码长度在8到16之间
	if len(password) < 8 || len(password) > 16 {
		return false
	}

	hasUpperCase := false
	hasLowerCase := false
	hasDigit := false
	hasSpecialChar := false

	for _, char := range password {
		ascii := int(char)

		// 检查大写字母
		if ascii >= 65 && ascii <= 90 {
			hasUpperCase = true
		}

		// 检查小写字母
		if ascii >= 97 && ascii <= 122 {
			hasLowerCase = true
		}

		// 检查数字
		if ascii >= 48 && ascii <= 57 {
			hasDigit = true
		}

		// 检查特殊字符
		if (ascii >= 33 && ascii <= 47) || (ascii >= 58 && ascii <= 64) || (ascii >= 91 && ascii <= 96) || (ascii >= 123 && ascii <= 126) {
			hasSpecialChar = true
		}
	}

	// 检查是否满足所有条件
	return hasUpperCase && hasLowerCase && hasDigit && hasSpecialChar
}

func signed(user models2.User, c *gin.Context) bool {
	jwt, err := utils.InitAuth()
	if err != nil {
		return false
	}
	id := strconv.Itoa(user.UserInfoId)
	rT, aT, _ := jwt.GenerateTokens(id)
	c.Header("Authorization", "Bearer "+aT)
	c.SetCookie("refresh_token", rT, 3600, "/", "localhost", false, true)
	return true
}

func User(c *gin.Context) models2.UserInfo {
	userinfoid, _ := c.Get("userInfo")
	userinfo := models2.FindUserInfo(userinfoid.(string))
	return userinfo
}

func upload(r *http.Request, w http.ResponseWriter, c *gin.Context) (string, error) {
	putPolicy := storage.PutPolicy{Scope: viper.GetString("qiniu.Scope")}
	mac := qbox.NewMac(viper.GetString("qiniu.QiniuAK"), viper.GetString("qiniu.QiniuSK"))
	upTocken := putPolicy.UploadToken(mac)

	cfg := storage.Config{Zone: &storage.ZoneHuabei, UseHTTPS: false, UseCdnDomains: false}
	file, head, err := r.FormFile("picture")
	if err != nil {
		utils.RespFail(c.Writer, "文件读取失败："+err.Error())
	}
	//fmt.Println(head.Header)
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	putExtra := storage.PutExtra{}
	fmt.Println(head.Filename)
	err = formUploader.Put(context.Background(), &ret, upTocken, head.Filename, file, head.Size, &putExtra)
	return url_ + head.Filename, err
}

func isLimited(c *gin.Context) bool {
	user := User(c)
	if user.Flag == 0 {
		utils.RespFail(c.Writer, "权限不够")
		return false
	}
	return true
}

func mergeSnacks(snacks []*models2.Snack) map[string]int {
	// 使用映射来存储每个零食的名字和合并后的库存量
	snackMap := make(map[string]int)

	for _, snack := range snacks {
		id := strconv.Itoa(int(snack.ID))
		if _, exists := snackMap[id]; exists {
			snackMap[id] = snackMap[id] + 1
		} else {
			// 否则，将这个零食添加到映射中
			snackMap[id] = 1
		}
	}

	return snackMap
}
