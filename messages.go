package main

import (
	"fmt"

	"bitbucket.org/djr2/tldr/color"
	"bitbucket.org/djr2/tldr/config"
	"bitbucket.org/djr2/tldr/platform"
)

func banner() {
	cfg := config.Config
	fmt.Print("" +
		color.Color(cfg.BannerColor1) + `   ___________   _____  _____  ` + "\n" +
		color.Color(cfg.BannerColor2) + `  /__   __/  /  /  _  \/  _  \ ` + "\n" +
		color.Color(cfg.BannerColor2) + `    /  / /  /  /  //  /  //  / ` + "\n" +
		color.Color(cfg.BannerColor1) + `   /  / /  /__/  //  /  / \  \ ` + "\n" +
		color.Color(cfg.BannerColor2) + `  /__/ /_____/______/__/   \_/ ` + color.ColorBold(cfg.TLDRColor) + "https://tldr.sh\n\n" + color.Reset,
	)
}

func version() {
	banner()
	cfg := config.Config
	fmt.Print(
		color.Color(cfg.BannerColor1) + `  version: ` + color.ColorBold(cfg.TLDRColor) + "2.1.0" + color.Reset + "\n",
	)
}

func noTldr(cmd string) {
	cfg := config.Config
	fmt.Println(color.Color(cfg.DescriptionColor) + "> Unable to find a tldr for " + color.Color(cfg.HeaderColor) + cmd + color.Reset)
}

func noLookup(cmd string, lang string, plat platform.Platform) {
	cfg := config.Config
	if cfg.LookupWarnings {
		fmt.Print(color.Color(cfg.DescriptionColor) + "> No tldr for ")
		fmt.Print(color.Color(cfg.HeaderDecorColor) + "[")
		fmt.Print(color.ColorBold(cfg.HeaderColor) + cmd)
		fmt.Print(color.Color(cfg.HeaderDecorColor) + "] language[")
		fmt.Print(color.ColorBold(cfg.HeaderColor) + lang)
		fmt.Print(color.Color(cfg.HeaderDecorColor) + "] platform[")
		fmt.Print(color.ColorBold(cfg.HeaderColor) + plat.String())
		fmt.Print(color.Color(cfg.HeaderDecorColor) + "]")
		fmt.Println(color.Reset)
	}
}
