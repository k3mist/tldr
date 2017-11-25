package platform

import (
	"flag"
	"runtime"
)

type Platform uint32

const (
	UNKNOWN Platform = iota
	COMMON
	LINUX
	OSX
	SUNOS
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

func (p Platform) String() string {
	if name, ok := platformMap[p]; ok {
		return name
	}
	return UNKNOWN.String()
}

func ParseFlag(p *flag.Flag) Platform {
	return Parse(p.Value.String())
}

func Parse(p string) Platform {
	for plat, name := range platformMap {
		if p == name {
			return plat
		}
	}
	return UNKNOWN
}

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
