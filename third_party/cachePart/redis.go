package cachePart

import (
	"context"
	"github.com/redis/go-redis/v9"
)

type Params struct {
	Host     string
	Port     string
	Password string
	DB       int
}

// NewRedis - establishing cachePart client connection with cachePart server
func NewRedis(cfg Params) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Host + ":" + cfg.Port,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	return client, nil
}
