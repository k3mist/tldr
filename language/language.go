package language

import (
	"os"
	"slices"
	"strings"

	"bitbucket.org/djr2/tldr/config"
)

var lang []string
var configLang string
var envLang string
var envLang2 string

func init() {
	configLang = config.Config.Language
	envLang = os.Getenv("LANG")
	envLang2 = os.Getenv("LANGUAGE")

	lang = setLang()
	if len(lang) == 0 {
		lang = append(lang, "en")
	}
}

func GetLanguage(which int) string {
	return lang[which]
}

func GetLanguages(which int) []string {
	return lang
}

func HasLanguage(language string) bool {
	return slices.Contains(lang, language)
}

func setLang() []string {
	var list []string

	if configLang != "" {
		list = append(list, parseLang(configLang)...)
		return list
	}

	if envLang != "" {
		list = append(list, parseLang(envLang)...)
		return list
	}

	if envLang2 != "" {
		list = append(list, parseLang(envLang2)...)
		return list
	}

	return list
}

var keepSuffix []string = []string{"BR", "PT", "TW"}

func findSuffix(locale []string) string {
	var found string
	for _, s := range keepSuffix {
		if locale[1] == s {
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
