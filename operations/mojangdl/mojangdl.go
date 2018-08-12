package mojangdl

import (
	"encoding/json"
	"fmt"
	"github.com/pufferpanel/pufferd/environments"
	"github.com/pufferpanel/pufferd/programs/operations/ops"
	"net/http"
	"github.com/pufferpanel/apufferi/logging"
	"errors"
	"os"
	"path"
	"io"
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
			logging.Debugf("Version %s json located, downloading from %s", version.Id, version.Url)
			env.DisplayToConsole(fmt.Sprintf("Version %s json located, downloading from %s", version.Id, version.Url))
			//now, get the version json for this one...
			return downloadServerFromJson(version.Url, op.Target, env)
		}
	}

	env.DisplayToConsole("Could not locate version " + op.Version)

	return errors.New("Version not located: " + op.Version)
}

func downloadServerFromJson(url, target string, env environments.Environment) error {
	client := &http.Client{}
	response, err := client.Get(url)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	var data MojangVersionJson

	json.NewDecoder(response.Body).Decode(&data)

	serverBlock := data.Downloads["server"]

	logging.Debugf("Version jar located, downloading from %s", serverBlock.Url)
	env.DisplayToConsole(fmt.Sprintf("Version jar located, downloading from %s", serverBlock.Url))

	return downloadFile(serverBlock.Url, target, env)
}

func downloadFile(url, fileName string, env environments.Environment) error {
	target, err := os.Create(path.Join(env.GetRootDirectory(), fileName))
	if err != nil {
		return err
	}
	defer target.Close()

	client := &http.Client{}

	logging.Debug("Downloading: " + url)

	response, err := client.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	_, err = io.Copy(target, response.Body)
	return err
}

type MojangDlOperationFactory struct {
}

func (of MojangDlOperationFactory) Create(op ops.CreateOperation) ops.Operation {
	version := op.OperationArgs["version"].(string)
	target := op.OperationArgs["target"].(string)
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
