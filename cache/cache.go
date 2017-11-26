package cache

import (
	"log"
	"os"

	"bitbucket.org/djr2/tldr/platform"
	"github.com/mitchellh/go-homedir"
)

var CacheDir string

func init() {
	h, err := homedir.Dir()
	if err != nil {
		log.Fatal(err)
	}

	CacheDir = h + "/" + ".tldr"

	if _, err := os.Stat(CacheDir); err != nil {
		if err := os.Mkdir(CacheDir, 0700); err != nil {
			log.Fatal(err)
		}
	}

	getAssets()
}

func newCacher(name string, plat platform.Platform) *cacher {
	return &cacher{name: name + ".md", plat: plat}
}

func validPlatform(plat platform.Platform) platform.Platform {
	if plat == platform.UNKNOWN {
		return platform.Actual()
	}
	return plat
}

func Find(name string, plat platform.Platform) (*os.File, platform.Platform) {
	cacher := newCacher(name, validPlatform(plat))
	return cacher.search(), cacher.plat
}

func Remove(name string, plat platform.Platform) {
	cacher := newCacher(name, validPlatform(plat))
	cacher.remove()
}
