package cache

import (
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"bitbucket.org/djr2/tldr/config"
	"bitbucket.org/djr2/tldr/pages"
	"bitbucket.org/djr2/tldr/platform"
)

type cacher struct {
	plat platform.Platform
	lang string
	name string
}

func (c *cacher) platformDir() string {
	return cacheDir + "/pages." + c.lang + "/" + c.plat.String()
}

func (c *cacher) file() string {
	return c.platformDir() + "/" + c.name
}

func (c *cacher) cmd() string {
	return strings.TrimSuffix(c.name, `.md`)
}

func (c *cacher) search() *os.File {
	var tried []platform.Platform
	c.plat = getPlatform(c.plat)
	tried = append(tried, c.plat)

	cached := c.find()
	if cached == nil {
		c.plat = platform.COMMON
		tried = append(tried, c.plat)
	}

	cached = c.find()
	if cached == nil {
		cached = c.extendedSearch(tried)
	}

	if cached != nil {
		info, err := cached.Stat()
		if err != nil {
			log.Fatal(err)
		}

		var hours int
		if expires := config.Config.CacheExpiration * 24; expires < 1 {
			hours = 1
		} else {
			hours = expires
		}

		if info.ModTime().Add(time.Hour * time.Duration(hours)).Before(time.Now()) {
			log.Println("Cache older than 30 days")
			return c.save()
		}
	}

	return cached
}

func (c *cacher) extendedSearch(tried []platform.Platform) *os.File {
	for _, plat := range platform.Platforms() {
		c.plat = getPlatform(platform.Parse(plat))
		if file := c.find(); file != nil {
			return file
		}
	}

	for _, plat := range tried {
		c.plat = plat
		if file := c.save(); file != nil {
			return file
		}
	}

	return nil
}

func (c *cacher) find() *os.File {
	for _, fileInfo := range c.readDir() {
		if fileInfo.Name() == c.name {
			file, err := os.Open(c.file())
			if err != nil {
				log.Fatal(err)
			}
			return file
		}
	}
	return nil
}

func (c *cacher) download() io.ReadCloser {
	page := pages.New(c.name, c.lang, c.plat)
	c.plat = page.Platform
	c.createDir()
	return page.Body()
}

func (c *cacher) save() *os.File {
	down := c.download()
	if down == nil {
		return nil
	}

	buf, err := io.ReadAll(down)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Create(c.file())
	if err != nil {
		log.Fatal(err)
	}

	ret, err := file.Write(buf)
	defer file.Close() // nolint: errcheck
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Created:", c.file(), "bytes:", strconv.Itoa(ret))
	return c.search()
}

func (c *cacher) remove() {
	if c.name == "clearall.md" {
		if err := os.RemoveAll(cacheDir); err != nil {
			log.Fatal(err)
		}
		log.Println("Cache cleared")
		os.Exit(0)
	}

	if c.search() == nil {
		log.Fatal("Command: ", c.cmd(), " not cached ", c.file())
	}

	if err := os.Remove(c.file()); err != nil {
		log.Fatal(err)
	}

	log.Println("Removed:", c.cmd(), c.file())
	os.Exit(0)
}

func (c *cacher) createDir() {
	_, err := os.Stat(c.platformDir())
	if err != nil {
		if err := os.Mkdir(c.platformDir(), 0700); err != nil {
			log.Fatal(err)
		}
	}
}

func (c *cacher) readDir() []os.DirEntry {
	c.createDir()
	srcDir, err := os.ReadDir(c.platformDir())
	if err != nil {
		log.Fatal(err)
	}
	return srcDir
}
