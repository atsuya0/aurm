package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func getLocalVer(pkgName string) (string, error) {
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
