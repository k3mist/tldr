package platform

import (
	"flag"
	"runtime"
)

// Platform holds the available TLDR page platforms.
type Platform uint32

const (
	// UNKNOWN is for unknown platform type
	UNKNOWN Platform = iota
	// COMMON is a reference to the common directory of the TLDR page assets.
	COMMON
	// LINUX is a reference to the linux directory of the TLDR page assets.
	LINUX
	// OSX is a reference to the osx directory of the TLDR page assets.
	OSX
	// SUNOS is a reference to the sunos directory of the TLDR page assets.
	SUNOS
	// WINDOWS is a reference to the windows directory of the TLDR page assets.
	WINDOWS
)

var platformMap = map[Platform]string{
	UNKNOWN: `unknown`,
	COMMON:  `common`,
	LINUX:   `linux`,
	OSX:     `osx`,
	SUNOS:   `sunos`,
	WINDOWS: `windows`,
}

// Strings provides the string based representation of the Platform
func (p Platform) String() string {
	if name, ok := platformMap[p]; ok {
		return name
	}
	return UNKNOWN.String()
}

// ParseFlag parses the provided command line platform Flag
func ParseFlag(p *flag.Flag) Platform {
	return Parse(p.Value.String())
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

// Actual returns the runtime platform. If a valid platform is not found it
// it returns COMMON.
func Actual() Platform {
	switch runtime.GOOS {
	case `freebsd`, `netbsd`, `openbsd`, `plan9`, `linux`:
		return LINUX
	case `darwin`:
		return OSX
	case `solaris`:
		return SUNOS
	case `windows`:
		return WINDOWS
	}
	return COMMON
}
