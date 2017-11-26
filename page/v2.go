package page

import (
	"bytes"
	"fmt"
	"regexp"

	"bitbucket.org/djr2/tldr/color"
	"bitbucket.org/djr2/tldr/config"
)

var (
	headerRxV2  = regexp.MustCompile(`^[=]+$`)
	exampleRxV2 = regexp.MustCompile(`^([\w]+)`)
	codeRxV2    = regexp.MustCompile(`^[\s]{4}`)
)

type pagev2 struct {
	lines [][]byte
	buf   *bytes.Buffer
}

func (p *pagev2) Lines() [][]byte {
	return p.lines
}

func (p *pagev2) Write(b []byte) {
	p.buf.Write(b)
}

func (p *pagev2) Print() {
	fmt.Println(p.buf.String() + color.Reset)
}

func (p *pagev2) header() []byte {
	cfg := config.Config
	return append(to_b(color.ColorBold(cfg.HeaderDecorColor)+"["+color.ColorBold(cfg.HeaderColor)), p.lines[0]...)
}

func (p *pagev2) example(line []byte) []byte {
	if exampleRxV2.Match(line) {
		cfg := config.Config
		return exampleRxV2.ReplaceAll(line, to_b(color.Color(cfg.HypenColor)+"- $1"+color.Color(cfg.ExampleColor)))
	}
	return nil
}

func (p *pagev2) code(line []byte) []byte {
	if codeRxV2.Match(line) {
		cfg := config.Config
		return codeRxV2.ReplaceAll(line, to_b(color.Color(cfg.SyntaxColor)))
	}
	return nil
}
