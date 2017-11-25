package page

import (
	"bytes"
	"fmt"
	"regexp"

	"bitbucket.org/djr2/tldr/color"
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
	return append(to_b(color.ColorBold(color.White)+"["+color.ColorBold(color.Blue)), p.lines[0]...)
}

func (p *pagev2) example(line []byte) []byte {
	if exampleRxV2.Match(line) {
		return exampleRxV2.ReplaceAll(line, to_b(color.Color(color.Normal)+"- $1"+color.ColorNormal(color.Cyan)))
	}
	return nil
}

func (p *pagev2) code(line []byte) []byte {
	if codeRxV2.Match(line) {
		return codeRxV2.ReplaceAll(line, to_b(color.ColorNormal(color.Red)))
	}
	return nil
}
