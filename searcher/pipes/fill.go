package pipes

import (
	"diggle/searcher/models"
	"diggle/searcher/repository"
)

func Fill(input *[]models.RawPage, size, page int, repo *repository.RedisRepository) *[]models.Page {

	start := (page - 1) * size
	end := page * size
	if end > len(*input) {
		end = len(*input)
	}

	pages := make([]models.Page, 0, end-start)

	for i := start; i < end; i++ {
		page := repo.GetPage((*input)[i].PageId)
		pages = append(pages, *page)
	}

	return &pages
}
