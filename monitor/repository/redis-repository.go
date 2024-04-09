package repository

import (
	"context"
	"os"
	"sync"

	"github.com/redis/go-redis/v9"
)

type RedisRepository struct {
	rh    *redis.Client
	ctx   *context.Context
	mutex *sync.Mutex
}

var keyPrefixes = struct {
	Pages string
}{
	Pages: "pages",
}

func (r *RedisRepository) GetPagesCount() int64 {
	return r.rh.HLen(*r.ctx, keyPrefixes.Pages).Val()
}

func NewRedisRepository() *RedisRepository {
	connectionAddress := os.Getenv("REDIS_CONNECTION_ADDRESS")
	password := os.Getenv("REDIS_PASSWORD")

	rdb := redis.NewClient(&redis.Options{
		Addr:     connectionAddress,
		Password: password,
		DB:       0,
	})
	ctx := context.Background()
	mu := new(sync.Mutex)
	return &RedisRepository{rdb, &ctx, mu}
}
