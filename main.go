package main

import (
	"flag"
	"os"
	"strings"

	"rsc.io/getopt"

	"bitbucket.org/djr2/tldr/cache"
	"bitbucket.org/djr2/tldr/config"
	"bitbucket.org/djr2/tldr/language"
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
var flagGet bool
var flagExtended bool
var flagWarn bool
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

	flagSet.StringVar(&flagLanguage, "L", "", "The desired language for the tldr page.")
	flagSet.Alias("L", "language")

	flagSet.BoolVar(&flagVersion, "version", false, "Display the version number.")

	flagSet.BoolVar(&flagGet, "g", false, "If a tldr page is not cached attempt to retrieve it.")
	flagSet.Alias("g", "get")

	// flagSet.BoolVar(&flagExtended, "e", false, "Perform an extended search. (default: on via config)")
	// flagSet.Alias("e", "extended")

	// flagSet.BoolVar(&flagWarn, "w", false, "Show search warnings when page, language, platform combination is not found.")
	// flagSet.Alias("w", "warn")

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
	language.SetLang()

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

	plat := platform.ParseFlag(flagPlatform)

	var lang string = flagLanguage
	if lang == "" {
		lang = language.GetLanguage(0)
	}

	if flagUpdate {
		banner()
		cache.Remove("clearall", lang, plat, false)
		cache.Create()
		return
	}

	if flagClear {
		banner()
		cache.Remove("clearall", lang, plat, true)
		return
	}

	if flagPageClear != "" {
		banner()
		cache.Remove(flagPageClear, lang, plat, true)
		return
	}

	getTldr(cmd, lang, plat)
}

func getTldr(cmd string, lang string, plat platform.Platform) {
	if cmd != "" {
		var lfile *os.File
		var llang string
		var lplat platform.Platform

		lfile, llang, lplat = cache.Find(cmd, lang, plat, flagGet)
		if lfile == nil {
			noLookup(cmd, lang, plat)
			for _, l := range language.GetLanguages() {
				if l != lang {
					lfile, llang, lplat = cache.Find(cmd, l, plat, flagGet)
				}
				if lfile != nil {
					break
				} else {
					noLookup(cmd, llang, lplat)
				}
			}
		}

		if lfile != nil {
			page.New(lfile, llang, lplat).Print()
		} else {
			noTldr(cmd)
		}

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
		if strings.HasPrefix(p, "--") {
			lastHyphen = i
			args = append(args, p)
		} else if strings.HasPrefix(p, "-") {
			lastHyphen = i
			args = append(args, p)
			if len(os.Args[1:]) > 1 {
				args = append(args, os.Args[(i+1)+1])
			}
		} else if lastHyphen < 0 {
			cmds = append(cmds, p)
		} else {
			lastHyphen = -1
		}
	}

	cmd = strings.Join(cmds, "-")
	return cmd, args
}
