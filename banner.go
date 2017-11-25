package main

import (
	"fmt"

	"bitbucket.org/djr2/tldr/color"
)

func banner() {
	fmt.Print("" +
		color.Color(color.Blue) + `   ___________   _____  _____  ` + "\n" +
		color.Color(color.Cyan) + `  /__   __/  /  /  _  \/  _  \ ` + "\n" +
		color.Color(color.Cyan) + `    /  / /  /  /  //  /  //  / ` + "\n" +
		color.Color(color.Blue) + `   /  / /  /__/  //  /  / \  \ ` + color.ColorBold(color.White) + "tldr.sh\n" +
		color.Color(color.Cyan) + `  /__/ /_____/______/__/   \_/ ` + color.ColorBold(color.DarkGray) + "bitbucket.org/djr2/tldr\n\n" + color.Reset,
	)
}
