package data

import (
	"log"
	"time"
	"yanglu/config"

	"github.com/patrickmn/go-cache"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var mainDb *gorm.DB

func InitMysql() {
	var err error

	mainDb, err = gorm.Open(mysql.Open(config.GetMysqlApi()), &gorm.Config{
		Logger: nil,
	})
	if err != nil {
		log.Fatalf("failed to connect database, dns: %v, err: %v", config.GetMysqlApi(), err)
	}

	sqlDB, err := mainDb.DB()
	if err != nil {
		log.Fatalf("failed to get sql db, dns: %v, err: %v", config.GetMysqlApi(), err)
	}

	sqlDB.SetMaxIdleConns(100)
	sqlDB.SetMaxOpenConns(1000)

}

func GetDB() *gorm.DB {
	return mainDb
}

//内存缓存
var C *cache.Cache

func InitMemoryCache() {
	C = cache.New(time.Second, time.Minute)
}
