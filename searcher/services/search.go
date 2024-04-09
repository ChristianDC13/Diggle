package services

import (
	fuzzysearch "diggle/searcher/fuzzy-search"
	"diggle/searcher/models"
	"diggle/searcher/repository"
	"strconv"
	"strings"
	"sync"
)

type SearchService struct {
	repo          *repository.RedisRepository
	fuzzySearcher *fuzzysearch.FuzzySearcher
}

func NewSearchService(repo *repository.RedisRepository, fuzzySearcher *fuzzysearch.FuzzySearcher) *SearchService {
	return &SearchService{repo, fuzzySearcher}
}

func (ss *SearchService) Search(tokens []string) (map[string][]models.SearchDBResult, error) {

	results := map[string][]models.SearchDBResult{}

	wg := sync.WaitGroup{}
	wg.Add(len(tokens))
	mu := sync.Mutex{}
	for _, token := range tokens {
		go func() {
			defer wg.Done()

			pages := ss.repo.GetPagesForWord(token)
			if len(pages) == 0 {
				fuzzTerms := ss.fuzzySearcher.Search(token)
				if len(fuzzTerms) == 0 {
					return
				}
				term := fuzzTerms[0]
				pages = ss.repo.GetPagesForWord(term)
			}

			if len(pages) == 0 {
				return
			}

			res := []models.SearchDBResult{}

			for page, info := range pages {
				pageInt, _ := strconv.Atoi(page)
				socore := strings.Split(info, ":")[0]
				position := strings.Split(info, ":")[1]
				scoreInt, _ := strconv.Atoi(socore)
				positionInt, _ := strconv.Atoi(position)
				res = append(res, models.SearchDBResult{PageId: int64(pageInt), Score: scoreInt, Position: positionInt})
			}

			mu.Lock()
			results[token] = res
			mu.Unlock()
		}()
	}
	wg.Wait()
	return results, nil
}
