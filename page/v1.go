package page

import (
	"bytes"
	"fmt"
	"regexp"

	"bitbucket.org/djr2/tldr/color"
)

var (
	headerRxV1  = regexp.MustCompile(`^#\s`)
	exampleRxV1 = regexp.MustCompile(`^-\s`)
	codeRxV1    = regexp.MustCompile("^`(.+)`$")
)

type pagev1 struct {
	lines []string
	buf   *bytes.Buffer
}

func (p *pagev1) Print() {
	p.buf.WriteString("\n")
	for i, line := range p.lines {
		if i == 0 {
			p.buf.WriteString("  " + p.header() + color.ColorBold(color.White) + "]" + "\n")
			continue
		}

		if desc := description(line); desc != "" {
			p.buf.WriteString("  " + desc + "\n")
			continue
		}

		if example := p.example(line); example != "" {
			p.buf.WriteString("\n  " + example + "\n")
			continue
		}

		if code := p.code(line); code != "" {
			p.buf.WriteString("    " + variable(code) + "\n")
			continue
		}
	}
	fmt.Println(p.buf.String() + color.Reset)
}

func (p *pagev1) header() string {
	return headerRxV1.ReplaceAllString(p.lines[0], color.ColorBold(color.White)+"["+color.ColorBold(color.Blue))
}

func (p *pagev1) example(line string) string {
	if exampleRxV1.MatchString(line) {
		return exampleRxV1.ReplaceAllString(line, color.Color(color.Normal)+"- "+color.ColorNormal(color.Cyan))
	}
	return ""
}

func (p *pagev1) code(line string) string {
	if codeRxV1.MatchString(line) {
		return codeRxV1.ReplaceAllString(line, color.ColorNormal(color.Red)+"$1")
	}
	return ""
}
