package Db

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"new_url/config"
	"time"
)

func NewRedis() *redis.Client {
	Rdb := redis.NewClient(&redis.Options{
		Addr:     config.GlobalCfg.Redis.RdbPort,
		Password: config.GlobalCfg.Redis.Password,
		DB:       config.GlobalCfg.Redis.Rdb,
	})

	// 2. 必须验证连接（关键！）
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err := Rdb.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("Redis 连接失败：%v", err)) // 启动时就报错，避免后续问题
	}

	return Rdb
}
