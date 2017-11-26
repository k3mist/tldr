package main

import (
	"fmt"

	"bitbucket.org/djr2/tldr/color"
	"bitbucket.org/djr2/tldr/config"
)

func banner() {
	cfg := config.Config
	fmt.Print("" +
		color.Color(cfg.BannerColor1) + `   ___________   _____  _____  ` + "\n" +
		color.Color(cfg.BannerColor2) + `  /__   __/  /  /  _  \/  _  \ ` + "\n" +
		color.Color(cfg.BannerColor2) + `    /  / /  /  /  //  /  //  / ` + "\n" +
		color.Color(cfg.BannerColor1) + `   /  / /  /__/  //  /  / \  \ ` + "\n" +
		color.Color(cfg.BannerColor2) + `  /__/ /_____/______/__/   \_/ ` + color.ColorBold(cfg.TLDRColor) + "tldr.sh\n\n" + color.Reset,
	)
}
