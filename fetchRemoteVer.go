package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
)

type pkg struct {
	Results []struct {
		Version string `json:"Version"`
	} `json:"results"`
}

func fetchRemoteVer(pkgName string) (string, error) {
	res, err := http.Get(aurHost + "/rpc/?v=5&type=info&arg[]=" + pkgName)
	if err != nil {
		return "", err
	}
	defer func() {
		if err = res.Body.Close(); err != nil {
			log.Fatalf("%+v\n", err)
		}
	}()

	body, err := ioutil.ReadAll(res.Body)
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
