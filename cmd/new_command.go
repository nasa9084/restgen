package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

type NewCommand struct {
	Directory string `short:"d" long:"directory" default:"." description:"target directory"`
	Args      struct {
		ApplicationName string `positional-arg-name:"APPLICATION_NAME"`
	} `positional-args:"yes"`
}

var dirs = []string{
	"api",
	"internal/pkg/models",
	"cmd",
}

func (cmd NewCommand) Execute([]string) error {
	if err := cmd.createDirectories(); err != nil {
		return err
	}
	if err := cmd.createSpecFile(); err != nil {
		return err
	}
	if err := cmd.createServerMain(); err != nil {
		return err
	}
	return nil
}

func (cmd NewCommand) createDirectories() error {
	if opts.Verbose {
		log.Printf("create directory: %s", cmd.Directory)
	}
	if err := os.MkdirAll(cmd.Directory, 0755); err != nil {
		return err
	}
	for _, dir := range dirs {
		dirPath := filepath.Join(cmd.Directory, dir)
		if opts.Verbose {
			log.Printf("create directory: %s", dirPath)
		}
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return err
		}
	}
	return nil
}

func (cmd NewCommand) createSpecFile() error {
	specFile := filepath.Join(cmd.Directory, "api", "spec.yaml")
	if opts.Verbose {
		log.Printf("create OpenAPI Spec file: %s\n", specFile)
	}
	if _, err := os.Stat(specFile); err == nil {
		return errors.New("spec.yaml has been existing")
	}
	defaultSpec, err := Assets.Open("/assets/default_spec.yaml.tmpl")
	if err != nil {
		return errors.Wrap(err, "opening from assets")
	}
	defer defaultSpec.Close()

	b, err := ioutil.ReadAll(defaultSpec)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(specFile, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := fmt.Fprintf(f, string(b), cmd.Args.ApplicationName); err != nil {
		return err
	}
	if opts.Verbose {
		log.Printf("generated spec.yaml\n%s", string(b))
	}
	return nil
}

func (cmd NewCommand) createServerMain() error {
	srvGo, err := Assets.Open("/assets/server_main.go.tmpl")
	if err != nil {
		return err
	}
	defer srvGo.Close()
	srvGoFile := filepath.Join(cmd.Directory, "cmd", "server", "main.go")
	f, err := os.OpenFile(srvGoFile, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	if opts.Verbose {
		log.Printf("generated %s", srvGoFile)
	}
	return nil
}
