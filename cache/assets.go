package cache

import (
	"archive/zip"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"bitbucket.org/djr2/tldr/pages"
)

func getAssets() { // nolint: gocyclo
	zipFile := cacheDir + "/assets.zip"
	if info, err := os.Stat(zipFile); err == nil {
		if info.ModTime().Add(time.Hour * 720).After(time.Now()) {
			return
		}
	}

	page := pages.Pages{}
	resp := page.Zip()

	contents, err := io.ReadAll(resp)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Create(zipFile)
	if err != nil {
		log.Fatal(err)
	}

	_, err = file.Write(contents)
	defer file.Close() // nolint: errcheck
	if err != nil {
		log.Fatal(err)
	}

	r, err := zip.OpenReader(zipFile)
	defer r.Close() // nolint: errcheck, megacheck
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range r.File {
		df, oerr := f.Open()
		if oerr != nil {
			log.Fatal(oerr)
		}

		filePath := filepath.Join(cacheDir, f.Name)
		if f.FileInfo().IsDir() {
			derr := os.MkdirAll(filePath, os.ModePerm)
			if derr != nil {
				log.Fatal(derr)
			}
		} else {
			var fileDir string
			if lastIndex := strings.LastIndex(filePath, string(os.PathSeparator)); lastIndex > -1 {
				fileDir = filePath[:lastIndex]
			}

			err = os.MkdirAll(fileDir, os.ModePerm)
			if err != nil {
				log.Fatal(err)
			}

			f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				log.Fatal(err)
			}

			_, err = io.Copy(f, df)
			if err != nil {
				log.Fatal(err)
			}
			if err := f.Close(); err != nil {
				log.Println(err)
			}
		}
		if err := df.Close(); err != nil {
			log.Println(err)
		}
	}
}
