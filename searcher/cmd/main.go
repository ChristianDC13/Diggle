package main

import (
	fuzzysearch "diggle/searcher/fuzzy-search"
	"diggle/searcher/handlers"
	"diggle/searcher/repository"
	"diggle/searcher/server"
	"diggle/searcher/services"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	repo := repository.NewRedisRepository()

	fuzzySearcher, err := fuzzysearch.NewFuzzySearcher(repo, 2)
	if err != nil {
		panic(err)
	}
	searchService := services.NewSearchService(repo, fuzzySearcher)
	handler := handlers.NewHandler(searchService, repo)
	server := server.NewServer(handler)

	server.Start()
}
