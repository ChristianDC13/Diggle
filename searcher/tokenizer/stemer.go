package tokenizer

import (
	"strings"

	"github.com/tebeka/snowball"
)

func StemWords(words []string, language string) []string {
	result := []string{}
	language = strings.ToLower(language)
	stemmer, err := snowball.New(language)
	if err != nil {
		return words
	}
	defer stemmer.Close()

	for _, word := range words {
		stemmed := stemmer.Stem(word)
		result = append(result, stemmed)
	}

	return result
}
