package main

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

type pkgDownloader struct {
	pkgNames []string
}

func getForeignPkgNames() ([]string, error) {
	cmd := exec.Command("pacman", "-Qmq")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return make([]string, 0), fmt.Errorf("%w", err)
	}

	if err := cmd.Start(); err != nil {
		return make([]string, 0), fmt.Errorf("%w", err)
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(stdout)
	if err := cmd.Wait(); err != nil {
		return make([]string, 0), fmt.Errorf("%w", err)
	}
	return strings.Split(strings.TrimSpace(buf.String()), "\n"), nil
}

func newPkgDownloader() (pkgDownloader, error) {
	pkgNames, err := getForeignPkgNames()
	if err != nil {
		return pkgDownloader{}, err
	}

	return pkgDownloader{pkgNames: pkgNames}, nil
}

func (p *pkgDownloader) fetchPkgIfNeeded() error {
	for _, pkgName := range p.pkgNames {
		localVer, err := getLocalVer(pkgName)
		exitError := &exec.ExitError{}
		if errors.As(err, &exitError) {
			if err := fetchPkg(pkgName); err != nil {
				return err
			}
			continue
		} else if err != nil {
			return err
		}
		remoteVer, err := fetchRemoteVer(pkgName)
		if err != nil {
			return err
		}
		if localVer != remoteVer {
			fmt.Printf("Download %s\n", pkgName)
			if err := fetchPkg(pkgName); err != nil {
				return err
			}
		}
	}
	return nil
}
