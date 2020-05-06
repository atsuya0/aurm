package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func list() error {
	pkgs, err := getPkgNames()
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	fmt.Println(strings.Join(pkgs, "\n"))
	return nil
}

func fetchPkgIfNeeded() error {
	pkgs, err := getPkgNames()
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
	var l = flag.Bool("l", false, "list")
	flag.Parse()
	if *l {
		if err := list(); err != nil {
			log.Fatalf("%+v\n", err)
		}
		return
	}
	if err := fetchPkgIfNeeded(); err != nil {
		log.Fatalf("%+v\n", err)
	}
}
