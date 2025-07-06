package main

import (
	"flag"
	"os"
	"strings"

	"rsc.io/getopt"

	"bitbucket.org/djr2/tldr/cache"
	"bitbucket.org/djr2/tldr/config"
	"bitbucket.org/djr2/tldr/page"
	"bitbucket.org/djr2/tldr/platform"
)

var flagSet *getopt.FlagSet

var flagClear bool
var flagPageClear string
var flagUpdate bool
var flagPlatform string
var flagPlatforms bool
var flagLanguage string
var flagVersion bool
var flagHelp bool

func init() {
	flagSet = getopt.NewFlagSet("", flag.ContinueOnError)

	flagSet.BoolVar(&flagClear, "clear", false, "Clear the entire page cache.")

	flagSet.StringVar(&flagPageClear, "c", "", "Clear cache for a specific tldr `page`.\n"+
		"\t-p is required if clearing cache for a specific platform.")

	flagSet.BoolVar(&flagUpdate, "u", false, "Update the local page cache.")
	flagSet.Alias("u", "update")

	flagSet.StringVar(&flagPlatform, "p", "", "Platform of the desired tldr page.")
	flagSet.Alias("p", "platform")

	flagSet.BoolVar(&flagPlatforms, "platforms", false, "Display a list of available platforms.")

	// TODO add language support
	flagSet.StringVar(&flagLanguage, "L", "", "The desired language for the tldr page.")
	flagSet.Alias("L", "language")

	flagSet.BoolVar(&flagVersion, "version", false, "Display the version number.")

	flagSet.BoolVar(&flagHelp, "help", false, "This usage output.")

	flagSet.String("debug", "disable", "Enables debug logging.")
	flagSet.SetOutput(new(logWriter))
}

func usage() {
	banner()
	flagSet.Usage()
}

func main() {
	config.Load()

	if len(os.Args[1:]) == 0 {
		usage()
		return
	}

	tldr()
}

func tldr() {
	setLogDebug()

	var cmd, args = getCmdArgs()
	if err := flagSet.Parse(args); err != nil {
		return
	}

	if flagHelp {
		usage()
		return
	}

	if flagVersion {
		version()
		return
	}

	if flagPlatforms {
		platform.Print()
		return
	}

	platform := platform.ParseFlag(flagPlatform)
	var language string = flagLanguage
	if language == "" {
		language = config.Config.Language
	}

	if flagUpdate {
		banner()
		cache.Remove("clearall", language, platform, false)
		cache.Create()
		return
	}

	if flagClear {
		banner()
		cache.Remove("clearall", language, platform, true)
		return
	}

	if flagPageClear != "" {
		banner()
		cache.Remove(flagPageClear, language, platform, true)
		return
	}

	if cmd != "" {
		page.New(cache.Find(cmd, language, platform)).Print()
		return
	}

	usage()
}

func getCmdArgs() (string, []string) {
	var cmd string
	var cmds []string
	var args []string

	var lastHyphen int = -1
	for i, p := range os.Args[1:] {
		// println(p)
		if strings.HasPrefix(p, "--") {
			lastHyphen = i
			args = append(args, p)
		} else if strings.HasPrefix(p, "-") {
			lastHyphen = i
			args = append(args, p)
			args = append(args, os.Args[(i+1)+1])
		} else if lastHyphen < 0 {
			cmds = append(cmds, p)
		} else {
			lastHyphen = -1
		}
	}
	// println(strings.Join(cmds, "-"))
	cmd = strings.Join(cmds, "-")
	return cmd, args
}
