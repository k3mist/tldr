package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"bitbucket.org/djr2/tldr/cache"
	"bitbucket.org/djr2/tldr/color"
	"bitbucket.org/djr2/tldr/page"
	"bitbucket.org/djr2/tldr/platform"
)

var flagSet *flag.FlagSet

// Custom time format
type logWriter struct{}

// Write with custom time format.
func (writer logWriter) Write(bytes []byte) (int, error) {
	if flagSet.Lookup("debug").Value.String() != "disable" {
		return fmt.Print("["+time.Now().UTC().Format("2006-01-02 15:04:05")+"] ", string(bytes))
	}
	return fmt.Print(string(bytes))
}

func init() {
	flagSet = flag.NewFlagSet("", flag.ContinueOnError)
	flagSet.String("p", "common", "Man page platform.\n\t  `platform` -- "+
		strings.Join(platform.Platforms(), ", "))
	flagSet.String("c", "", "Clear cache for a man page.\n\t  `name` -- "+
		"Use `clearcache` to clear entire cache.\n\t  -p is required if clearing cache for a specific platform.")
	flagSet.String("debug", "disable", "Enables debugging.")
	log.SetOutput(new(logWriter))
}

func main() {
	if err := tldr(); err != nil {
		os.Exit(1)
	}
}

func tldr() error {
	fmt.Print(Banner())

	if err := flagSet.Parse(os.Args[1:]); err != nil {
		return nil
	}

	if flagSet.Lookup("debug").Value.String() != "disable" {
		log.SetFlags(log.Lshortfile)
	} else {
		log.SetFlags(0)
	}

	if clear := flagSet.Lookup("c"); clear.Value.String() != "" {
		cache.Remove(clear.Value.String(), platform.Parse(flagSet.Lookup("p")))
		return nil
	}

	if len(flagSet.Args()) > 0 {
		page.Print(cache.Find(flagSet.Arg(0), platform.Parse(flagSet.Lookup("p"))))
	} else if len(os.Args[1:]) == 0 {
		flagSet.Usage()
	}

	return nil
}

func Banner() string {
	return "" +
		color.ColorNormal(color.Blue) + `   ___________   _____  _____  ` + "\n" +
		color.ColorNormal(color.Cyan) + `  /__   __/  /  /  _  \/  _  \ ` + "\n" +
		color.ColorNormal(color.Cyan) + `    /  / /  /  /  //  /  //  / ` + "\n" +
		color.ColorNormal(color.Blue) + `   /  / /  /__/  //  /  / \  \ ` + "\n" +
		color.ColorNormal(color.Cyan) + `  /__/ /_____/______/__/   \_/ ` +
		color.ColorBold(color.White) + "https://tldr.sh\n\n" + color.Reset
}
