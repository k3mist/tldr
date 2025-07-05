package cache

import (
	"log"
	"os"

	"github.com/mitchellh/go-homedir"

	"bitbucket.org/djr2/tldr/platform"
)

var cacheDir string

func init() {
	h, err := homedir.Dir()
	if err != nil {
		log.Fatal(err)
	}

	cacheDir = h + "/" + ".tldr"

	if _, err := os.Stat(cacheDir); err != nil {
		if err := os.Mkdir(cacheDir, 0700); err != nil {
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

// Find attempts to find the requested tldr page from the local cache. If a
// local cache page is not found it will attempt to retrieve the page from
// tldr pages repository
func Find(name string, plat platform.Platform) (*os.File, platform.Platform) {
	cacher := newCacher(name, validPlatform(plat))
	return cacher.search(), cacher.plat
}

// Remove will delete a local tldr page from the cache or if `clearall` is
// provided as the name it will remove all tldr pages from the cache.
func Remove(name string, plat platform.Platform) {
	cacher := newCacher(name, validPlatform(plat))
	cacher.remove()
}
