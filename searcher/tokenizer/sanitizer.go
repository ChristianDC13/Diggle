package tokenizer

import (
	"regexp"
	"strings"
)

var sanitizingPatterns = map[string]string{
	"english":    "[^a-zA-Z0-9]+",
	"spanish":    "[^a-zA-Z0-9áéíóúÁÉÍÓÚñÑüÜ]+",
	"german":     "[^a-zA-Z0-9äöüÄÖÜß]+",
	"russian":    "[^a-zA-Z0-9а-яА-Я]+",
	"arabic":     "[^a-zA-Z0-9ء-ي]+",
	"chinese":    "[^a-zA-Z0-9\u4E00-\u9FA5]+",
	"japanese":   "[^a-zA-Z0-9\u3040-\u309F\u30A0-\u30FF\u4E00-\u9FA5]+",
	"korean":     "[^a-zA-Z0-9\uAC00-\uD7AF\u1100-\u11FF\u3130-\u318F]+",
	"thai":       "[^a-zA-Z0-9ก-๙]+",
	"persian":    "[^a-zA-Z0-9آ-ی]+",
	"indonesian": "[^a-zA-Z0-9]+",
	"malay":      "[^a-zA-Z0-9]+",
	"filipino":   "[^a-zA-Z0-9]+",
	"portuguese": "[^a-zA-Z0-9áàãâéêíóôõúçÁÀÃÂÉÊÍÓÔÕÚÇ]+",
	"french":     "[^a-zA-Z0-9àâçéèêëîïôûùüÿñæœÀÂÇÉÈÊËÎÏÔÛÙÜŸÑÆŒ]+",
	"italian":    "[^a-zA-Z0-9àèéìíîòóùú]+",
	"dutch":      "[^a-zA-Z0-9]+",
	"polish":     "[^a-zA-Z0-9ąćęłńóśźżĄĆĘŁŃÓŚŹŻ]+",
	"ukrainian":  "[^a-zA-Z0-9а-яА-ЯіїєґІЇЄҐ]+",
	"romanian":   "[^a-zA-Z0-9ăâîșțĂÂÎȘȚ]+",
	"greek":      "[^a-zA-Z0-9α-ωΑ-ΩάέήίϊΐόύϋΰώΆΈΉΊΌΎΏ]+",
	"hungarian":  "[^a-zA-Z0-9áéíóöőúüűÁÉÍÓÖŐÚÜŰ]+",
	"finnish":    "[^a-zA-Z0-9äöåÄÖÅ]+",
	"swedish":    "[^a-zA-Z0-9åäöÅÄÖ]+",
	"norwegian":  "[^a-zA-Z0-9æøåÆØÅ]+",
	"danish":     "[^a-zA-Z0-9æøåÆØÅ]+",
	"icelandic":  "[^a-zA-Z0-9áðéíóúýþæöÁÐÉÍÓÚÝÞÆÖ]+",
	"czech":      "[^a-zA-Z0-9áčďéěíňóřšťúůýžÁČĎÉĚÍŇÓŘŠŤÚŮÝŽ]+",
	"slovak":     "[^a-zA-Z0-9áäčďéíĺľňóôŕšťúýžÁÄČĎÉÍĹĽŇÓÔŔŠŤÚÝŽ]+",
	"croatian":   "[^a-zA-Z0-9čćdžđšžČĆDŽĐŠŽ]+",
	"bosnian":    "[^a-zA-Z0-9čćdžđšžČĆDŽĐŠŽ]+",
	"albanian":   "[^a-zA-Z0-9ëçëÇ]+",
	"estonian":   "[^a-zA-Z0-9äöõüšžÄÖÕÜŠŽ]+",
	"latvian":    "[^a-zA-Z0-9āčēģīķļņšūžĀČĒĢĪĶĻŅŠŪŽ]+",
	"lithuanian": "[^a-zA-Z0-9ąčęėįšųūžĄČĘĖĮŠŲŪŽ]+",
	"georgian":   "[^a-zA-Z0-9ა-ჰ]+",
	"default":    "[^a-zA-Z0-9]+",
}

func Sanitize(input []string, language string) []string {
	var sanitized []string
	language = strings.ToLower(language)
	pattern := sanitizingPatterns[language]
	if pattern == "" {
		pattern = sanitizingPatterns["default"]
	}
	for _, word := range input {
		item := regexp.MustCompile(pattern).ReplaceAllString(word, "")
		item = strings.ToLower(item)
		if item != "" {

			sanitized = append(sanitized, item)
		}
	}
	return sanitized
}
