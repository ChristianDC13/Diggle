package main

import (
	"diggle/monitor/repository"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func main() {
	godotenv.Load()
	repo := repository.NewRedisRepository()
	maxPages := os.Getenv("MAX_PAGES")
	maxPagesInt, err := strconv.Atoi(maxPages)
	if err != nil {
		panic(err)
	}

	i := int64(0)
	go showStatus(&i, maxPagesInt)

	for {
		time.Sleep(1 * time.Second)
		i = repo.GetPagesCount()
		if i == int64(maxPagesInt) {
			break
		}
	}

}

func showStatus(i *int64, maxPages int) {

	timeStart := time.Now()
	p := message.NewPrinter(language.English)
	for {
		time.Sleep(1 * time.Second)
		fmt.Print("\033[H\033[2J")
		parcentage := int((*i * 100) / int64(maxPages))
		p.Printf("Crawling...\n%d / %d\n", *i, maxPages)
		bar := ""
		for j := 0; j < parcentage; j++ {
			bar += "="
		}
		fmt.Printf("[%s%s]%d%%\n", bar, strings.Repeat(" ", 100-parcentage), parcentage)

		timeElapsed := time.Since(timeStart).Seconds()
		timeRemaining := (timeElapsed / float64(*i)) * float64(int64(maxPages)-*i)
		hoursRemaining := int(timeRemaining) / 3600
		minutesRemaining := int(timeRemaining) % 3600 / 60
		secondsRemaining := int(timeRemaining) % 60
		fmt.Printf("Time remaining: %02dh:%02dm:%02ds\n", hoursRemaining, minutesRemaining, secondsRemaining)
	}
}
