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

// Pages provides the retrieval of the TLDR assets and repository pages.
type Pages struct {
	Name     string
	Platform platform.Platform
	cfg      config.Options
}

// New creates a new Pages instance
func New(name string, plat platform.Platform) *Pages {
	return &Pages{
		Name:     name,
		Platform: plat,
		cfg:      config.Config,
	}
}

func (p *Pages) url() string {
	var uri string
	if p.cfg.PagesURI != "" {
		uri = p.cfg.PagesURI
	} else {
		uri = raw
	}
	return uri + p.Platform.String() + "/" + p.Name
}

// Zip returns the result of the downloaded TLDR pages zip file.
func (p *Pages) Zip() io.ReadCloser {
	var uri string
	if p.cfg.ZipURI != "" {
		uri = p.cfg.ZipURI
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

// Body returns the result of a http lookup from the main TLDR repository for a
// TLDR page that was not found in the local cache.
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
