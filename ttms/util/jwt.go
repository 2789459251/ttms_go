package utils

import (
	"TTMS_go/ttms/models"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/locales/vo"
	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
	"strings"
	"time"
)

type JWTMiddleware struct {
	AccessSecret  []byte
	RefreshSecret []byte
	Timeout       int
	MaxRefresh    int
}

type Claims struct {
	User models.UserInfo
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

func (jm *JWTMiddleware) GenerateToken(user models.UserInfo) (string, string, int64) {
	aT := Claims{user, jwt.StandardClaims{
		Issuer:    "Zy",
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Hour * time.Duration(jm.Timeout)).Unix(),
	},
	}
	rT := Claims{user, jwt.StandardClaims{
		Issuer:    "Zy",
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Hour * time.Duration(jm.MaxRefresh)).Unix(),
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
	return assessTokenSigned, refreshTokenSigned, time.Now().Unix()
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
		fmt.Errorf("解析refresh token失败：%s", err.Error())
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
		if isUpd {
			c.Abort()
			RespFail(c.Writer, "令牌失效")
			return
		}
		c.Set("userInfo", parseToken.User)
		c.Next()
	}
}
func (jm *JWTMiddleware) RefreshHandler(c *gin.Context) {
	var req vo.RefreshTokenRequest
	//请求json绑定
	if err := c.ShouldBind(&req); err != nil {
		fmt.Errorf("获取失败,%s", err.Error())
		RespFail(c.Writer, err.Error())
		return
	}

}
