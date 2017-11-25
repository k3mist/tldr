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

func (p *pagev2) Print() {
	p.buf.WriteString("\n")
	for i, line := range p.lines {
		if i == 0 {
			p.buf.Write(to_b("  "))
			p.buf.Write(p.header())
			p.buf.Write(p.lines[0])
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

func (p *pagev2) header() []byte {
	return to_b(color.ColorBold(color.White) + "[" + color.ColorBold(color.Blue))
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
