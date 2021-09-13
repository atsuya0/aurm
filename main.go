package main

import (
	"log"
)

const (
	aurHost = "https://aur.archlinux.org"
)

func init() {
	log.SetFlags(log.Llongfile)
}

func main() {
	pkgDownloader, err := newPkgDownloader()
	if err != nil {
		log.Fatalf("%+v\n", err)
	}

	if pkgDownloader.fetchPkgIfNeeded() != nil {
		log.Fatalf("%+v\n", err)
	}
}
