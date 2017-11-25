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

func (p *pagev1) Print() {
	p.buf.WriteString("\n")
	for i, line := range p.lines {
		if i == 0 {
			p.buf.Write(to_b("  "))
			p.buf.Write(p.header())
			p.buf.Write(to_b(color.ColorBold(color.White) + "]" + "\n"))
			continue
		}

		if desc := description(line); desc != nil {
			p.buf.Write(to_b("  "))
			p.buf.Write(desc)
			p.buf.Write(to_b("\n"))
			continue
		}

		if example := p.example(line); example != nil {
			p.buf.Write(to_b("\n  "))
			p.buf.Write(example)
			p.buf.Write(to_b("\n"))
			continue
		}

		if code := p.code(line); code != nil {
			p.buf.Write(to_b("    "))
			p.buf.Write(variable(code))
			p.buf.Write(to_b("\n"))
			continue
		}
	}
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
