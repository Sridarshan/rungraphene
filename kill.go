package main

import (
	"os"
	"fmt"
	"path/filepath"
	
	"github.com/urfave/cli"
	log "github.com/Sirupsen/logrus"
)

var killCommand = cli.Command{
	Name:        "kill",
	Usage:       "TODO",
	ArgsUsage:   "<container-id> <signal>",
	Description: "TODO",
	Action: func(context *cli.Context) error {
		log.Println("Nothing yet")
		return nil
	},
}

var deleteCommand = cli.Command {
	Name:        "delete",
	Usage:       "TODO",
	ArgsUsage:   "",
	Description: "TODO",
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