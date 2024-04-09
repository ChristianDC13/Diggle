package scrapper

import (
	"context"
	htmlextractor "diggle/crawler/htmlextractor"
	domainExtractor "diggle/crawler/url-tools/domain-extractor"
	"diggle/crawler/webpage"
	"io"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/html"
)

type Scrapper struct {
	httpClient *http.Client
}

func NewScrapper(httpClient *http.Client) *Scrapper {
	return &Scrapper{httpClient}
}

func (s *Scrapper) Scrape(siteUrl string) (*webpage.WebPage, error) {

	outboundLinks := make(map[string]bool)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", siteUrl, nil)

	if err != nil {
		return nil, err
	}

	res, err := s.httpClient.Do(req)

	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	doc, _ := html.Parse(strings.NewReader(string(b)))

	page := htmlextractor.GetPageInfo(doc)
	page.URL = siteUrl

	for _, link := range page.Links {
		if _, err := domainExtractor.ExtractDomain(link); err == nil {
			outboundLinks[link] = true
		}
	}

	for domain := range outboundLinks {
		page.OutboundLinks = append(page.OutboundLinks, domain)
	}

	return page, nil
}
