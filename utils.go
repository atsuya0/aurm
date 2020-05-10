package main

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	aurHost = "https://aur.archlinux.org"
	aurPath = "/cgit/aur.git"
)

func getDataPath() (string, error) {
	dataPath := os.Getenv("XDG_DATA_HOME")
	if dataPath == "" {
		userHomeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("%w", err)
		}
		return filepath.Join(userHomeDir, ".local/share/aurm/packages.txt"), nil
	}
	return filepath.Join(dataPath, "aurm/packages.txt"), nil
}
