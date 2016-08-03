package redis

import (
	"sync"
	"time"

	"brooce/config"

	redis "gopkg.in/redis.v3"
)

var redisClient *redis.Client
var once sync.Once
var threads = config.Config.TotalThreads() + 10

func Get() *redis.Client {
	once.Do(func() {
		redisClient = redis.NewClient(&redis.Options{
			Addr:         config.Config.Redis.Host,
			Password:     config.Config.Redis.Password,
			MaxRetries:   10,
			PoolSize:     threads,
			DialTimeout:  5 * time.Second,
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 5 * time.Second,
			PoolTimeout:  1 * time.Second,
		})
	})

	return redisClient
}

func FlushList(src, dst string) (err error) {
	redisClient := Get()
	for err == nil {
		_, err = redisClient.RPopLPush(src, dst).Result()
	}

	if err == redis.Nil {
		err = nil
	}

	return
}
