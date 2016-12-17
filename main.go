package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
	"os"
)

const (
	usage                 = "OCI runtime for Graphene Library OS"
	logFile               = "/home/sridarshan/rungraphene_log"
	fifo_name             = "exec.fifo"
	rungrapheneWorkdir    = "/run/rungraphene"
	grapheneBootstrap     = "/home/sridarshan/524_1/Graphene"
	containerInfoJsonFile = "container_info.json"
	manifestTemplate      = "exec.manifest"
)

type ContainerCreateInfo struct {
	Bundle  string `json:"bundle"`
	PidFile string `json:"pidfile"`
	Console string `json:"console"`
	Id      string `json:"id"`
}

func main() {
	app := cli.NewApp()
	app.Name = "rungraphene"
	app.Usage = usage

	app.Commands = []cli.Command{
		createCommand,
		startCommand,
		killCommand,
		reexecCommand,
		deleteCommand,
	}
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug",
			Usage: "enable debug output for logging",
		},
		cli.StringFlag{
			Name:  "log",
			Value: "/dev/null",
			Usage: "set the log file path where internal debug information is written",
		},
		cli.StringFlag{
			Name:  "log-format",
			Value: "text",
			Usage: "set the format used by logs ('text' (default), or 'json')",
		},
		cli.StringFlag{
			Name:  "root",
			Value: "/run/runc",
			Usage: "root directory for storage of container state (this should be located in tmpfs)",
		},
		cli.StringFlag{
			Name:  "criu",
			Value: "criu",
			Usage: "path to the criu binary used for checkpoint and restore",
		},
		cli.BoolFlag{
			Name:  "systemd-cgroup",
			Usage: "enable systemd cgroup support, expects cgroupsPath to be of form \"slice:prefix:name\" for e.g. \"system.slice:runc:434234\"",
		},
	}

	// Setup the logging
	logfile, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening the log file")
	}
	defer logfile.Close()
	log.SetOutput(logfile)

	if err := app.Run(os.Args); err != nil {
		fmt.Println("error starting the app")
	}
}
