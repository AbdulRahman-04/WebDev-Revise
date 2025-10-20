package utils

import (
	"context"
	"fmt"

	"github.com/AbdulRahman-04/07EvenetManagement-API/server/config"
	"github.com/redis/go-redis/v9"
)

var (
	RedisClient*redis.Client
	ctx = context.Background()
)

func RedisConnect() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr: config.AppConfig.Redis.Host,
		Password: config.AppConfig.Redis.Password,
		DB: config.AppConfig.Redis.DB,
	})

	_ ,err := RedisClient.Ping(ctx).Result()
	if err != nil {
		fmt.Println("Ping to redis failed")
		return
	}

	fmt.Println("Redis Connected✅")
}

func RedisSet(key string , value string) error {
	return RedisClient.Set(ctx, key, value, 0).Err()
}

func RedisGet(key string) (string, error){
	return  RedisClient.Get(ctx, key).Result()
}

func RedisDel(key string) error {
	return RedisClient.Del(ctx, key).Err()
}