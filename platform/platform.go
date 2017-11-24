package platform

import (
	"flag"
	"runtime"
)

type Platform uint32

const (
	COMMON Platform = iota
	LINUX
	OSX
	SUNOS
	WINDOWS
)

var platformMap = map[Platform]string{
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
	return COMMON.String()
}

func Parse(p *flag.Flag) Platform {
	for plat, name := range platformMap {
		if p.Value.String() == name {
			return plat
		}
	}
	return COMMON
}

func Platforms() []string {
	var platforms []string
	for _, name := range platformMap {
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
