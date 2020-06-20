package main

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
)

const (
	aurHost = "https://aur.archlinux.org"
)

func fetchPkgIfNeeded() error {
	pkgNames, err := getForeignPkgNames()
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	for _, pkgName := range pkgNames {
		localVer, err := getLocalVer(pkgName)
		exitError := &exec.ExitError{}
		if errors.As(err, &exitError) {
			if err := fetchPkg(pkgName); err != nil {
				return fmt.Errorf("%w", err)
			}
			continue
		} else if err != nil {
			return fmt.Errorf("%w", err)
		}
		remoteVer, err := fetchRemoteVer(pkgName)
		if err != nil {
			return fmt.Errorf("%w", err)
		}
		if localVer != remoteVer {
			if err := fetchPkg(pkgName); err != nil {
				return fmt.Errorf("%w", err)
			}
		}
	}
	return nil
}

func main() {
	if err := fetchPkgIfNeeded(); err != nil {
		log.Fatalf("%+v\n", err)
	}
}
