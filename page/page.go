package page

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	"bitbucket.org/djr2/tldr/color"
)

type parser uint32

const (
	v1 parser = iota
	v2
)

type Page interface {
	Print()
	header() string
	example(line string) string
	code(line string) string
}

var (
	descRx = regexp.MustCompile(`^>\s`)
	varRx  = regexp.MustCompile(`{{([\w\s\\/~!@#$%^&*()\[\]:;"'<,>?.]+)}}`)
)

func NewPage(file *os.File) Page {
	b, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	parse := v1
	contents := strings.Split(string(b), "\n")
	if headerRxV2.MatchString(contents[1]) {
		parse = v2
		contents[1] = contents[0]
		contents = contents[1:]
		return &pagev2{contents, &bytes.Buffer{}, parse}
	}
	return &pagev1{contents, &bytes.Buffer{}, parse}
}

func description(line string) string {
	if descRx.MatchString(line) {
		return descRx.ReplaceAllString(line, color.Color(color.Normal))
	}
	return ""
}

func variable(line string) string {
	return varRx.ReplaceAllString(line, color.Color(color.Normal)+"$1"+color.ColorNormal(color.Red))
}
