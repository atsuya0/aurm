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
	pkgs, err := getForeignPkgNames()
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	for _, pkg := range pkgs {
		localVer, err := getLocalVer(pkg)
		exitError := &exec.ExitError{}
		if errors.As(err, &exitError) {
			if err := fetchPkg(pkg); err != nil {
				return fmt.Errorf("%w", err)
			}
			continue
		} else if err != nil {
			return fmt.Errorf("%w", err)
		}
		remoteVer, err := fetchRemoteVer(pkg)
		if err != nil {
			return fmt.Errorf("%w", err)
		}
		if localVer != remoteVer {
			if err := fetchPkg(pkg); err != nil {
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
