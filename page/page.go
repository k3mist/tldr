package page

import (
	"bytes"
	"io"
	"log"
	"os"
	"regexp"

	"bitbucket.org/djr2/tldr/color"
	"bitbucket.org/djr2/tldr/config"
	"bitbucket.org/djr2/tldr/platform"
)

// Page provides the Print method for the final generated output of a TLDR page.
type Page interface {
	Print()
}

// Parser provides the interface for parsing a TLDR page.
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

// New creates a parsed TLDR page. It parsers the provided file and returns the
// parsed TLDR page.
func New(file *os.File, plat platform.Platform) Page { // nolint: interfacer
	b, err := io.ReadAll(file)
	defer file.Close() // nolint: errcheck
	if err != nil {
		log.Fatal(err)
	}
	lines := bytes.Split(b, toB("\n"))
	if headerRxV2.Match(lines[1]) {
		p := &pagev2{lines, &bytes.Buffer{}}
		Parse(p, plat)
		return p
	}
	p := &pagev1{lines, &bytes.Buffer{}}
	Parse(p, plat)
	return p
}

// Parse takes a Parser interface and the current platform of the document
// that is to be parsed and the parses the internal lines writing to the
// internal buffer of the parser.
func Parse(p Parser, plat platform.Platform) {
	cfg := config.Config
	p.Write(toB("\n"))
	for i, line := range p.Lines() {
		if i == 0 {
			p.Write(toB("  "))
			p.Write(toB(color.ColorBold(cfg.HeaderDecorColor) + "["))
			p.Write(p.Header())
			p.Write(toB(color.ColorBold(cfg.HeaderDecorColor) + " - "))
			if plat == platform.Actual() || plat == platform.COMMON {
				p.Write(toB(color.Color(cfg.PlatformColor) + plat.String()))
			} else {
				p.Write(toB(color.Color(cfg.PlatformAltColor) + plat.String()))
			}
			p.Write(toB(color.ColorBold(cfg.HeaderDecorColor) + "]\n"))
			continue
		}

		if desc := p.Description(line); desc != nil {
			p.Write(toB("  "))
			p.Write(desc)
			p.Write(toB("\n"))
			continue
		}

		if example := p.Example(line); example != nil {
			p.Write(toB("\n  "))
			p.Write(example)
			p.Write(toB("\n"))
			continue
		}

		if syntax := p.Syntax(line); syntax != nil {
			p.Write(toB("    "))
			p.Write(p.Variable(syntax))
			p.Write(toB("\n"))
			continue
		}
	}
}

func toB(str string) []byte {
	return []byte(str)
}
