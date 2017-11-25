package page

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"regexp"

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
	lines := bytes.Split(b, to_b("\n"))
	if headerRxV2.Match(lines[1]) {
		lines[1] = lines[0]
		lines = lines[1:]
		p = &pagev2{lines, &bytes.Buffer{}}
	} else {
		p = &pagev1{lines, &bytes.Buffer{}}
	}
	return p
}

func description(line []byte) []byte {
	if descRx.Match(line) {
		return descRx.ReplaceAll(line, to_b(color.Color(color.Normal)))
	}
	return nil
}

func variable(line []byte) []byte {
	return varRx.ReplaceAll(line, to_b(color.Color(color.Normal)+"$1"+color.ColorNormal(color.Red)))
}

func to_b(str string) []byte {
	return []byte(str)
}
