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
	"log"
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
		localVer, err := p.getLocalVer(pkgName)
		exitError := &exec.ExitError{}
		if errors.As(err, &exitError) {
			if err := p.fetchPkg(pkgName); err != nil {
				return err
			}
			continue
		} else if err != nil {
			return err
		}
		remoteVer, err := p.fetchRemoteVer(pkgName)
		if err != nil {
			return err
		}
		if localVer != remoteVer {
			fmt.Printf("Download %s\n", pkgName)
			if err := p.fetchPkg(pkgName); err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *pkgDownloader) getLocalVer(pkgName string) (string, error) {
	cmd := exec.Command("pacman", "-Qi", pkgName)
	os.Setenv("LANG", "C")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	if err := cmd.Start(); err != nil {
		return "", fmt.Errorf("%w", err)
	}
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		text := scanner.Text()
		if strings.HasPrefix(text, "Version") {
			ver := strings.TrimSpace(text[strings.Index(text, ":")+1:])
			if err := cmd.Wait(); err != nil {
				return "", fmt.Errorf("%w", err)
			}
			return ver, nil
		}
	}
	if err := cmd.Wait(); err != nil {
		return "", fmt.Errorf("%w", err)
	}
	return "", errors.New("Cannot get package local version.")
}

type pkg struct {
	Results []struct {
		Version string `json:"Version"`
	} `json:"results"`
}

func (p *pkgDownloader) fetchRemoteVer(pkgName string) (string, error) {
	res, err := http.Get(aurHost + "/rpc/?v=5&type=info&arg[]=" + pkgName)
	if err != nil {
		return "", err
	}
	defer func() {
		if err = res.Body.Close(); err != nil {
			log.Fatalf("%+v\n", err)
		}
	}()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	pkg := pkg{}
	err = json.Unmarshal(body, &pkg)
	if err != nil {
		return "", err
	}
	if len(pkg.Results) < 1 {
		return "", errors.New("Cannot fetch the remote package version")
	}
	return pkg.Results[0].Version, nil
}

func (p *pkgDownloader) fetchPkg(pkgName string) error {
	fileName := pkgName + ".tar.gz"
	res, err := http.Get(aurHost + path.Join("/cgit/aur.git/snapshot", fileName))
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	defer func() {
		if err = res.Body.Close(); err != nil {
			log.Fatalf("%+v\n", err)
		}
	}()
	if res.StatusCode == http.StatusNotFound {
		return errors.New(res.Status)
	}

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
		if errors.Is(err, io.EOF) {
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
