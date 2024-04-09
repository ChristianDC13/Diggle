package repository

import (
	"context"
	"diggle/searcher/models"
	"encoding/json"
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
	Ranks           string
	WordFrequencies string
}{
	Word:            "w:",
	Pages:           "pages",
	PageToId:        "p-i:",
	IdToPage:        "i-p:",
	OutboundLinks:   "out:",
	Ranks:           "ranks",
	WordFrequencies: "words-freq",
}

func (r *RedisRepository) GetPagesForWord(word string) map[string]string {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	res, err := r.rh.HGetAll(*r.ctx, keyPrefixes.Word+word).Result()
	if err != nil {
		return nil
	}
	return res
}

func (r *RedisRepository) GetPage(id int64) *models.Page {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	key := strconv.FormatInt(id, 10)
	res, err := r.rh.HGet(*r.ctx, keyPrefixes.Pages, key).Result()

	var page models.Page

	json.Unmarshal([]byte(res), &page)

	if err != nil {
		return nil
	}
	return &page
}

func (r *RedisRepository) GetRank(id int64) float64 {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	key := strconv.FormatInt(id, 10)
	res, err := r.rh.HGet(*r.ctx, keyPrefixes.Ranks, key).Result()

	if err != nil {
		return 0
	}

	rank, _ := strconv.ParseFloat(res, 64)
	return rank
}

func (r *RedisRepository) GetRanks(pageIds []string) map[string]float64 {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	res, err := r.rh.HMGet(*r.ctx, keyPrefixes.Ranks, pageIds...).Result()

	if err != nil {
		return nil
	}

	ranks := map[string]float64{}

	for i, rank := range res {
		if rank != nil {
			rankFloat, _ := strconv.ParseFloat(rank.(string), 64)
			ranks[pageIds[i]] = rankFloat
		}
	}

	return ranks
}

func (r *RedisRepository) GetAllWords() ([]models.WordFrequency, error) {

	res, err := r.rh.HGetAll(*r.ctx, keyPrefixes.WordFrequencies).Result()
	if err != nil {
		return nil, err
	}
	frequencies := []models.WordFrequency{}

	for word, count := range res {
		countInt, _ := strconv.Atoi(count)

		frequencies = append(frequencies, models.WordFrequency{
			Word:  word,
			Count: countInt,
		})
	}
	return frequencies, nil
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
	mu := &sync.Mutex{}
	return &RedisRepository{rdb, &ctx, mu}
}
