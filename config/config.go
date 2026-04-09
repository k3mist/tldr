package config

import (
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/mitchellh/go-homedir"

	"bitbucket.org/djr2/tldr/color"
)

// Config provides the configuration variables read from config.json
var Config Options

// Options defines the available configuration options
type Options struct {
	PagesURI         string `json:"pages_uri"`
	ZipURI           string `json:"zip_uri"`
	Language         string `json:"language"`
	CacheExpiration  int    `json:"cache_expiration"`
	ExtendedSearch   bool   `json:"extended_search"`
	LookupWarnings   bool   `json:"lookup_warnings"`
	BannerColor1     int    `json:"banner_color_1"`
	BannerColor2     int    `json:"banner_color_2"`
	TLDRColor        int    `json:"tldr_color"`
	HeaderColor      int    `json:"header_color"`
	HeaderDecorColor int    `json:"header_decor_color"`
	PlatformColor    int    `json:"platform_color"`
	PlatformAltColor int    `json:"platform_alt_color"`
	DescriptionColor int    `json:"description_color"`
	ExampleColor     int    `json:"example_color"`
	HyphenColor      int    `json:"hyphen_color"`
	SyntaxColor      int    `json:"syntax_color"`
	VariableColor    int    `json:"variable_color"`
}

// Load looks for the config.json file in $HOME/.tldr. If the configuration
// file is not found it will create one with the default configuration options.
func Load() {
	h, err := homedir.Dir()
	if err != nil {
		log.Fatal(err)
	}

	f := h + "/" + ".tldr/config.json"
	if _, err := os.Stat(f); err != nil {
		create(f)
	}

	file, err := os.Open(f)
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()

	b, err := io.ReadAll(file)
	if err != nil {
		log.Println(err)
		return
	}

	if err := json.Unmarshal(b, &Config); err != nil {
		log.Println(err)
	}
}

func create(f string) {
	vars := Options{
		PagesURI:         "",
		ZipURI:           "",
		Language:         "",
		CacheExpiration:  30,
		ExtendedSearch:   true,
		LookupWarnings:   false,
		BannerColor1:     color.Cyan,
		BannerColor2:     color.Blue,
		TLDRColor:        color.White,
		HeaderColor:      color.Blue,
		HeaderDecorColor: color.White,
		PlatformColor:    color.DarkGray,
		PlatformAltColor: color.BrightPurple,
		DescriptionColor: color.Normal,
		ExampleColor:     color.Cyan,
		HyphenColor:      color.Normal,
		SyntaxColor:      color.Red,
		VariableColor:    color.Normal,
	}

	file, err := os.Create(f)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	j, err := json.MarshalIndent(vars, "", "")
	if err != nil {
		log.Fatal(err)
	}

	if _, err := file.Write(j); err != nil {
		log.Fatal(err)
	}
}
