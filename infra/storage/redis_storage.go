package storage

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStorage struct {
	client *redis.Client
}

func (rds *RedisStorage) Get(c context.Context, key string) (interface{}, error) {
	return rds.client.Get(c, key).Result()
}

func (rds *RedisStorage) Set(c context.Context, key string, value any, expiration time.Duration) error {
	return rds.client.Set(c, key, value, expiration).Err()
}

func (rds *RedisStorage) HGet(c context.Context, key string, receiver any) error {
	val, err := rds.client.Get(c, key).Bytes()
	if err != nil {
		return err
	}
	err = json.Unmarshal(val, receiver)
	if err != nil {
		return err
	}
	return nil
}

func (rds *RedisStorage) HSet(c context.Context, key string, value any, expiration time.Duration) error {
	val, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return rds.client.Set(c, key, val, expiration).Err()
}

func (rds *RedisStorage) Del(c context.Context, key string) error {
	_, err := rds.client.Del(c, key).Result()
	if err != nil {
		return err
	}
	return nil
}
