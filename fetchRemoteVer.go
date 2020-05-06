package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"path"
	"strings"
)

func fetchRemoteVer(pkgName string) (string, error) {
	res, err := http.Get(aurHost + path.Join(aurPath, "plain/PKGBUILD?h=") + pkgName)
	if err != nil {
		return "", nil
	}
	defer func() {
		if err = res.Body.Close(); err != nil {
			log.Fatalf("%+v\n", err)
		}
	}()

	scanner := bufio.NewScanner(res.Body)
	var ver, rel string
	for scanner.Scan() {
		text := scanner.Text()
		if strings.HasPrefix(text, "pkgver=") {
			ver = text[len("pkgver="):]
		} else if strings.HasPrefix(text, "pkgrel=") {
			rel = text[len("pkgrel="):]
		}
	}
	return fmt.Sprintf("%s-%s", ver, rel), nil
}
