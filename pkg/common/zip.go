package common

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

func DownloadZipFromURL(url, destFile string) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}

	zipBuffer, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return err
	}
	LogOK(fmt.Sprintf("Read %v bytes from %v", len(zipBuffer), url))

	//TODO (doug4j@gmail.com): Ensure the this is the correct write mode "os.ModePerm"
	err = ioutil.WriteFile(destFile, zipBuffer, os.ModePerm)
	if err != nil {
		err = fmt.Errorf("Could not read values file: " + err.Error())
		LogError(err.Error())
		return err
	}
	LogOK(fmt.Sprintf("Wrote temp %v bytes of runtime bundle template to %v", len(zipBuffer), destFile))
	return nil
}

func UnzipFromFileToDir(src, dest string) (outputDir string, err error) {
	//https://stackoverflow.com/questions/20357223/easy-way-to-unzip-file-with-golang
	r, err := zip.OpenReader(src)
	if err != nil {
		return "", err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	os.MkdirAll(dest, 0755)

	var topLevelDirProcessed bool

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(dest, f.Name)
		if !topLevelDirProcessed {
			outputDir = path
			topLevelDirProcessed = true
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), f.Mode())
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return "", err
		}
	}

	return outputDir, nil
}
