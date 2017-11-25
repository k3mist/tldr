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

type Page interface {
	Print()
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
	var p Page
	lines := strings.Split(string(b), "\n")
	if headerRxV2.MatchString(lines[1]) {
		lines[1] = lines[0]
		lines = lines[1:]
		p = &pagev2{lines, &bytes.Buffer{}}
	} else {
		p = &pagev1{lines, &bytes.Buffer{}}
	}
	return p
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
