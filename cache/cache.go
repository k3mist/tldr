package cache

import (
	"log"
	"os"

	"github.com/mitchellh/go-homedir"

	"bitbucket.org/djr2/tldr/platform"
)

var configDir string
var cacheDir string

// Create local cache assets
func Create() {
	h, err := homedir.Dir()
	if err != nil {
		log.Fatal(err)
	}

	configDir = h + "/" + ".tldr"
	cacheDir = configDir + "/cache"

	if _, err := os.Stat(cacheDir); err != nil {
		if err := os.MkdirAll(cacheDir, 0700); err != nil {
			log.Fatal(err)
		}
	}

	getAssets()
}

func init() {
	Create()
}

func newCacher(name string, lang string, plat platform.Platform) *cacher {
	return &cacher{name: name + ".md", lang: lang, plat: plat}
}

func getPlatform(plat platform.Platform) platform.Platform {
	if plat == platform.UNKNOWN {
		return platform.Actual()
	}
	return plat
}

// Find attempts to find the requested tldr page from the local cache. If a
// local cache page is not found it will attempt to retrieve the page from
// tldr pages repository
func Find(name string, lang string, plat platform.Platform) (*os.File, string, platform.Platform) {
	cacher := newCacher(name, lang, getPlatform(plat))
	return cacher.search(), lang, cacher.plat
}

// Remove will delete a local tldr page from the cache or if `clearall` is
// provided as the name it will remove all tldr pages from the cache.
func Remove(name string, lang string, plat platform.Platform, exit bool) {
	cacher := newCacher(name, lang, getPlatform(plat))
	cacher.remove(exit)
}
