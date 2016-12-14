package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	oci "github.com/opencontainers/runtime-spec/specs-go"
	"github.com/urfave/cli"
)

var createCommand = cli.Command{
	Name:        "create",
	Usage:       "TODO",
	ArgsUsage:   "TODO",
	Description: "TODO",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "bundle, b",
			Value: "",
			Usage: `path to the root of the bundle directory, defaults to the current directory`,
		},
		cli.StringFlag{
			Name:  "console",
			Value: "",
			Usage: "specify the pty slave path for use with the container",
		},
		cli.StringFlag{
			Name:  "pid-file",
			Value: "",
			Usage: "specify the file to write the process id to",
		},
		cli.BoolFlag{
			Name:  "no-pivot",
			Usage: "do not use pivot root to jail process inside rootfs.  This should be used whenever the rootfs is on top of a ramdisk",
		},
		cli.BoolFlag{
			Name:  "no-new-keyring",
			Usage: "do not create a new session keyring for the container.  This will cause the container to inherit the calling processes session key",
		},
	},
	Action: func(context *cli.Context) error {
		id := context.Args().First()
		bundle := context.String("bundle")
		pidfile := context.String("pid-file")
		console := context.String("console")

		containerDir := filepath.Join(rungrapheneWorkdir, id)
		log.Print("In create")

		_, err := os.Stat(containerDir)
		if err != nil {
			// this should not happen, cause delete/kill should have cleaned up properly

			// but for now, lets clean it up here
			os.RemoveAll(containerDir)
		}
		err = os.MkdirAll(containerDir, 0644)
		if err != nil {
			log.Println("Error creating work dir for container")
		}

		container := ContainerCreateInfo{
			Bundle:  bundle,
			Console: console,
			PidFile: pidfile,
			Id:      id,
		}

		containerJson, _ := json.Marshal(container)
		err = ioutil.WriteFile(filepath.Join(containerDir, containerInfoJsonFile), containerJson, 0644)
		if err != nil {
			log.Println("Error writing container create info to json")
		}

		/*
			cmd := exec.Command("ginit", id)
			if err := cmd.Start(); err != nil {
				if exErr, ok := err.(*exec.Error); ok {
					if exErr.Err == exec.ErrNotFound || exErr.Err == os.ErrNotExist {
						log.Println("ginit not installed on system")
					} else {
						log.Println("Error starting ginit", err)
					}
				}
				// can we return an error instead? TODO
				return nil
			}
		*/
		// copyGraphene(bundle)
		createPidFile(pidfile, 11)
		log.Println("Returning from create command")
		return nil
	},
}

func readSpec(bundle string) (oci.Spec, error) {
	var spec oci.Spec
	f, err := os.Open(filepath.Join(bundle, "config.json"))
	if err != nil {
		return spec, err
	}
	defer f.Close()
	if err := json.NewDecoder(f).Decode(&spec); err != nil {
		return spec, err
	}
	return spec, nil
}

func copyGraphene(bundle string) {
	log.Println("Starting bootstrap copy")
	cmd := exec.Command("cp", grapheneBootstrap, bundle, "-r")
	if err := cmd.Run(); err != nil {
		log.Println(err.Error())
		return
	}
	log.Println("Graphene bootstrap copied into bundle")
}

func createPidFile(path string, pid int) error {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_EXCL|os.O_SYNC, 0666)
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = fmt.Fprintf(f, "%d", pid)
	f.Close()
	return err
}
