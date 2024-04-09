package tokenizer

import (
	"strings"

	"github.com/rylans/getlang"
)

func Tokenize(text string) []string {
	info := getlang.FromString(text)
	language := info.LanguageName()
	confidence := info.Confidence()

	result := strings.Fields(text)

	if confidence < 0.6 {
		language = "english"
	}

	result = Sanitize(result, language)

	result, err := RemoveStopWords(result, language)
	if err != nil {
		return result
	}
	result = StemWords(result, language)

	return result
}
