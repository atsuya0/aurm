package main

import (
	"archive/tar"
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
)

type pkgDownloader struct {
	pkgNames []string
}

func getForeignPkgNames() ([]string, error) {
	errMsg := "Failed to get the foreign package names: %w"

	cmd := exec.Command("pacman", "-Qmq")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return make([]string, 0), fmt.Errorf(errMsg, err)
	}

	if err := cmd.Start(); err != nil {
		return make([]string, 0), fmt.Errorf(errMsg, err)
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(stdout)
	if err := cmd.Wait(); err != nil {
		return make([]string, 0), fmt.Errorf(errMsg, err)
	}
	return strings.Split(strings.TrimSpace(buf.String()), "\n"), nil
}

func newPkgDownloader() (pkgDownloader, error) {
	pkgNames, err := getForeignPkgNames()
	if err != nil {
		return pkgDownloader{}, fmt.Errorf("Failed to build the package downloader: %w", err)
	}

	return pkgDownloader{pkgNames: pkgNames}, nil
}

func (p *pkgDownloader) fetchPkgIfNeeded() error {
	errMsg := "Failed to fetch the package: %w"

	for _, pkgName := range p.pkgNames {
		localVer, err := p.getLocalVer(pkgName)
		exitError := &exec.ExitError{}
		if errors.As(err, &exitError) {
			if err := p.fetchPkg(pkgName); err != nil {
				return fmt.Errorf(errMsg, err)
			}
			continue
		} else if err != nil {
			return fmt.Errorf(errMsg, err)
		}
		remoteVer, err := p.fetchRemoteVer(pkgName)
		if err != nil {
			return fmt.Errorf(errMsg, err)
		}
		if localVer != remoteVer {
			fmt.Printf("Download %s\n", pkgName)
			if err := p.fetchPkg(pkgName); err != nil {
				return fmt.Errorf(errMsg, err)
			}
		}
	}
	return nil
}

func (p *pkgDownloader) getLocalVer(pkgName string) (string, error) {
	errMsg := "Failed to get the package local version"

	cmd := exec.Command("pacman", "-Qi", pkgName)
	os.Setenv("LANG", "C")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", fmt.Errorf(errMsg+": %w", err)
	}

	if err := cmd.Start(); err != nil {
		return "", fmt.Errorf(errMsg+": %w", err)
	}
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		text := scanner.Text()
		if strings.HasPrefix(text, "Version") {
			ver := strings.TrimSpace(text[strings.Index(text, ":")+1:])
			if err := cmd.Wait(); err != nil {
				return "", fmt.Errorf(errMsg+": %w", err)
			}
			return ver, nil
		}
	}
	if err := cmd.Wait(); err != nil {
		return "", fmt.Errorf(errMsg+": %w", err)
	}
	return "", errors.New(errMsg)
}

type pkg struct {
	Results []struct {
		Version string `json:"Version"`
	} `json:"results"`
}

func (p *pkgDownloader) fetchRemoteVer(pkgName string) (ver string, err error) {
	errMsg := "Failed to fetch the package remote version"

	res, err := http.Get(aurHost + "/rpc/?v=5&type=info&arg[]=" + pkgName)
	if err != nil {
		return "", fmt.Errorf(errMsg+": %w", err)
	}
	defer func() {
		if deferErr := res.Body.Close(); deferErr != nil {
			err = fmt.Errorf(errMsg+": %w", deferErr)
		}
	}()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf(errMsg+": %w", err)
	}
	pkg := pkg{}
	err = json.Unmarshal(body, &pkg)
	if err != nil {
		return "", fmt.Errorf(errMsg+": %w", err)
	}
	if len(pkg.Results) < 1 {
		return "", errors.New(errMsg)
	}
	return pkg.Results[0].Version, nil
}

func (p *pkgDownloader) fetchPkg(pkgName string) (err error) {
	errMsg := "Failed to fetch the package: %w"

	fileName := pkgName + ".tar.gz"
	res, err := http.Get(aurHost + path.Join("/cgit/aur.git/snapshot", fileName))
	if err != nil {
		return fmt.Errorf(errMsg, err)
	}
	defer func() {
		if deferErr := res.Body.Close(); deferErr != nil {
			err = fmt.Errorf(errMsg, err)
		}
	}()
	if res.StatusCode == http.StatusNotFound {
		return fmt.Errorf("%w", errors.New(res.Status))
	}

	gzipReader, err := gzip.NewReader(res.Body)
	if err != nil {
		return fmt.Errorf(errMsg, err)
	}
	defer func() {
		if deferErr := gzipReader.Close(); deferErr != nil {
			err = fmt.Errorf(errMsg, err)
		}
	}()

	tarReader := tar.NewReader(gzipReader)
	for {
		header, err := tarReader.Next()
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return fmt.Errorf(errMsg, err)
		}
		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.Mkdir(header.Name, 0700); err != nil {
				return fmt.Errorf(errMsg, err)
			}
		case tar.TypeReg:
			newFile, err := os.Create(header.Name)
			if err != nil {
				return fmt.Errorf(errMsg, err)
			}
			defer func() {
				if deferErr := newFile.Close(); deferErr != nil {
					err = fmt.Errorf(errMsg, err)
				}
			}()
			if _, err := io.Copy(newFile, tarReader); err != nil {
				return fmt.Errorf(errMsg, err)
			}
		}
	}
	return nil
}
