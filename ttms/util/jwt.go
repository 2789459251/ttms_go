package utils

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
	"strings"
	"time"
)

type user string
type JWTMiddleware struct {
	AccessSecret  []byte
	RefreshSecret []byte
	Timeout       int
	MaxRefresh    int
}

type Claims struct {
	User string `json:"user"`
	jwt.StandardClaims
}

func InitAuth() (*JWTMiddleware, error) {
	authMiddleware := &JWTMiddleware{
		AccessSecret:  []byte(viper.GetString("jwt.AKey")),
		RefreshSecret: []byte(viper.GetString("jwt.RKey")),
		Timeout:       viper.GetInt("jwt.Timeout"),
		MaxRefresh:    viper.GetInt("jwt.MaxRefresh"),
	}
	return authMiddleware, nil
}

func (jm *JWTMiddleware) GenerateTokens(user string) (string, string, int64) {
	fmt.Println("ggg:", user)
	fmt.Println(jm.Timeout)
	aT := Claims{user, jwt.StandardClaims{
		Issuer:    "Zy",
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Second * time.Duration(jm.Timeout)).Unix(),
	},
	}
	rT := Claims{user, jwt.StandardClaims{
		Issuer:    "Zy",
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Second * time.Duration(jm.MaxRefresh)).Unix(),
	},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, aT)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, rT)
	assessTokenSigned, err := accessToken.SignedString(jm.AccessSecret)
	if err != nil {
		fmt.Errorf("获取token失败，Secret错误")
		return "", "", 0
	}
	refreshTokenSigned, err := refreshToken.SignedString(jm.RefreshSecret)
	if err != nil {
		fmt.Errorf("获取token失败，Secret错误")
		return "", "", 0
	}

	Red.Set(context.Background(), refreshTokenSigned, true, time.Duration(rT.ExpiresAt))
	return assessTokenSigned, refreshTokenSigned, time.Now().Unix()
}
func (jm *JWTMiddleware) GetAccessToken(user string) (string, error) {
	aT := Claims{user, jwt.StandardClaims{
		Issuer:    "Zy",
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Second * time.Duration(jm.Timeout)).Unix(),
	}}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, aT)
	accessTokenSigned, err := accessToken.SignedString(jm.AccessSecret)
	return accessTokenSigned, err
}
func (jm *JWTMiddleware) ParseRefreshToken(refreshTokenString string) (*Claims, bool, error) {
	refreshToken, err := jwt.ParseWithClaims(refreshTokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) { return jm.RefreshSecret, nil })
	if err != nil {
		fmt.Errorf("解析refresh token失败：%s", err.Error())
		return nil, false, err
	}
	if claims, ok := refreshToken.Claims.(*Claims); ok && refreshToken.Valid {
		return claims, true, nil
	}
	return nil, false, errors.New("invaild refresh token")
}
func (jm *JWTMiddleware) ParseAccessToken(accessTokenString string) (*Claims, bool, error) {
	accessToken, err := jwt.ParseWithClaims(accessTokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) { return jm.RefreshSecret, nil })
	if err != nil {
		fmt.Errorf("解析access token失败：%s", err.Error())
		return nil, false, err
	}
	if claims, ok := accessToken.Claims.(*Claims); ok && accessToken.Valid {
		return claims, true, nil
	}
	return nil, false, errors.New("invaild  accesstoken")
}

func (jm *JWTMiddleware) JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.Abort()
			RespFail(c.Writer, "传入空令牌")
			return
		}
		parts := strings.Split(authHeader, " ")
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.Abort()
			RespFail(c.Writer, "令牌不规范")
			return
		}
		parseToken, isUpd, err := jm.ParseAccessToken(parts[1])
		if err != nil {
			fmt.Errorf("")
			c.Abort()
			RespFail(c.Writer, "令牌解析错误")
			return
		}
		if !isUpd {
			c.Abort()
			fmt.Println(jm)
			RespFail(c.Writer, "access令牌失效")
			return
		}
		fmt.Println("ppp:", parseToken.User)
		c.Set("userInfo", parseToken.User)
		c.Next()
	}
}
func (jm *JWTMiddleware) RefreshHandler(c *gin.Context) {
	refreshStr, err := c.Cookie("refresh_token")
	if refreshStr == "" || err != nil {
		RespFail(c.Writer, "令牌有误refresh失败")
		return
	}

	refresh, isUpd, err := jm.ParseRefreshToken(refreshStr)

	if !isUpd || err != nil {
		RespFail(c.Writer, "第二重验证失败，refresh_token失效或者解析出错")
		return
	}

	flag, err := Red.Get(context.Background(), refreshStr).Result()
	if flag == "" || err != nil {
		RespFail(c.Writer, "第三重验证失败，refresh_token已移除，请重新登录")
		return
	}

	access, err := jm.GetAccessToken(refresh.User)
	if err != nil {
		RespFail(c.Writer, "生成新access失败："+err.Error())
		return
	}
	RespOk(c.Writer, access, "刷新成功")
}
