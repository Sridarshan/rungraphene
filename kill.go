package main

import (
	"fmt"

	"github.com/urfave/cli"
)

var killCommand = cli.Command{
	Name:        "kill",
	Usage:       "TODO",
	ArgsUsage:   "<container-id> <signal>",
	Description: "TODO",
	Action: func(context *cli.Context) error {
		fmt.Println("Nothing yet")
		return nil
	},
}
