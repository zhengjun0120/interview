package database

import (
	"ai_interview/conf"
	"ai_interview/pkg/zlog"
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var client *redis.Client

func InitRedis() {
	host := conf.GetConfig().Redis.Host
	port := conf.GetConfig().Redis.Port
	password := conf.GetConfig().Redis.Password
	dbn := conf.GetConfig().Redis.DB

	_client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
		DB:       dbn,
	})

	if err := _client.Ping(context.Background()).Err(); err != nil {
		panic("redis连接失败" + err.Error())
	}
	client = _client
	zlog.Infof("redis连接成功")
}

func GetRedis() *redis.Client {
	return client
}
