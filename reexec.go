package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"path/filepath"
	
	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
)

var reexecCommand = cli.Command{
	Name:        "reexec",
	Usage:       "Not to be called externally",
	ArgsUsage:   "TODO",
	Description: "TODO",
	Action: func(context *cli.Context) error {
		id := os.Args[2]
		fifo_path := filepath.Join(rungrapheneWorkdir, id, fifo_name)
		if err := syscall.Mkfifo(fifo_path, 0644); err != nil {
			log.Println("error creating fifo, ", err.Error())
			log.Println("fifo_path: ", fifo_path)
			return err
		}
		fifo, err := os.OpenFile(fifo_path, os.O_WRONLY, 0)
		if err != nil {
			log.Println("error opening fifo, ", err.Error())
			return err
		}
		
		if _, err := syscall.Write(int(fifo.Fd()), []byte("0")); err != nil {
			log.Println("write 0 to fifo failed")
			return err
		}
		
		log.Println("Other side of fifo read, now conitnuing..")
		
		container, err := getContainer(id)
		if err != nil {
			log.Println("error getting container: ", err)
			return err
		}
		log.Println("Container: ", container)
		
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
		
		err = dupStdio(container.Console)
		if err != nil {
			log.Println("error duping, ", err)
			return err
		}
		
		palpath := filepath.Join(grapheneBootstrap, loader, "pal")
		/*
		This would have been better, instead of exec.Command way
		if err := syscall.Exec(palpath, []string{manifestFile}, os.Environ()); err != nil {
			fmt.Println("error execing pal, ", err.Error())
		}
		*/
		
		cmd := exec.Command(palpath, manifestFile)
		cmd.Stdout = os.Stdout
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr
		
		if err := cmd.Run(); err != nil {
			fmt.Println("error running pal")
		}
		return nil
	},
}


// open is a clone of os.OpenFile without the O_CLOEXEC used to open the pty slave.
func open(slavePath string, flag int) (*os.File, error) {
	r, e := syscall.Open(slavePath, flag, 0)
	if e != nil {
		return nil, &os.PathError{
			Op:   "open",
			Path: slavePath,
			Err:  e,
		}
	}
	return os.NewFile(uintptr(r), slavePath), nil
}

// dupStdio opens the slavePath for the console and dups the fds to the current
// processes stdio, fd 0,1,2.
func dupStdio(slavePath string) error {
	slave, err := open(slavePath, syscall.O_RDWR)
	if err != nil {
		return err
	}
	fd := int(slave.Fd())
	for _, i := range []int{0, 1, 2} {
		if err := syscall.Dup3(fd, i, 0); err != nil {
			return err
		}
	}
	return nil
}
