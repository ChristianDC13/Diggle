package repository

import (
	"context"
	"diggle/crawler/webpage"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/redis/go-redis/v9"
)

type RedisRepository struct {
	rh    *redis.Client
	ctx   *context.Context
	mutex *sync.Mutex
}

var keyPrefixes = struct {
	Word            string
	Pages           string
	PageToId        string
	IdToPage        string
	OutboundLinks   string
	PageQueue       string
	VisitedPages    string
	WordFrequencies string
}{
	Word:            "w:",
	Pages:           "pages",
	PageToId:        "p-i:",
	IdToPage:        "i-p:",
	OutboundLinks:   "out:",
	PageQueue:       "queue",
	VisitedPages:    "visited",
	WordFrequencies: "words-freq",
}

func (r *RedisRepository) AddPage(page *webpage.WebPage) (bool, error) {

	r.mutex.Lock()
	defer r.mutex.Unlock()

	exists, _ := r.getPageId(page.URL)
	if exists {
		return false, nil
	}

	pageId, err := r.savePageData(page)
	if err != nil {
		return false, err
	}

	err = r.savePageLinks(page)
	if err != nil {
		return false, err
	}

	for _, word := range page.Words {
		key := keyPrefixes.Word + word.Word
		scorePosition := fmt.Sprint(word.Score, ":", word.Position)
		err := r.rh.HSet(*r.ctx, key, pageId, scorePosition).Err()

		if err != nil {
			return false, err
		}

		err = r.rh.HIncrBy(*r.ctx, keyPrefixes.WordFrequencies, word.Word, 1).Err()

		if err != nil {
			return false, err
		}
	}
	return true, nil
}

func (r *RedisRepository) GetPagesCount() int64 {
	return r.rh.HLen(*r.ctx, keyPrefixes.Pages).Val()
}

func (r *RedisRepository) savePageData(page *webpage.WebPage) (int64, error) {

	pageId := r.rh.HLen(*r.ctx, keyPrefixes.Pages).Val()
	info := map[string]interface{}{
		"url":      page.URL,
		"title":    page.Title,
		"abstract": page.Abstract,
	}

	jsonData, err := json.Marshal(info)
	if err != nil {
		return pageId, err
	}

	err = r.rh.HSet(*r.ctx, keyPrefixes.Pages, pageId, string(jsonData)).Err()
	if err != nil {
		return pageId, err
	}

	return pageId, nil
}

func (r *RedisRepository) savePageLinks(page *webpage.WebPage) error {

	pageExists := r.rh.HExists(*r.ctx, keyPrefixes.Pages, page.URL).Val()
	if pageExists {
		return nil
	}

	_, pageId := r.getPageId(page.URL)
	key := fmt.Sprint(keyPrefixes.OutboundLinks, pageId)
	for _, link := range page.OutboundLinks {
		_, outLink := r.getPageId(link)
		r.rh.HSet(*r.ctx, key, outLink, outLink).Err()
	}
	return nil
}

func (r *RedisRepository) getPageId(url string) (bool, int64) {
	pageExists := r.rh.HExists(*r.ctx, keyPrefixes.PageToId, url).Val()
	if pageExists {
		res := r.rh.HGet(*r.ctx, keyPrefixes.PageToId, url).Val()
		pageId, _ := strconv.ParseInt(res, 10, 64)
		return true, pageId
	}

	pageId := r.rh.HLen(*r.ctx, keyPrefixes.Pages).Val()
	r.rh.HSet(*r.ctx, keyPrefixes.PageToId, url, pageId)

	r.rh.HSet(*r.ctx, keyPrefixes.IdToPage, pageId, url)
	return false, pageId
}

func (r *RedisRepository) EnqueuePage(url string) (bool, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	alreadyVisited := r.rh.HExists(*r.ctx, keyPrefixes.VisitedPages, url).Val()
	if alreadyVisited {
		return false, nil
	}

	err := r.rh.LPush(*r.ctx, keyPrefixes.PageQueue, url).Err()
	if err != nil {
		return false, err
	}

	err = r.rh.HSet(*r.ctx, keyPrefixes.VisitedPages, url, true).Err()

	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *RedisRepository) DequeuePage() (string, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	url, err := r.rh.RPop(*r.ctx, keyPrefixes.PageQueue).Result()
	if err != nil {
		return "", err
	}

	return url, nil
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
