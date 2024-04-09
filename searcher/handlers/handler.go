package handlers

import (
	"diggle/searcher/repository"
	"diggle/searcher/services"
)

type Handler struct {
	searchService *services.SearchService
	repo          *repository.RedisRepository
}

func NewHandler(
	searchService *services.SearchService,
	repo *repository.RedisRepository,

) *Handler {
	return &Handler{
		searchService: searchService,
		repo:          repo,
	}
}
