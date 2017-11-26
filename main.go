package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"bitbucket.org/djr2/tldr/cache"
	"bitbucket.org/djr2/tldr/config"
	"bitbucket.org/djr2/tldr/page"
	"bitbucket.org/djr2/tldr/platform"
)

var flagSet *flag.FlagSet

func init() {
	flagSet = flag.NewFlagSet("", flag.ContinueOnError)
	flagSet.String("p", "", "platform of the tldr page\n\t  `platform` -- "+
		strings.Join(platform.Platforms(), ", "))
	flagSet.String("c", "", "clear cache for a tldr page\n\t  `page` -- "+
		"Use `clearall` to clear entire cache\n\t  -p is required if clearing cache for a specific platform")
	flagSet.String("debug", "disable", "enables debug logging")
	log.SetOutput(new(logWriter))
}

func main() {
	tldr()
}

func tldr() {
	if err := flagSet.Parse(os.Args[1:]); err != nil {
		return
	}

	setLogDebug()

	config.Load()
	plat := platform.ParseFlag(flagSet.Lookup("p"))

	if clear := flagSet.Lookup("c"); clear.Value.String() != "" {
		banner()
		cache.Remove(clear.Value.String(), plat)
		return
	}

	if len(flagSet.Args()) > 0 {
		page.New(cache.Find(flagSet.Arg(0), plat)).Print()
	} else if len(os.Args[1:]) == 0 {
		banner()
		flagSet.Usage()
	}
}
