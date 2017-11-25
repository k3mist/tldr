package page

import (
	"bytes"
	"fmt"
	"regexp"

	"bitbucket.org/djr2/tldr/color"
)

var (
	headerRxV1  = regexp.MustCompile(`^#\s`)
	exampleRxV1 = regexp.MustCompile(`^(-\s)`)
	codeRxV1    = regexp.MustCompile("^`(.+)`$")
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
	return headerRxV1.ReplaceAll(p.lines[0], to_b(color.ColorBold(color.White)+"["+color.ColorBold(color.Blue)))
}

func (p *pagev1) example(line []byte) []byte {
	if exampleRxV1.Match(line) {
		return exampleRxV1.ReplaceAll(line, to_b(color.Color(color.Normal)+"$1"+color.ColorNormal(color.Cyan)))
	}
	return nil
}

func (p *pagev1) code(line []byte) []byte {
	if codeRxV1.Match(line) {
		return codeRxV1.ReplaceAll(line, to_b(color.ColorNormal(color.Red)+"$1"))
	}
	return nil
}
