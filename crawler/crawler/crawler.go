package crawler

import (
	"diggle/crawler/repository"
	"diggle/crawler/scrapper"
	domainExtractor "diggle/crawler/url-tools/domain-extractor"
	"flag"
	"fmt"
	"sync"
	"time"
)

type Crawler struct {
	scrapper   *scrapper.Scrapper
	repository *repository.Repository
	mu         *sync.Mutex
}

func (c *Crawler) Crawl(maxPages, simultaneousRequests int) {

	showStatusFlag := flag.Bool("show-status", false, "Show status of the crawler")

	flag.Parse()

	ch := make(chan struct{}, simultaneousRequests)
	var i int64

	if *showStatusFlag {
		go c.showStatus(&i, maxPages)
	}

	for {
		i = (*c.repository).GetPagesCount()
		if i >= int64(maxPages) {
			break
		}
		ch <- struct{}{}
		go func() {

			url, err := (*c.repository).DequeuePage()
			if err != nil {
				<-ch
				return
			}
			if url == "" {
				<-ch
				return
			}

			page, err := c.scrapper.Scrape(url)

			if err != nil {
				<-ch
				return
			}

			for _, u := range page.Links {
				_, err = domainExtractor.ExtractDomain(u)

				if err != nil {
					continue
				}
				(*c.repository).EnqueuePage(u)
			}

			added, err := (*c.repository).AddPage(page)
			if err != nil {
				<-ch
				return
			}
			if added {
				i++
			}

			<-ch
		}()
	}

}

func NewCrawler(scrapper *scrapper.Scrapper, repo repository.Repository) *Crawler {
	mu := &sync.Mutex{}
	return &Crawler{scrapper: scrapper, repository: &repo, mu: mu}
}

func (c *Crawler) Start(firstURLs []string, maxPages, simultaneousRequests int) {

	for _, url := range firstURLs {
		(*c.repository).EnqueuePage(url)
	}

	now := time.Now()
	c.Crawl(maxPages, simultaneousRequests)
	elapsed := time.Since(now)

	hours := int(elapsed.Hours())
	elapsed -= time.Duration(hours) * time.Hour
	minutes := int(elapsed.Minutes())
	elapsed -= time.Duration(minutes) * time.Minute
	seconds := int(elapsed.Seconds())

	fmt.Printf("Crawling time: %dh:%dm:%ds\n", hours, minutes, seconds)
}
