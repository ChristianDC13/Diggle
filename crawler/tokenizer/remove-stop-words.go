package tokenizer

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var languageAbrv = map[string]string{
	"afrikaans":      "af",
	"albanian":       "sq",
	"arabic":         "ar",
	"armenian":       "hy",
	"azerbaijani":    "az",
	"basque":         "eu",
	"belarusian":     "be",
	"bengali":        "bn",
	"bosnian":        "bs",
	"bulgarian":      "bg",
	"catalan":        "ca",
	"chinese":        "zh",
	"croatian":       "hr",
	"czech":          "cs",
	"danish":         "da",
	"dutch":          "nl",
	"english":        "en",
	"esperanto":      "eo",
	"estonian":       "et",
	"finnish":        "fi",
	"french":         "fr",
	"galician":       "gl",
	"georgian":       "ka",
	"german":         "de",
	"greek":          "el",
	"gujarati":       "gu",
	"haitian creole": "ht",
	"hebrew":         "he",
	"hindi":          "hi",
	"hungarian":      "hu",
	"icelandic":      "is",
	"indonesian":     "id",
	"irish":          "ga",
	"italian":        "it",
	"japanese":       "ja",
	"javanese":       "jv",
	"kannada":        "kn",
	"kazakh":         "kk",
	"khmer":          "km",
	"korean":         "ko",
	"kurdish":        "ku",
	"kyrgyz":         "ky",
	"lao":            "lo",
	"latin":          "la",
	"latvian":        "lv",
	"lithuanian":     "lt",
	"luxembourgish":  "lb",
	"macedonian":     "mk",
	"malagasy":       "mg",
	"malay":          "ms",
	"malayalam":      "ml",
	"maltese":        "mt",
	"maori":          "mi",
	"marathi":        "mr",
	"mongolian":      "mn",
	"myanmar":        "my",
	"nepali":         "ne",
	"norwegian":      "no",
	"odia":           "or",
	"pashto":         "ps",
	"persian":        "fa",
	"polish":         "pl",
	"portuguese":     "pt",
	"punjabi":        "pa",
	"romanian":       "ro",
	"russian":        "ru",
	"samoan":         "sm",
	"scots gaelic":   "gd",
	"serbian":        "sr",
	"sesotho":        "st",
	"shona":          "sn",
	"sindhi":         "sd",
	"sinhala":        "si",
	"slovak":         "sk",
	"slovenian":      "sl",
	"somali":         "so",
	"spanish":        "es",
	"sundanese":      "su",
	"swahili":        "sw",
	"swedish":        "sv",
	"tajik":          "tg",
	"tamil":          "ta",
	"tatar":          "tt",
	"telugu":         "te",
	"thai":           "th",
	"turkish":        "tr",
	"turkmen":        "tk",
	"ukrainian":      "uk",
	"urdu":           "ur",
	"uyghur":         "ug",
	"uzbek":          "uz",
	"vietnamese":     "vi",
	"welsh":          "cy",
	"xhosa":          "xh",
	"yiddish":        "yi",
	"yoruba":         "yo",
	"zulu":           "zu",
}

var stopWords = map[string]map[string]struct{}{}

func RemoveStopWords(words []string, language string) ([]string, error) {

	language = strings.ToLower(language)

	if len(stopWords[language]) == 0 {
		err := loadStopWords(language)
		if err != nil {
			return words, err
		}
	}

	currentStopWords := stopWords[language]

	if len(currentStopWords) == 0 {
		return words, nil
	}

	result := []string{}

	for _, word := range words {

		if _, ok := currentStopWords[word]; ok {
			continue
		}
		result = append(result, word)
	}

	return result, nil
}

func loadStopWords(language string) error {

	_, b, _, _ := runtime.Caller(0)

	root := filepath.Join(filepath.Dir(b))

	lng := languageAbrv[language]

	file, err := os.Open(root + "/stopwords/" + lng + ".json")

	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	sw := []string{}
	err = decoder.Decode(&sw)

	mapSw := map[string]struct{}{}

	for _, word := range sw {
		mapSw[word] = struct{}{}
	}

	if err != nil {
		return err
	}

	stopWords[language] = mapSw

	return nil
}
