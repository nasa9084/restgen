package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/nasa9084/restgen/internal/pkg/generator"
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
	"internal/pkg/httperr",
	"cmd/server",
}

func (cmd NewCommand) Execute([]string) error {
	if err := cmd.createDirectories(); err != nil {
		return err
	}
	if err := cmd.createMakefile(); err != nil {
		return err
	}
	if err := cmd.createSpecFile(); err != nil {
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

func (cmd NewCommand) createMakefile() error {
	makefile := filepath.Join(cmd.Directory, "Makefile")
	if _, err := os.Stat(makefile); err == nil {
		log.Print("Makefile has been existing")
		return nil
	}
	out, err := generator.Template("/makefile.tmpl", makefile, cmd.Args.ApplicationName)
	if err != nil {
		return err
	}

	if opts.Verbose {
		log.Printf("generated Makefile\n%s", out)
	}
	return nil
}

func (cmd NewCommand) createSpecFile() error {
	specFile := filepath.Join(cmd.Directory, "api", "spec.yaml")
	if _, err := os.Stat(specFile); err == nil {
		log.Print("spec.yaml has been existing")
		return nil
	}
	out, err := generator.Template("/default_spec.yaml.tmpl", specFile, cmd.Args.ApplicationName)
	if err != nil {
		return err
	}
	if opts.Verbose {
		log.Printf("generated spec.yaml\n%s", out)
	}
	return nil
}
