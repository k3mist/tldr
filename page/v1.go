package page

import (
	"bytes"
	"fmt"
	"regexp"

	"bitbucket.org/djr2/tldr/color"
	"bitbucket.org/djr2/tldr/config"
)

var (
	headerRxV1  = regexp.MustCompile(`^#\s`)
	exampleRxV1 = regexp.MustCompile(`^(-\s)`)
	syntaxRxV1  = regexp.MustCompile("^`(.+)`$")
)

type pagev1 struct {
	lines [][]byte
	buf   *bytes.Buffer
}

func (p *pagev1) Lines() [][]byte {
	return p.lines
}

func (p *pagev1) Write(b []byte) {
	p.buf.Write(b)
}

func (p *pagev1) Print() {
	fmt.Println(p.buf.String() + color.Reset)
}

func (p *pagev1) header() []byte {
	cfg := config.Config
	return headerRxV1.ReplaceAll(p.lines[0], to_b(color.ColorBold(cfg.HeaderDecorColor)+"["+color.ColorBold(cfg.HeaderColor)))
}

func (p *pagev1) example(line []byte) []byte {
	if exampleRxV1.Match(line) {
		cfg := config.Config
		return exampleRxV1.ReplaceAll(line, to_b(color.Color(cfg.HypenColor)+"$1"+color.Color(cfg.ExampleColor)))
	}
	return nil
}

func (p *pagev1) syntax(line []byte) []byte {
	if syntaxRxV1.Match(line) {
		cfg := config.Config
		return syntaxRxV1.ReplaceAll(line, to_b(color.Color(cfg.SyntaxColor)+"$1"))
	}
	return nil
}
