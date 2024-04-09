package pipes

import (
	"diggle/searcher/models"
	"sort"
)

func Sort(input *[]models.RawPage) *[]models.RawPage {

	sort.Slice(*input, func(i, j int) bool {

		if (*input)[i].HitsCount != (*input)[j].HitsCount {
			return (*input)[i].HitsCount > (*input)[j].HitsCount
		}
		if (*input)[i].PositionalDistance != (*input)[j].PositionalDistance {
			return (*input)[i].PositionalDistance < (*input)[j].PositionalDistance
		}

		if (*input)[i].HitsSum != (*input)[j].HitsSum {
			return (*input)[i].HitsSum > (*input)[j].HitsSum
		}

		return (*input)[i].Rank > (*input)[j].Rank
	})
	return input
}
