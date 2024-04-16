package main

import (
	models2 "TTMS_go/ttms/models"
	utils "TTMS_go/ttms/util"
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var DB *gorm.DB

func main_() {
	viper.SetConfigName("app")
	viper.AddConfigPath("config")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("config app inited")
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)

	DB, _ = gorm.Open(mysql.Open(viper.GetString("mysql.dns")),
		&gorm.Config{Logger: newLogger})
	utils.DB.AutoMigrate(models2.User{})
	utils.DB.AutoMigrate(models2.UserInfo{})
	utils.DB.AutoMigrate(models2.Ticket{})
	utils.DB.AutoMigrate(models2.Theatre{})
	utils.DB.AutoMigrate(models2.Movie{})
	utils.DB.AutoMigrate(models2.Snack{})
}
