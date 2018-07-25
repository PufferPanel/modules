package main

import (
	"encoding/json"
	"fmt"
	"github.com/pufferpanel/pufferd/environments"
	"github.com/pufferpanel/pufferd/programs/operations/ops"
	"net/http"
)

const VERSION_JSON = "https://launchermeta.mojang.com/mc/game/version_manifest.json"

type MojangDl struct {
	Version string
	Target  string
}

func (op MojangDl) Run(env environments.Environment) error {
	client := &http.Client{}

	response, err := client.Get(VERSION_JSON)
	defer response.Body.Close()
	if err != nil {
		return err
	}

	var data MojangLauncherJson

	json.NewDecoder(response.Body).Decode(&data)

	for _, version := range data.Versions {
		if version.Id == op.Version {
			fmt.Println(version)
			//now, get the version json for this one...
			downloadServerFromJson(version.Url, op.Target, env)
		}
	}

	return nil
}

func downloadServerFromJson(url, file string, env environments.Environment) error {
	client := &http.Client{}
	response, err := client.Get(url)
	defer response.Body.Close()
	if err != nil {
		return err
	}

	var data MojangVersionJson

	json.NewDecoder(response.Body).Decode(&data)

	fmt.Println(data.Downloads["server"])

	return nil
}

type MojangDlOperationFactory struct {
}

func (of MojangDlOperationFactory) Create(op ops.CreateOperation) ops.Operation {
	version := op.DataMap["version"].(string)
	target := op.DataMap["target"].(string)
	return MojangDl{Version: version, Target: target}
}

func (of MojangDlOperationFactory) Key() string {
	return "mojangdl"
}

type MojangLauncherJson struct {
	Versions []MojangLauncherVersion `json:"versions"`
}

type MojangLauncherVersion struct {
	Id   string `json:"id"`
	Url  string `json:"url"`
	Type string `json:"type"`
}

type MojangVersionJson struct {
	Downloads map[string]MojangDownloadType `json:"downloads"`
}

type MojangDownloadType struct {
	Sha1 string `json:"sha1"`
	Size uint64 `json:"size"`
	Url  string `json:"url"`
}

var Factory MojangDlOperationFactory
