package repository

import (
	"context"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

type RedisRepository struct {
	rh  *redis.Client
	ctx *context.Context
}

var keyPrefixes = struct {
	Pages         string
	OutboundLinks string
	PageRank      string
}{
	Pages:         "pages",
	OutboundLinks: "out:",
	PageRank:      "ranks",
}

func (r *RedisRepository) GetPages() ([]string, error) {
	pages := r.rh.HKeys(*r.ctx, keyPrefixes.Pages).Val()
	return pages, nil
}

func (r *RedisRepository) GetOutboundLinks(pageId string) ([]string, error) {
	outboundKey := fmt.Sprint(keyPrefixes.OutboundLinks, pageId)
	outboundLinks := r.rh.HKeys(*r.ctx, outboundKey).Val()
	return outboundLinks, nil
}

func (r *RedisRepository) SetPageRank(pageId string, rank float64) error {
	err := r.rh.HSet(*r.ctx, keyPrefixes.PageRank, pageId, rank).Err()
	return err
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
	return &RedisRepository{rdb, &ctx}
}
