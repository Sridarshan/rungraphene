package main

import (
	"fmt"
	"os"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
)

var killCommand = cli.Command{
	Name:        "kill",
	Usage:       "kills a container",
	ArgsUsage:   "<container-id> <signal>",
	Description: "kills a container",
	Action: func(context *cli.Context) error {
		log.Println("Nothing yet")
		return nil
	},
}

var deleteCommand = cli.Command{
	Name:        "delete",
	Usage:       "deletes a container",
	ArgsUsage:   "",
	Description: "deletes a container",
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("Not enough args to delete command")
		}
		id := context.Args().First()
		containerDir := filepath.Join(rungrapheneWorkdir, id)
		if _, err := os.Stat(containerDir); os.IsNotExist(err) {
			log.Println("delete error, container does not exists")
			return nil
		}
		if err := os.RemoveAll(containerDir); err != nil {
			log.Println("Error removing container dir")
		}
		return nil
	},
}
