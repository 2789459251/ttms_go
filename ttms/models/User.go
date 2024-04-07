package models

import (
	utils "TTMS_go/ttms/util"
	"gorm.io/gorm"
)

// validate:"min=6,max=10" `valid:"matches(^1[3-9]{1}\\d{9})"`
// `valid:"matches(^1[3-9]{1}\\d{9}$)"`
// `valid:"matches(^(?=.*[A-Za-z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{8,}$)"`
type User struct {
	gorm.Model
	Phone      string   `validate:"required"`
	Password   string   `validate:"required"`
	UserInfoId int      `json:"user_info_id"`
	UserInfo   UserInfo `gorm:"foreignKey:user_info_id"`
}

func (table *User) TableName() string {
	return "user_basic"
}
func CreateUser(user User) *gorm.DB {
	userInfo := UserInfo{}
	user.UserInfoId = int(user.ID)
	//Todo: 这里的user.userInfoId应该让userinfo创建结束后赋值,但是也必须考虑token中userInfoId
	utils.DB.Create(&userInfo)
	user.UserInfo = userInfo
	return utils.DB.Create(&user)
}
func FindUserByPhone(phone string) User {
	user := User{}
	utils.DB.Where("phone=?", phone).First(&user)
	return user
}
func FindUserById(id string) User {
	user := User{}
	utils.DB.Where("id=?", id).First(&user)
	return user
}
func EditUserPassword(password, phone string) {
	user := FindUserByPhone(phone)
	user.Password, _ = utils.GetPwd(password)
	utils.DB.Updates(&user)
}
func FindUserByUserInfoId(id string) User {
	user := User{}
	utils.DB.Where("user_info_id?", id).First(&user)
	return user
}
func RefreshUserInfo(id string, userInfo UserInfo) {
	user := FindUserByUserInfoId(id)
	user.UserInfo = userInfo
	utils.DB.Updates(&user)
}
