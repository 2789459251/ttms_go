package main

import (
	router "TTMS_go/ttms/app/api"
	"TTMS_go/ttms/domain/models"
	dto "TTMS_go/ttms/domain/models/dao"
	utils "TTMS_go/ttms/util"
)

func main() {
	utils.InitMysql()
	utils.InitRedis()
	utils.DB.AutoMigrate(models.User{})
	utils.DB.AutoMigrate(dto.UserInfo{})
	utils.DB.AutoMigrate(models.Ticket{})
	utils.DB.AutoMigrate(models.Place{})
	utils.DB.AutoMigrate(models.Movie{})
	utils.DB.AutoMigrate(models.Snack{})
	r := router.Router()
	r.Run("0.0.0.0:8082")
}
