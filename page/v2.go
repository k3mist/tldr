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
	lines  []string
	buf    *bytes.Buffer
	parser parser
}

func (p *pagev2) Print() {
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

func (p *pagev2) header() string {
	return color.ColorBold(color.White) + "[" + color.ColorBold(color.Blue) + p.lines[0]
}

func (p *pagev2) example(line string) string {
	if exampleRxV2.MatchString(line) {
		return color.Color(color.Normal) + "- " + color.ColorNormal(color.Cyan) + line
	}
	return ""
}

func (p *pagev2) code(line string) string {
	if codeRxV2.MatchString(line) {
		return codeRxV2.ReplaceAllString(line, color.ColorNormal(color.Red))
	}
	return ""
}
