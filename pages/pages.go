package pages

import (
	"io"
	"log"
	"net/http"
	"os"

	"bitbucket.org/djr2/tldr/platform"
)

const (
	raw = "https://raw.github.com/tldr-pages/tldr/master/pages/"
)

type Pages struct {
	Name     string
	Platform string
}

func (p *Pages) url() string {
	return raw + p.Platform + "/" + p.Name
}

func (p *Pages) Body() io.ReadCloser {
	log.Println("Retrieving:", p.url())
	cnr, err := http.Get(p.url())
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	if cnr.StatusCode != http.StatusOK {
		log.Println("Problem getting:", p.Name, "Server Error:", cnr.StatusCode)

		p.Platform = platform.Actual().String()
		log.Println("Trying by platform:", p.Platform)

		log.Println("Retrieving:", p.url())
		pmr, err := http.Get(p.url())
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		if pmr.StatusCode != http.StatusOK {
			log.Println("Problem getting:", p.Platform, "Server Error:", pmr.StatusCode)
			os.Exit(1)
		}
		return pmr.Body
	}
	return cnr.Body
}
