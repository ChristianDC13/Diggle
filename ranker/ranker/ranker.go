package ranker

import (
	"diggle/ranker/repository"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/dcadenas/pagerank"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type Ranker struct {
	repo         *repository.Repository
	pageRanks    *map[string][]float64
	mu           *sync.Mutex
	pages        *[]string
	deadEnds     *map[string]struct{}
	deadEndsList *[]string
}

func NewRanker(repo repository.Repository) *Ranker {
	pr := make(map[string][]float64)
	mu := sync.Mutex{}
	return &Ranker{
		repo:      &repo,
		pageRanks: &pr,
		mu:        &mu,
	}
}

func (r *Ranker) Rank() error {

	// timeStart := time.Now()
	p := message.NewPrinter(language.English)
	pages, err := (*r.repo).GetPages()

	if err != nil {
		return err
	}

	r.pages = &pages

	graph := pagerank.New()
	start := time.Now()
	fmt.Printf("Building graph")
	i := 0
	// go showStatus(&i, len(*r.pages))

	for _, page := range *r.pages {
		links, _ := (*r.repo).GetOutboundLinks(page)

		for _, link := range links {
			currentPageInt, _ := strconv.Atoi(page)
			linkInt, _ := strconv.Atoi(link)
			graph.Link(currentPageInt, linkInt)
		}
		if i%1000 == 0 {
			fmt.Print("\033[H\033[2J")
			p.Printf("Creating graph %d / %d\n", i, len(pages))
			parcentage := (i * 100) / len(pages)
			bar := ""
			for j := 0; j < parcentage; j++ {
				bar += "="
			}
			fmt.Printf("[%s%s]%d%%\n", bar, strings.Repeat(" ", 100-parcentage), parcentage)
		}
		i++

	}

	fmt.Printf("Graph built in %v\n", time.Since(start))

	probability_of_following_a_link := 0.85 // The bigger the number, less probability we have to teleport to some random link
	tolerance := 0.0001                     // the smaller the number, the more exact the result will be but more CPU cycles will be needed

	startTime := time.Now()
	i = 0
	fmt.Print("Starting ranking\n")

	graph.Rank(probability_of_following_a_link, tolerance, func(identifier int, rank float64) {
		pageId := strconv.Itoa(identifier)
		(*r.repo).SetPageRank(pageId, rank)
		if i%1000 == 0 {
			fmt.Print("\033[H\033[2J")
			p.Printf("Ranking %d / %d\n", i, len(pages))
			parcentage := (i * 100) / len(pages)
			bar := ""
			for j := 0; j < parcentage; j++ {
				bar += "="
			}
			fmt.Printf("[%s%s]%d%%\n", bar, strings.Repeat(" ", 100-parcentage), parcentage)
		}
		i++
	})
	fmt.Print("\033[H\033[2J")
	fmt.Printf("Ranking Time: %v\n", time.Since(startTime))

	return nil
}
