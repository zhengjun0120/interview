package database

import (
	"ai_interview/conf"
	"ai_interview/pkg/zlog"
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
	zlog.Infof("mysql连接成功")
}

func GetDB() *gorm.DB {
	return db
}
