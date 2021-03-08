package config

import (
	"fmt"
	"github.com/go-redis/redis"
)

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

func NewRedisClient(conf *RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		Password: conf.Password,
		DB:       conf.DB,
	})
	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}
