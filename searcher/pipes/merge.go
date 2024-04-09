package pipes

import (
	"diggle/searcher/models"
)

func Merge(input *map[string][]models.SearchDBResult) *[]models.RawPage {
	resultsMap := map[int64]*models.RawPage{}

	for _, results := range *input {
		for _, result := range results {
			if _, ok := resultsMap[result.PageId]; !ok {
				resultsMap[result.PageId] = &models.RawPage{
					PageId:             result.PageId,
					HitsCount:          0,
					HitsSum:            0,
					Rank:               0,
					PositionalDistance: 0,
				}
			}

			resultsMap[result.PageId].HitsCount++
			resultsMap[result.PageId].HitsSum += result.Score
			resultsMap[result.PageId].PositionalDistance += result.Position
		}
	}

	results := make([]models.RawPage, 0, len(resultsMap))
	for _, result := range resultsMap {
		results = append(results, *result)
	}

	return &results

}
