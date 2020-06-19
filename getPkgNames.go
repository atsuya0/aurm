package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

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
