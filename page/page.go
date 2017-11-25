package page

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"regexp"

	"bitbucket.org/djr2/tldr/color"
	"bitbucket.org/djr2/tldr/platform"
)

type Page interface {
	Print()
	Write(p []byte)
	Lines() [][]byte
	header() []byte
	example(line []byte) []byte
	code(line []byte) []byte
}

var (
	descRx = regexp.MustCompile(`^>\s`)
	varRx  = regexp.MustCompile(`{{([\w\s\\/~!@#$%^&*()\[\]:;"'<,>?.]+)}}`)
)

func NewPage(file *os.File, plat platform.Platform) Page {
	b, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
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
	parse(p, plat)
	return p
}

func parse(p Page, plat platform.Platform) {
	p.Write(to_b("\n"))
	for i, line := range p.Lines() {
		if i == 0 {
			p.Write(to_b("  "))
			p.Write(p.header())
			p.Write(to_b(color.ColorBold(color.White) + " - "))
			p.Write(to_b(color.Color(color.DarkGray) + plat.String()))
			p.Write(to_b(color.ColorBold(color.White) + "]" + "\n"))
			continue
		}

		if desc := description(line); desc != nil {
			p.Write(to_b("  "))
			p.Write(desc)
			p.Write(to_b("\n"))
			continue
		}

		if example := p.example(line); example != nil {
			p.Write(to_b("\n  "))
			p.Write(example)
			p.Write(to_b("\n"))
			continue
		}

		if code := p.code(line); code != nil {
			p.Write(to_b("    "))
			p.Write(variable(code))
			p.Write(to_b("\n"))
			continue
		}
	}
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
