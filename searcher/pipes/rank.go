package pipes

import (
	"diggle/searcher/models"
	"diggle/searcher/repository"
	"fmt"
)

func Rank(input *[]models.RawPage, repo *repository.RedisRepository) *[]models.RawPage {

	ids := make([]string, 0, len(*input))
	for _, page := range *input {
		idStr := fmt.Sprintf("%d", page.PageId)
		ids = append(ids, idStr)
	}

	ranks := repo.GetRanks(ids)

	for i, page := range *input {
		key := fmt.Sprintf("%d", page.PageId)
		rank, ok := ranks[key]
		if !ok {
			continue
		}
		(*input)[i].Rank = rank
	}

	return input
}
