package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func getPkgNames() ([]string, error) {
	path, err := getDataPath()
	if err != nil {
		return make([]string, 0), fmt.Errorf("%w", err)
	}
	fp, err := os.Open(path)
	if err != nil {
		return make([]string, 0), fmt.Errorf("%w", err)
	}
	defer func() {
		if err = fp.Close(); err != nil {
			log.Fatalf("%+v\n", err)
		}
	}()
	scanner := bufio.NewScanner(fp)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, nil
}
