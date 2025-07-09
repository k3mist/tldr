package language

import (
	"os"
	"slices"
	"strings"

	"bitbucket.org/djr2/tldr/config"
)

var languages []string

// Get a language by index
func GetLanguage(which int) string {
	return languages[which]
}

// Get all languages
func GetLanguages() []string {
	return languages
}

// Check if a language exists
func HasLanguage(which string) bool {
	return slices.Contains(languages, which)
}

// Set the languages
func SetLang() {
	var list []string

	if config.Config.Language != "" {
		list = append(list, parseLang(config.Config.Language)...)
	} else if os.Getenv("LANG") != "" {
		list = append(list, parseLang(os.Getenv("LANG"))...)
	} else if os.Getenv("LANGUAGE") != "" {
		list = append(list, parseLang(os.Getenv("LANGUAGE"))...)
	}

	languages = list
	if len(languages) == 0 || !slices.Contains(languages, "en") {
		languages = append(languages, "en")
	}
}

var keepSuffix []string = []string{"BR", "PT", "TW"}

func findSuffix(locale []string) string {
	var found string
	for _, s := range keepSuffix {
		if len(locale) == 2 && locale[1] == s {
			found = strings.Join(locale, "_")
			return found
		} else {
			found = locale[0]
			return found
		}
	}
	return found
}

func parseLang(lang string) []string {
	var split []string = strings.Split(lang, ":")
	var list []string
	for _, l := range split {
		var language []string = strings.Split(l, ".")
		var locale []string = strings.Split(language[0], "_")
		var found string = findSuffix(locale)
		list = append(list, found)
	}
	return list
}
