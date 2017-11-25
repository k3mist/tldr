package cache

import (
	"archive/zip"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"bitbucket.org/djr2/tldr/pages"
)

func getAssets() {
	zipFile := cacheDir + "/assets.zip"
	if info, err := os.Stat(zipFile); err == nil {
		if info.ModTime().Add(time.Hour * 720).After(time.Now()) {
			return
		}
	}

	page := pages.Pages{}
	resp := page.Zip()

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Create(zipFile)
	if err != nil {
		log.Fatal(err)
	}

	_, err = file.Write(contents)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	var filenames []string

	r, err := zip.OpenReader(zipFile)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			log.Fatal(err)
		}

		fpath := filepath.Join(cacheDir, f.Name)
		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
		} else {
			var fdir string
			if lastIndex := strings.LastIndex(fpath, string(os.PathSeparator)); lastIndex > -1 {
				fdir = fpath[:lastIndex]
			}

			err = os.MkdirAll(fdir, os.ModePerm)
			if err != nil {
				log.Fatal(err)
			}
			f, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()

			_, err = io.Copy(f, rc)
			if err != nil {
				log.Fatal(err)
			}

		}
	}
}
