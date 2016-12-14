package main

import (
	// "bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	log "github.com/Sirupsen/logrus"
	oci "github.com/opencontainers/runtime-spec/specs-go"
	"github.com/urfave/cli"
)

const (
	loader = "LibOS/shim/test/native"
)

var startCommand = cli.Command{
	Name:        "start",
	Usage:       "TODO",
	ArgsUsage:   "[containerd-id]",
	Description: "TODO",
	Action: func(context *cli.Context) error {
		id := context.Args().First()
		container, err := getContainer(id)
		if err != nil {
			log.Println(err)
			return err
		}

		/*
			// open the fifo file and readA
			path := filepath.Join(rungrapheneWorkdir, id, fifo_name)
			fifo, err := os.Open(path)
			if err != nil {
				log.Println(err)
				return err
			}
			defer fifo.Close()
			data, err := ioutil.ReadAll(fifo)
			if err != nil {
				log.Println(err)
				return err
			}
			if len(data) == 0 {
				return fmt.Errorf("Container already started")
			}
		*/
		log.Println(container)

		spec, err := readSpec(container.Bundle)
		if err != nil {
			log.Println("Error reading the spec file", err)
			return err
		}
		manifestFile, err := generateManifestFile(container, spec)
		if err != nil {
			log.Println("Error writing manifest file", err)
			return err
		}

		cf, err := os.OpenFile(container.Console, os.O_RDWR, 0644)
		defer cf.Close()
		if err != nil {
			log.Println("Error in opening the console")
			return err
		}
		
		/*
		go io.Copy(cf, os.Stdout)
		go io.Copy(cf, os.Stderr)
		go io.Copy(os.Stdin, cf)
		*/
		

		cmd := exec.Command(filepath.Join(grapheneBootstrap, loader, "pal"), manifestFile)
		// cmd = exec.Command("echo", "dummy")
		log.Println(cmd.Args)
		cmd.Stdout = cf
		cmd.Stderr = cf
		cmd.Stdin = cf
		
		err = cmd.Run()
		log.Println("Run error: ", err)
		return err
		

		if err := cmd.Start(); err != nil {
			log.Println(err.Error())
			return err
		}
		go func() {
			// io.Copy(cf, os.Stdout)
			log.Println("Stdout closed")
		}()

		go func() {
			io.Copy(os.Stdin, cf)
			log.Println("Stdin closed")
		}()

		// io.Copy(cf, os.Stderr)
		log.Println("Stderr closed")

		io.WriteString(cf, "End\n")
		err = cf.Close()

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
