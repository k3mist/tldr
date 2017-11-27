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
}

type Parser interface {
	Write(p []byte)
	Lines() [][]byte
	Header() []byte
	Description(line []byte) []byte
	Example(line []byte) []byte
	Syntax(line []byte) []byte
	Variable(line []byte) []byte
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
	lines := bytes.Split(b, to_b("\n"))
	if headerRxV2.Match(lines[1]) {
		p := &pagev2{lines, &bytes.Buffer{}}
		parse(p, plat)
		return p
	}
	p := &pagev1{lines, &bytes.Buffer{}}
	parse(p, plat)
	return p
}

func parse(p Parser, plat platform.Platform) {
	cfg := config.Config
	p.Write(to_b("\n"))
	for i, line := range p.Lines() {
		if i == 0 {
			p.Write(to_b("  "))
			p.Write(to_b(color.ColorBold(cfg.HeaderDecorColor) + "["))
			p.Write(p.Header())
			p.Write(to_b(color.ColorBold(cfg.HeaderDecorColor) + " - "))
			p.Write(to_b(color.Color(cfg.PlatformColor) + plat.String()))
			p.Write(to_b(color.ColorBold(cfg.HeaderDecorColor) + "]" + "\n"))
			continue
		}

		if desc := p.Description(line); desc != nil {
			p.Write(to_b("  "))
			p.Write(desc)
			p.Write(to_b("\n"))
			continue
		}

		if example := p.Example(line); example != nil {
			p.Write(to_b("\n  "))
			p.Write(example)
			p.Write(to_b("\n"))
			continue
		}

		if syntax := p.Syntax(line); syntax != nil {
			p.Write(to_b("    "))
			p.Write(p.Variable(syntax))
			p.Write(to_b("\n"))
			continue
		}
	}
}

func to_b(str string) []byte {
	return []byte(str)
}
