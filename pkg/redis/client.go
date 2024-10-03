package redis

import goRedis "github.com/redis/go-redis/v9"

func NewRedisClient(host string, password string) *goRedis.Client {
	client := goRedis.NewClient(&goRedis.Options{
		Addr:     host,
		Password: password,
		DB:       0,
	})
	return client
}
