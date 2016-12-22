package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	log "github.com/Sirupsen/logrus"
	oci "github.com/opencontainers/runtime-spec/specs-go"
	"github.com/urfave/cli"
)

var startCommand = cli.Command{
	Name:        "start",
	Usage:       "call to start an already created container",
	ArgsUsage:   "[containerd-id]",
	Description: "call to start an already created container",
	Action: func(context *cli.Context) error {
		id := context.Args().First()
		path := filepath.Join(rungrapheneWorkdir, id, fifo_name)
		fifo, err := os.Open(path)
		if err != nil {
			log.Println("error opening fifo: ", err.Error())
			return err
		}
		defer fifo.Close()
		data, err := ioutil.ReadAll(fifo)
		if err != nil {
			log.Println("error reading fifo: ", err.Error())
			return err
		}
		if len(data) == 0 {
			return fmt.Errorf("Container already started")
		}

		return nil
	},
}

func generateManifestFile(container ContainerCreateInfo, spec oci.Spec) (string, error) {
	manifestInputPath := filepath.Join(grapheneBootstrap, manifestTemplate)
	content, err := ioutil.ReadFile(manifestInputPath)
	if err != nil {
		return "", err
	}
	lines := strings.Split(string(content), "\n")
	containerRoot := spec.Root.Path
	if filepath.IsAbs(containerRoot) == false {
		containerRoot = filepath.Join(container.Bundle, spec.Root.Path)
	}

	execPath := filepath.Join(containerRoot, spec.Process.Args[0])
	for i, line := range lines {
		if strings.Contains(line, "DOCKER_EXEC_PATH") {
			lines[i] = strings.Replace(line, "DOCKER_EXEC_PATH", "file:"+execPath, 1)
		}
		if strings.Contains(line, "DOCKER_ROOT") {
			lines[i] = strings.Replace(line, "DOCKER_ROOT", "file:"+containerRoot, 1)
		}
	}
	outfilePath := filepath.Join(rungrapheneWorkdir, container.Id, "exec.manifest")
	err = ioutil.WriteFile(outfilePath, []byte(strings.Join(lines, "\n")), 0644)
	if err != nil {
		return "", err
	}
	return outfilePath, nil
}

func getContainer(id string) (ContainerCreateInfo, error) {
	var containerInfo ContainerCreateInfo

	containerDir := filepath.Join(rungrapheneWorkdir, id)
	if _, err := os.Stat(containerDir); err != nil {
		return containerInfo, err
	}
	containerJsonFile, err := os.Open(filepath.Join(containerDir, containerInfoJsonFile))
	if err != nil {
		return containerInfo, err
	}
	defer containerJsonFile.Close()

	if err := json.NewDecoder(containerJsonFile).Decode(&containerInfo); err != nil {
		return containerInfo, err
	}
	return containerInfo, nil
}
