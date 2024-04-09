package handlers

import (
	"diggle/searcher/pipes"
	"diggle/searcher/tokenizer"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

func (h *Handler) HandleSearch(w http.ResponseWriter, r *http.Request) {
	timeStart := time.Now()
	searchTerm := r.URL.Query().Get("q")
	page := r.URL.Query().Get("page")
	size := r.URL.Query().Get("size")

	if page == "" {
		page = "1"
	}

	if size == "" {
		size = "10"
	}

	pageInt, _ := strconv.Atoi(page)
	sizeInt, _ := strconv.Atoi(size)

	tokens := tokenizer.Tokenize(searchTerm)

	pages, _ := h.searchService.Search(tokens)

	merged := pipes.Merge(&pages)
	ranked := pipes.Rank(merged, h.repo)
	sorted := pipes.Sort(ranked)
	filled := pipes.Fill(sorted, sizeInt, pageInt, h.repo)

	elapsed := time.Since(timeStart).Milliseconds()

	totalResults := len(*sorted)
	pagesCount := float32(totalResults) / float32(sizeInt)
	if totalResults%sizeInt != 0 {
		pagesCount++
	}

	response := map[string]interface{}{
		"pages":        filled,
		"time":         elapsed,
		"totalResults": len(*sorted),
		"page":         pageInt,
		"pageSize":     sizeInt,
		"pagesCount":   int(pagesCount),
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(response)

}
