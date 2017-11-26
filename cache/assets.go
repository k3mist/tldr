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

	contents, err := ioutil.ReadAll(resp)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Create(zipFile)
	if err != nil {
		log.Fatal(err)
	}

	_, err = file.Write(contents)
	file.Close()
	if err != nil {
		log.Fatal(err)
	}

	var fileNames []string

	r, err := zip.OpenReader(zipFile)
	defer r.Close()
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range r.File {
		df, err := f.Open()
		if err != nil {
			log.Fatal(err)
		}

		filePath := filepath.Join(cacheDir, f.Name)
		fileNames = append(fileNames, filePath)

		if f.FileInfo().IsDir() {
			os.MkdirAll(filePath, os.ModePerm)
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
			f.Close()
		}
		df.Close()
	}
}
