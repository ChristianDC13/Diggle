package fuzzysearch

import (
	"diggle/searcher/repository"
	"fmt"
	"sort"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type FuzzySearcher struct {
	repo      *repository.RedisRepository
	bkt       *BKTree
	tolerance int
}

func (f *FuzzySearcher) loadDictionary() error {
	words, err := f.repo.GetAllWords()
	if err != nil {
		return err
	}

	p := message.NewPrinter(language.English)
	for i := 0; i < len(words); i++ {
		if i%1000 == 0 {
			fmt.Print("\033[H\033[2J")
			p.Printf("Loading dictionary: %d/%d\n", i, len(words))
		}
		word := words[i]
		f.bkt.Add(word)
	}
	fmt.Print("\033[H\033[2J")
	p.Printf("Dictionary loaded: %d words\n", len(words))
	return nil
}

func (f *FuzzySearcher) Search(word string) []string {
	results := f.bkt.Search(word, f.tolerance)

	sort.Slice(results, func(i, j int) bool {
		if results[i].distance != results[j].distance {
			return results[i].distance < results[j].distance
		}
		return results[i].frequency > results[j].frequency

	})

	var res []string
	min := -1
	for _, r := range results {

		if min == -1 {
			min = r.distance
		}

		if r.distance > min {
			break
		}

		res = append(res, r.word)
	}

	return res
}

func NewFuzzySearcher(repo *repository.RedisRepository, tolerance int) (*FuzzySearcher, error) {
	bkt := NewBKTree()
	f := &FuzzySearcher{repo, bkt, tolerance}
	err := f.loadDictionary()
	return f, err
}
