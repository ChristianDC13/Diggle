package crawler

import (
	"diggle/crawler/repository"
	"diggle/crawler/scrapper"
	domainExtractor "diggle/crawler/url-tools/domain-extractor"
	"diggle/crawler/webqueue"
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

func (c *Crawler) Crawl(maxPages, simultaneousRequests int, wq *webqueue.WebQueue) {

	showStatusFlag := flag.Bool("show-status", false, "Show status of the crawler")

	flag.Parse()

	ch := make(chan struct{}, simultaneousRequests)
	var i int64

	if *showStatusFlag {
		go c.showStatus(&i, maxPages)
	}

	for {

		if i >= int64(maxPages) {
			break
		}
		ch <- struct{}{}
		go func() {
			c.mu.Lock()
			url := wq.Dequeue()
			c.mu.Unlock()
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
				c.mu.Lock()
				wq.Enqueue(u)
				c.mu.Unlock()
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

	wq := webqueue.NewWebQueue()
	for _, url := range firstURLs {
		wq.Enqueue(url)
	}

	now := time.Now()
	c.Crawl(maxPages, simultaneousRequests, wq)
	elapsed := time.Since(now)

	hours := int(elapsed.Hours())
	elapsed -= time.Duration(hours) * time.Hour
	minutes := int(elapsed.Minutes())
	elapsed -= time.Duration(minutes) * time.Minute
	seconds := int(elapsed.Seconds())

	fmt.Printf("Crawling time: %dh:%dm:%ds\n", hours, minutes, seconds)
}
