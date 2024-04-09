package htmlparser

import (
	"diggle/crawler/tokenizer"
	"diggle/crawler/webpage"

	"golang.org/x/net/html"
)

func GetPageInfo(doc *html.Node) *webpage.WebPage {
	page := &webpage.WebPage{}
	currentPosition := 0
	doTraverse(doc, page, &currentPosition)
	return page
}

func doTraverse(doc *html.Node, page *webpage.WebPage, currentPosition *int) {

	var traverse func(n *html.Node) *html.Node

	traverse = func(n *html.Node) *html.Node {

		for c := n.FirstChild; c != nil; c = c.NextSibling {

			if c.Type == html.ElementNode && c.Data == "a" {
				for _, a := range c.Attr {
					if a.Key == "href" {
						page.Links = append(page.Links, a.Val)
					}
				}
			}

			if c.Type == html.ElementNode && c.Data == "title" {
				if c.FirstChild != nil {
					page.Title = c.FirstChild.Data
				}
			}

			if c.Type == html.ElementNode && c.Data == "meta" {
				for _, a := range c.Attr {
					if a.Key == "name" && a.Val == "description" || a.Val == "og:description" || a.Val == "twitter:description" {
						content := ""
						for _, a := range c.Attr {
							if a.Key == "content" {
								content = a.Val
							}
						}
						page.Abstract = content
					}
				}
			}

			if page.Abstract == "" || len(page.Abstract) < 20 {
				for attempt := 0; attempt < 4; attempt++ {
					if c.Type == html.ElementNode && c.Data == "p" && page.Abstract == "" {
						if c.FirstChild != nil {
							page.Abstract = c.FirstChild.Data
						}
						if len(page.Abstract) > 20 {
							break
						}
					}
				}
			}

			if len(page.Abstract) > 150 {
				page.Abstract = page.Abstract[:150]
			}

			textElements := map[string]bool{
				"title": true,
				"h1":    true,
				"h2":    true,
				"h3":    true,
				"h4":    true,
				"h5":    true,
				"h6":    true,
				"p":     true,
			}

			if c.Type == html.ElementNode && textElements[c.Data] {
				extractWords(c, page, currentPosition)
			}

			res := traverse(c)

			if res != nil {
				return res
			}
		}

		return nil
	}

	traverse(doc)
}

func extractWords(c *html.Node, page *webpage.WebPage, currentPosition *int) {

	scoreTable := map[string]int{
		"title": 20,
		"h1":    10,
		"h2":    5,
		"h3":    4,
		"h4":    3,
		"h5":    2,
		"p":     1,
	}

	if c.FirstChild != nil {
		wordsMap := make(map[string]struct {
			Score    int
			Position int
		})
		sentence := c.FirstChild.Data

		words := tokenizer.Tokenize(sentence)
		tag := c.Data

		for _, word := range words {
			if word != "" {
				if wordInfo, ok := wordsMap[word]; ok && scoreTable[tag] < wordInfo.Score {
					continue
				}
				wordsMap[word] = struct {
					Score    int
					Position int
				}{Score: scoreTable[tag], Position: *currentPosition}
				(*currentPosition)++
			}
		}
		for word, wordData := range wordsMap {
			page.Words = append(page.Words, webpage.Word{Word: word, Score: wordData.Score, Position: wordData.Position})
		}
	}

}
