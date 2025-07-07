package platform

import (
	// "flag"
	"bytes"
	"fmt"
	"runtime"

	"bitbucket.org/djr2/tldr/color"
	"bitbucket.org/djr2/tldr/config"
)

// Platform holds the available TLDR page platforms.
type Platform uint32

const (
	// UNKNOWN is for unknown platform type
	UNKNOWN Platform = iota
	// COMMON is a reference to the common directory of the TLDR page assets.
	COMMON
	// OSX is a reference to the osx directory of the TLDR page assets.
	OSX
	// WINDOWS is a reference to the windows directory of the TLDR page assets.
	WINDOWS
	// LINUX is a reference to the linux directory of the TLDR page assets.
	LINUX
	// ANDROID is a reference to the android directory of the TLDR page assets.
	ANDROID
	// FREEBSD is a reference to the freebsd directory of the TLDR page assets.
	FREEBSD
	// OPENBSD is a reference to the openbsd directory of the TLDR page assets.
	OPENBSD
	// NETBSD is a reference to the netbsd directory of the TLDR page assets.
	NETBSD
	// SUNOS is a reference to the sunos directory of the TLDR page assets.
	SUNOS
)

var platformMap = map[Platform]string{
	UNKNOWN: `unknown`,
	COMMON:  `common`,
	OSX:     `osx`,
	WINDOWS: `windows`,
	LINUX:   `linux`,
	ANDROID: `android`,
	FREEBSD: `freebsd`,
	OPENBSD: `openbsd`,
	NETBSD:  `netbsd`,
	SUNOS:   `sunos`,
}

// Get the platform
func GetPlatform(plat Platform) Platform {
	if plat == UNKNOWN {
		return Actual()
	}
	return plat
}

// String provides the string based representation of the Platform
func (p Platform) String() string {
	if name, ok := platformMap[p]; ok {
		return name
	}
	return UNKNOWN.String()
}

// ParseFlag parses the provided command line platform Flag
func ParseFlag(p string) Platform {
	return Parse(p)
}

// Parse returns the Platform if valid. If the provided platform is not valid
// it returns an UNKNOWN platform.
func Parse(p string) Platform {
	for plat, name := range platformMap {
		if p == name {
			return plat
		}
	}
	return UNKNOWN
}

// Platforms returns the string array of the available Platforms.
func Platforms() []string {
	var platforms []string
	for _, name := range platformMap {
		if name == `unknown` {
			continue
		}
		platforms = append(platforms, name)
	}
	return platforms
}

// Print prints the platforms
func Print() {
	var buf bytes.Buffer
	cfg := config.Config

	buf.Write(toB("\n"))
	buf.Write(toB("  "))
	buf.Write(toB(color.ColorBold(cfg.HeaderDecorColor) + "["))
	buf.Write(toB(color.ColorBold(cfg.HeaderColor) + "platforms"))
	buf.Write(toB(color.ColorBold(cfg.HeaderDecorColor) + "]\n"))

	buf = write(Actual().String(), true, buf)

	for _, name := range platformMap {
		if name == `unknown` {
			continue
		}

		if name == Actual().String() {
			continue
		}

		buf = write(name, false, buf)
	}

	fmt.Print(buf.String() + color.Reset)
}

func write(name string, actual bool, buf bytes.Buffer) bytes.Buffer {
	cfg := config.Config
	buf.Write(toB("  "))
	if actual {
		buf.Write(toB(color.Color(cfg.HyphenColor) + "> "))
		buf.Write(toB(color.Color(cfg.SyntaxColor) + name))
		buf.Write(toB(color.Color(cfg.DescriptionColor) + " (detected)"))
	} else {
		buf.Write(toB(color.Color(cfg.HyphenColor) + "- "))
		buf.Write(toB(color.Color(cfg.DescriptionColor) + name))
	}
	buf.Write(toB("\n"))
	return buf
}

// Actual returns the runtime platform. If a valid platform is not found it
// it returns COMMON.
func Actual() Platform {
	switch runtime.GOOS {
	case `linux`:
		return LINUX
	case `freebsd`:
		return FREEBSD
	case `netbsd`:
		return NETBSD
	case `openbsd`:
		return OPENBSD
	case `darwin`:
		return OSX
	case `solaris`:
		return SUNOS
	case `windows`:
		return WINDOWS
	}
	return COMMON
}

func toB(str string) []byte {
	return []byte(str)
}
