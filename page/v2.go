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
	syntaxRxV2  = regexp.MustCompile(`^[\s]{4}`)
)

type pagev2 struct {
	lns [][]byte
	buf *bytes.Buffer
}

func (p *pagev2) Lines() [][]byte {
	if headerRxV2.Match(p.lns[1]) {
		p.lns[1] = p.lns[0]
		p.lns = p.lns[1:]
	}
	return p.lns
}

func (p *pagev2) Write(b []byte) {
	p.buf.Write(b)
}

func (p *pagev2) Print() {
	fmt.Println(p.buf.String() + color.Reset)
}

func (p *pagev2) Header() []byte {
	cfg := config.Config
	return append(toB(color.ColorBold(cfg.HeaderColor)), p.lns[0]...)
}

func (p *pagev2) Description(line []byte) []byte {
	cfg := config.Config
	if descRx.Match(line) {
		return descRx.ReplaceAll(line, toB(color.Color(cfg.DescriptionColor)))
	}
	return nil
}

func (p *pagev2) Example(line []byte) []byte {
	if exampleRxV2.Match(line) {
		cfg := config.Config
		return exampleRxV2.ReplaceAll(line, toB(color.Color(cfg.HypenColor)+"- $1"+color.Color(cfg.ExampleColor)))
	}
	return nil
}

func (p *pagev2) Syntax(line []byte) []byte {
	if syntaxRxV2.Match(line) {
		cfg := config.Config
		return syntaxRxV2.ReplaceAll(line, toB(color.Color(cfg.SyntaxColor)))
	}
	return nil
}

func (p *pagev2) Variable(line []byte) []byte {
	cfg := config.Config
	return varRx.ReplaceAll(line, toB(color.Color(cfg.VariableColor)+"$1"+color.Color(cfg.SyntaxColor)))
}
