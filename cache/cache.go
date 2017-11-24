package cache

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"strings"

	"bitbucket.org/djr2/tldr/platform"
	"github.com/mitchellh/go-homedir"
)

var cacheDir string

const repository = "https://raw.github.com/tldr-pages/tldr/master/pages/"

func init() {
	h, err := homedir.Dir()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	cacheDir = h + "/" + ".tldr"

	if _, err := os.Stat(cacheDir); err != nil {
		if err := os.Mkdir(cacheDir, 0700); err != nil {
			log.Println(err)
			os.Exit(1)
		}
	}
}

func newCacher(name string, plat platform.Platform) *cacher {
	return &cacher{name: name + ".md", platform: plat.String()}
}

func Find(name string, plat platform.Platform) *os.File {
	cacher := newCacher(name, plat)
	cached := cacher.search()
	if cached != nil {
		return cached
	}
	return cacher.create()
}

func Remove(name string, plat platform.Platform) {
	cacher := newCacher(name, plat)
	cacher.remove()
}

type cacher struct {
	platform string
	name     string
}

func (c *cacher) platformDir() string {
	return cacheDir + "/" + c.platform
}

func (c *cacher) file() string {
	return c.platformDir() + "/" + c.name
}

func (c *cacher) search() *os.File {
	for _, fileInfo := range c.readDir() {
		if fileInfo.Name() == c.name {
			file, err := os.Open(c.file())
			if err != nil {
				log.Println(err)
				os.Exit(1)
			}
			return file
		}
	}
	return nil
}

func (c *cacher) create() *os.File {
	url := repository + c.platform + "/" + c.name
	log.Println("Retrieving:", url)
	response, err := http.Get(url)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	if response.StatusCode != http.StatusOK {
		log.Println("Server Error:", response.StatusCode)
		os.Exit(1)
	}

	buf, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	file, err := os.Create(c.file())
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	ret, err := file.Write(buf)
	defer file.Close()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	log.Println("Created:", c.file(), "bytes:", strconv.Itoa(ret), "\n")
	return c.search()
}

func (c *cacher) remove() {
	if c.name == "clearcache.md" {
		if err := os.RemoveAll(cacheDir); err != nil {
			log.Println(err)
			os.Exit(1)
		}
		log.Println("Cache cleared")
		os.Exit(0)
	}

	file := c.search()
	if file == nil {
		log.Println("Command:", strings.TrimRight(c.name, ".md"), "not cached")
		os.Exit(1)
	}

	if err := os.Remove(c.file()); err != nil {
		log.Println(err)
		os.Exit(1)
	}

	log.Println("Removed:", strings.TrimRight(c.name, ".md"), c.file())
	os.Exit(0)
}

func (c *cacher) readDir() []os.FileInfo {
	_, err := os.Stat(c.platformDir())
	if err != nil {
		if err := os.Mkdir(c.platformDir(), 0700); err != nil {
			log.Println(err)
			os.Exit(1)
		}
	}

	srcDir, err := ioutil.ReadDir(c.platformDir())
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	return srcDir
}
