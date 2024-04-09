package main

import (
	"diggle/crawler/crawler"
	"diggle/crawler/repository"
	"diggle/crawler/scrapper"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load()

	firstURLsStr := os.Getenv("FIRST_WEB_URLS")
	firstURLs := strings.Split(firstURLsStr, ",")

	maxPagesStr := os.Getenv("MAX_PAGES")
	maxPages, err := strconv.Atoi(maxPagesStr)

	if err != nil {
		panic(err)
	}

	simultaneousRequestsStr := os.Getenv("SIMULTANEOUS_REQUESTS")
	simultaneousRequests, err := strconv.Atoi(simultaneousRequestsStr)

	if err != nil {
		panic(err)
	}

	t := http.DefaultTransport.(*http.Transport).Clone()
	t.DisableKeepAlives = true
	httpClient := &http.Client{Transport: t}
	scrapper := scrapper.NewScrapper(httpClient)

	repo := repository.NewRedisRepository()

	crawler := crawler.NewCrawler(scrapper, repo)

	crawler.Start(firstURLs, maxPages, simultaneousRequests)

}
