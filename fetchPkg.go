package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
)

func fetchPkg(pkgName string) error {
	fileName := pkgName + ".tar.gz"
	res, err := http.Get(aurHost + path.Join(aurPath, "snapshot", fileName))
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	defer func() {
		if err = res.Body.Close(); err != nil {
			log.Fatalf("%+v\n", err)
		}
	}()

	gzipReader, err := gzip.NewReader(res.Body)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	defer func() {
		if err = gzipReader.Close(); err != nil {
			log.Fatalf("%+v\n", err)
		}
	}()

	tarReader := tar.NewReader(gzipReader)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("%w", err)
		}
		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.Mkdir(header.Name, 0700); err != nil {
				return fmt.Errorf("%w", err)
			}
		case tar.TypeReg:
			newFile, err := os.Create(header.Name)
			if err != nil {
				return fmt.Errorf("%w", err)
			}
			defer func() {
				if err = newFile.Close(); err != nil {
					log.Fatalf("%+v\n", err)
				}
			}()
			if _, err := io.Copy(newFile, tarReader); err != nil {
				return fmt.Errorf("%w", err)
			}
		}
	}
	return nil
}
