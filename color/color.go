package color

import (
	"strconv"
	"strings"
)

// nolint: varcheck, deadcode, megacheck
const (
	Reset        = "\033[0m"
	Normal       = 0
	Bold         = 1
	Dim          = 2
	Underline    = 4
	Blink        = 5
	Reverse      = 7
	Hidden       = 8
	Black        = 30
	Red          = 31
	Green        = 32
	Yellow       = 33
	Blue         = 34
	Magenta      = 35
	Cyan         = 36
	BrightGray   = 37
	DarkGray     = 90
	BrightRed    = 91
	BrightGreen  = 92
	BrightYellow = 93
	BrightBlue   = 94
	BrightPurple = 95
	BrightCyan   = 96
	White        = 97
)

func Color(code int, flags ...int) string {
	var strFlags []string
	for _, f := range flags {
		strFlags = append(strFlags, strconv.Itoa(f))
	}
	return "\033[" + strings.Join(strFlags, "") + ";" + strconv.Itoa(code) + "m"
}

func ColorNormal(code int) string {
	return "\033[" + strconv.Itoa(Normal) + ";" + strconv.Itoa(code) + "m"
}

func ColorBold(code int) string {
	return "\033[" + strconv.Itoa(Bold) + ";" + strconv.Itoa(code) + "m"
}
