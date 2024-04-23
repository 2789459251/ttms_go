package main

import (
	router "TTMS_go/ttms/app/api"
	utils "TTMS_go/ttms/util"
)

func main() {
	utils.InitConfig()
	utils.InitMysql()
	utils.InitRedis()
	//utils.DB.Migrator().DropTable(models2.User{})
	//utils.DB.Migrator().DropTable(models2.UserInfo{})
	//
	//utils.DB.AutoMigrate(models2.User{})
	//utils.DB.AutoMigrate(models2.UserInfo{})
	//utils.DB.AutoMigrate(models2.Ticket{})
	//
	//utils.DB.Migrator().DropTable(models.Theatre{})
	//utils.DB.AutoMigrate(models.Theatre{})
	//
	//utils.DB.Migrator().DropTable(models.Play{})
	//utils.DB.AutoMigrate(models.Play{})

	//utils.DB.Migrator().DropTable(models2.Movie{})
	//utils.DB.AutoMigrate(models2.Movie{})
	//utils.DB.Migrator().DropTable(models2.Snack{})
	//utils.DB.AutoMigrate(models2.Snack{})
	r := router.Router()
	r.Run("0.0.0.0:8082")
}
