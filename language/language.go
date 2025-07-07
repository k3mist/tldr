package language

import (
	"os"
	"slices"
	"strings"

	"bitbucket.org/djr2/tldr/config"
)

var lang []string

func init() {
	lang = setLang()
	if len(lang) == 0 {
		lang = append(lang, "en")
	}
}

func GetLanguage(which int) string {
	return lang[which]
}

func GetLanguages() []string {
	return lang
}

func HasLanguage(which string) bool {
	return slices.Contains(lang, which)
}

func setLang() []string {
	var list []string

	if config.Config.Language != "" {
		list = append(list, parseLang(config.Config.Language)...)
		return list
	}

	if os.Getenv("LANG") != "" {
		list = append(list, parseLang(os.Getenv("LANG"))...)
		return list
	}

	if os.Getenv("LANGUAGE") != "" {
		list = append(list, parseLang(os.Getenv("LANGUAGE"))...)
		return list
	}

	return list
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
