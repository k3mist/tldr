package main

import (
	"flag"
	"os"
	"strings"

	"github.com/rsc/getopt"

	"bitbucket.org/djr2/tldr/cache"
	"bitbucket.org/djr2/tldr/config"
	"bitbucket.org/djr2/tldr/page"
	"bitbucket.org/djr2/tldr/platform"
)

var flagSet *getopt.FlagSet

var flagClear bool
var flagPageClear string
var flagPlatform string
var flagPlatforms bool
var flagLanguage string
var flagVersion bool

func init() {
	flagSet = getopt.NewFlagSet("", flag.ContinueOnError)

	flagSet.BoolVar(&flagClear, "clear", false, "Clear the entire page cache")

	flagSet.StringVar(&flagPageClear, "c", "", "Clear cache for a specific tldr `page`\n"+
		"\t-p is required if clearing cache for a specific platform")

	flagSet.StringVar(&flagPlatform, "p", "", "Platform of the desired tldr page.")
	flagSet.Alias("p", "platform")

	flagSet.BoolVar(&flagPlatforms, "platforms", false, "Display a list of available platforms")

	// TODO add language support
	flagSet.StringVar(&flagLanguage, "L", "", "The desired language for the tldr page. (WIP)")
	flagSet.Alias("L", "language")

	flagSet.BoolVar(&flagVersion, "version", false, "Display the version number")

	flagSet.String("debug", "disable", "Enables debug logging")
	flagSet.SetOutput(new(logWriter))
}

func main() {
	config.Load()

	if len(os.Args[1:]) == 0 {
		banner()
		flagSet.Usage()
		return
	}

	tldr()
}

func tldr() {
	setLogDebug()

	var cmd string = os.Args[1]
	var args []string

	if strings.HasPrefix(cmd, "-") {
		args = os.Args[1:]
	} else {
		args = os.Args[2:]
	}

	if err := flagSet.Parse(args); err != nil {
		return
	}

	if flagPlatforms {
		platform.Print()
		return
	}

	platform := platform.ParseFlag(flagPlatform)

	if flagClear {
		banner()
		cache.Remove("clearall", platform)
		return
	}

	if flagPageClear != "" {
		banner()
		cache.Remove(flagPageClear, platform)
		return
	}

	if len(os.Args[1:]) > 0 {
		page.New(cache.Find(cmd, platform)).Print()
	}
}
