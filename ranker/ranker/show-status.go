package ranker

import (
	"fmt"
	"runtime"
	"strings"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func showStatus(current *int, total int) {

	p := message.NewPrinter(language.English)
	for {

		time.Sleep(100 * time.Millisecond)

		fmt.Print("\033[H\033[2J")
		parcentage := (*current * 100) / total
		fmt.Println("Ranking...")
		p.Printf("Page: %d / %d\n", *current, total)

		bar := ""
		for j := 0; j < parcentage; j++ {
			bar += "="
		}
		fmt.Printf("[%s%s]%d%%\n", bar, strings.Repeat(" ", 100-parcentage), parcentage)
		rutines := runtime.NumGoroutine()
		fmt.Printf("Goroutines: %d\n", rutines)

		if *current == total-1 {
			break
		}
	}
}
