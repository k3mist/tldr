package pages

import (
	"io"
	"log"
	"net/http"

	"bitbucket.org/djr2/tldr/config"
	"bitbucket.org/djr2/tldr/platform"
)

const (
	zip = "https://tldr.sh/assets/tldr.zip"
	raw = "https://raw.github.com/tldr-pages/tldr/master/pages/"
)

type Pages struct {
	Name     string
	Platform platform.Platform
}

func (p *Pages) url() string {
	var uri string
	if config.Config.PagesURI != "" {
		uri = config.Config.PagesURI
	} else {
		uri = raw
	}
	return uri + p.Platform.String() + "/" + p.Name
}

func (p *Pages) Zip() io.ReadCloser {
	var uri string
	if config.Config.ZipURI != "" {
		uri = config.Config.ZipURI
	} else {
		uri = zip
	}
	zpr, err := http.Get(uri)
	if err != nil {
		log.Fatal(err)
	}
	if zpr.StatusCode != http.StatusOK {
		log.Fatal("Problem getting: ", zip, " Server Error: ", zpr.StatusCode)
	}
	return zpr.Body
}

func (p *Pages) Body() io.ReadCloser {
	log.Println("Retrieving:", p.url())
	cnr, err := http.Get(p.url())
	if err != nil {
		log.Fatal(err)
	}
	if cnr.StatusCode != http.StatusOK {
		log.Println("Problem getting:", p.Name, "Server Error:", cnr.StatusCode)
		return nil
	}
	return cnr.Body
}
