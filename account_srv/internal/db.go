package internal

import (
	"fmt"
	"github.com/i-coder-robot/mic-trainning-lessons/account_srv/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"
)

var DB *gorm.DB
var err error

func InitDB() {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		logger.Config{
			SlowThreshold:             time.Second, // 慢 SQL 阈值
			LogLevel:                  logger.Info, // 日志级别
			IgnoreRecordNotFoundError: true,        // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  true,        // 禁用彩色打印
		},
	)
	//启动mysql,到Github/mic-basic-lessons/week4/docker-compose up -d
	conn := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", "root", "123456", "192.168.0.106", 3308, "happy_account_mic_traning")
	DB, err = gorm.Open(mysql.Open(conn), &gorm.Config{
		Logger: newLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, //表明用英文单数形式
		},
	})
	if err != nil {
		panic("数据库连接失败:" + err.Error())
	}
	err = DB.AutoMigrate(&model.Account{})
	if err != nil {
		fmt.Println(err)
	}
}
