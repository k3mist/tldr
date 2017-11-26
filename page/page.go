package page

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"regexp"

	"bitbucket.org/djr2/tldr/color"
	"bitbucket.org/djr2/tldr/config"
	"bitbucket.org/djr2/tldr/platform"
)

type Page interface {
	Print()
	Write(p []byte)
	Lines() [][]byte
	header() []byte
	example(line []byte) []byte
	syntax(line []byte) []byte
}

var (
	descRx = regexp.MustCompile(`^>\s`)
	varRx  = regexp.MustCompile(`{{([\w\s\\/~!@#$%^&*()\[\]:;"'<,>?.]+)}}`)
)

func New(file *os.File, plat platform.Platform) Page {
	b, err := ioutil.ReadAll(file)
	defer file.Close()
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
	cfg := config.Config
	p.Write(to_b("\n"))
	for i, line := range p.Lines() {
		if i == 0 {
			p.Write(to_b("  "))
			p.Write(p.header())
			p.Write(to_b(color.ColorBold(cfg.HeaderDecorColor) + " - "))
			p.Write(to_b(color.Color(cfg.PlatformColor) + plat.String()))
			p.Write(to_b(color.ColorBold(cfg.HeaderDecorColor) + "]" + "\n"))
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

		if syntax := p.syntax(line); syntax != nil {
			p.Write(to_b("    "))
			p.Write(variable(syntax))
			p.Write(to_b("\n"))
			continue
		}
	}
}

func description(line []byte) []byte {
	cfg := config.Config
	if descRx.Match(line) {
		return descRx.ReplaceAll(line, to_b(color.Color(cfg.DescriptionColor)))
	}
	return nil
}

func variable(line []byte) []byte {
	cfg := config.Config
	return varRx.ReplaceAll(line, to_b(color.Color(cfg.VariableColor)+"$1"+color.Color(cfg.SyntaxColor)))
}

func to_b(str string) []byte {
	return []byte(str)
}
