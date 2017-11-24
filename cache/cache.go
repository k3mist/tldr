package cache

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"bitbucket.org/djr2/tldr/pages"
	"bitbucket.org/djr2/tldr/platform"
	"github.com/mitchellh/go-homedir"
)

var cacheDir string

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
		info, err := cached.Stat()
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		if info.ModTime().Add(time.Hour * 120).Before(time.Now()) {
			log.Println("Cache older than 5 days")
			return cacher.save()
		}
		return cached
	}
	cacher.platform = plat.String()
	return cacher.save()
}

func Remove(name string, plat platform.Platform) {
	cacher := newCacher(name, plat)
	cacher.remove()
}

type cacher struct {
	platform string
	name     string
	page     *pages.Pages
}

func (c *cacher) platformDir() string {
	return cacheDir + "/" + c.platform
}

func (c *cacher) file() string {
	return c.platformDir() + "/" + c.name
}

func (c *cacher) cmd() string {
	return strings.TrimSuffix(c.name, `.md`)
}

func (c *cacher) meta() string {
	return c.platformDir() + "/" + c.cmd() + ".json"
}

func (c *cacher) search() *os.File {
	cached := c.find()
	if cached == nil {
		c.platform = platform.Actual().String()
	}
	return c.find()
}

func (c *cacher) find() *os.File {
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

func (c *cacher) download() io.ReadCloser {
	c.page = &pages.Pages{c.name, c.platform}
	body := c.page.Body()
	c.platform = c.page.Platform
	c.createDir()
	return body
}

func (c *cacher) save() *os.File {
	buf, err := ioutil.ReadAll(c.download())
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

	log.Println("Created:", c.file(), "bytes:", strconv.Itoa(ret))
	return c.search()
}

func (c *cacher) remove() {
	if c.name == "clearall.md" {
		if err := os.RemoveAll(cacheDir); err != nil {
			log.Println(err)
			os.Exit(1)
		}
		log.Println("Cache cleared")
		os.Exit(0)
	}

	if c.search() == nil {
		log.Println("Command:", c.cmd(), "not cached", c.file())
		os.Exit(1)
	}

	if err := os.Remove(c.file()); err != nil {
		log.Println(err)
		os.Exit(1)
	}

	log.Println("Removed:", c.cmd(), c.file())
	os.Exit(0)
}

func (c *cacher) createDir() {
	_, err := os.Stat(c.platformDir())
	if err != nil {
		if err := os.Mkdir(c.platformDir(), 0700); err != nil {
			log.Println(err)
			os.Exit(1)
		}
	}
}

func (c *cacher) readDir() []os.FileInfo {
	c.createDir()
	srcDir, err := ioutil.ReadDir(c.platformDir())
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	return srcDir
}
