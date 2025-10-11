package redis

import (
	"fmt"

	"github.com/Romasmi/advanced-url-shortener/internal/config"
	"github.com/redis/go-redis/v9"
)

type RedisConnection struct {
	Rdb    *redis.Client
	Config *config.Config
}

func (c *RedisConnection) Connect() {
	c.Rdb = redis.NewClient(&redis.Options{
		Addr:     c.Config.Redis.Host,
		Username: c.Config.Redis.Username,
		Password: c.Config.Redis.Password,
		DB:       0,
	})
}

func (c *RedisConnection) Close() {
	if c.Rdb != nil {
		err := c.Rdb.Close()
		if err != nil {
			fmt.Printf("Erro while closing Redis connection %v\n", err)
		}
	}
}
