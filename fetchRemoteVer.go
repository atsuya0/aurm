package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"path"
	"strings"
)

const (
	ver = "pkgver="
	rel = "pkgrel="
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
		if strings.HasPrefix(text, ver) {
			ver = text[len(ver):]
		} else if strings.HasPrefix(text, rel) {
			rel = text[len(rel):]
		}
	}
	return fmt.Sprintf("%s-%s", ver, rel), nil
}
