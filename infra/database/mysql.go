package database

import (
	"ai_interview/conf"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDB() {
	dsn := conf.GetConfig().Mysql.Dsn
	_db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		panic("mysql连接失败" + err.Error())
	}
	db = _db
	fmt.Println("mysql连接成功")
}

func GetDB() *gorm.DB {
	return db
}
