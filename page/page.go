package page

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	"bitbucket.org/djr2/tldr/color"
)

var (
	header      = regexp.MustCompile(`^#\s`)
	description = regexp.MustCompile(`^>\s`)
	example     = regexp.MustCompile(`^-\s`)
	codeStart   = regexp.MustCompile(`^.([a-z])`)
	codeEnd     = regexp.MustCompile("`$")
	variable    = regexp.MustCompile(`{{([\w\s\\/~!@#$%^&*()\[\]:;"'<,>?.]+)}}`)
)

func Print(file *os.File) {
	b, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	contents := string(b)
	buf := bytes.Buffer{}
	for _, line := range strings.Split(contents, "\n") {
		var str string
		if header.MatchString(line) {
			str = header.ReplaceAllString(line, color.ColorBold(color.White)+"["+color.ColorBold(color.Blue))
			buf.WriteString("  " + str + color.ColorBold(color.White) + "]" + "\n")
			continue
		}
		if description.MatchString(line) {
			str = description.ReplaceAllString(line, color.Color(color.Normal))
			buf.WriteString("  " + str + "\n")
			continue
		}
		if example.MatchString(line) {
			str = example.ReplaceAllString(line, color.Color(color.Normal)+"- "+color.ColorNormal(color.Cyan))
			buf.WriteString("\n  " + str + "\n")
			continue
		}
		if codeStart.MatchString(line) {
			str = codeStart.ReplaceAllString(line, color.ColorNormal(color.Red)+"$1")
			str = codeEnd.ReplaceAllString(str, "")
			str = variable.ReplaceAllString(str, color.Color(color.Normal)+"$1"+color.ColorNormal(color.Red))
			buf.WriteString("    " + str + "\n")
			continue
		}

	}
	fmt.Println(buf.String())
}
