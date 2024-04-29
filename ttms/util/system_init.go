package utils

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/olivere/elastic/v7"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var (
	DB  *gorm.DB
	Red *redis.Client
	ES  *elastic.Client
)

func EsClient() {
	client, err := elastic.NewClient(
		elastic.SetURL(viper.GetString("es.url")),
		elastic.SetBasicAuth(viper.GetString("es.basicauth.username"), viper.GetString("es.basicauth.password")))
	elastic.SetSniff(viper.GetBool("es.sniff"))
	if err != nil {
		fmt.Println("连接失败：", err.Error())
		return
	}
	ES = client
	fmt.Println("连接Es客户端：", ES)
}

func InitConfig() {
	viper.SetConfigName("app")
	viper.AddConfigPath("conf")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("config app inited")
}
func InitMysql() {
	//自定义日志模版,打印sql语句
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)
	fmt.Println(viper.GetString("mysql.dns"))
	DB, _ = gorm.Open(mysql.Open(viper.GetString("mysql.dns")),
		&gorm.Config{Logger: newLogger})

}
func InitRedis() {
	Red = redis.NewClient(&redis.Options{
		Addr:         viper.GetString("redis.addr"),
		Password:     viper.GetString("redis.password"),
		DB:           viper.GetInt("redis.DB"),
		PoolSize:     viper.GetInt("redis.poolSize"),
		MinIdleConns: viper.GetInt("redis.minIdleConn"),
	})
}
