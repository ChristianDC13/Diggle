package main

import (
	"diggle/ranker/ranker"
	"diggle/ranker/repository"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	repo := repository.NewRedisRepository()
	r := ranker.NewRanker(repo)
	err := r.Rank()
	if err != nil {
		panic(err)
	}
	log.Println("Ranking completed")

}
