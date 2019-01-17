package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"

	flags "github.com/jessevdk/go-flags"
	"github.com/pkg/errors"
)

var (
	version  string
	revision string
)

type options struct {
	Verbose bool `short:"v" long:"verbose" description:"show verbose log"`

	// subcommands
	New      NewCommand      `command:"new"`
	Generate GenerateCommand `command:"generate" subcommands-optional:"yes"`
	Version  VersionCommand  `command:"version"`
}

var opts options

func main() {
	if err := exec(); err != nil {
		log.Printf("error: %v", err)
		os.Exit(1)
	}
	os.Exit(0)
}

func exec() error {
	parser := flags.NewParser(&opts, flags.HelpFlag|flags.PassDoubleDash)
	if _, err := parser.Parse(); err != nil {
		if fe, ok := err.(*flags.Error); ok && fe.Type == flags.ErrHelp {
			parser.WriteHelp(os.Stdout)
			return nil
		}
		return err
	}
	return nil
}

type VersionCommand struct{}

func (cmd VersionCommand) Execute([]string) error {
	if opts.Verbose {
		fmt.Printf("Version: %s\n", version)
		fmt.Printf("Revision: %s\n", revision)
		return nil
	}
	fmt.Println(version)
	return nil
}

type NewCommand struct {
	Directory string `short:"d" long:"directory" default:"." description:"target directory"`
}

var dirs = []string{
	"api",
	"internal/pkg/models",
}

func (cmd NewCommand) Execute([]string) error {
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
	specFile := filepath.Join(cmd.Directory, "api", "spec.yaml")
	if opts.Verbose {
		log.Printf("create OpenAPI Spec file: %s\n", specFile)
	}
	if _, err := os.Stat(specFile); err == nil {
		return errors.New("spec.yaml has been existing")
	}
	var buf bytes.Buffer
	buf.WriteString("---")
	buf.WriteString("\nopenapi: ")
	buf.WriteString("\ninfo: ")
	buf.WriteString("\n\ttitle: ")
	buf.WriteString("\n\tversion: ")
	buf.WriteString("\npaths: ")
	f, err := os.OpenFile(specFile, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.Write(buf.Bytes()); err != nil {
		return err
	}
	if opts.Verbose {
		log.Printf("generate spec.yaml\n%s", buf.String())
	}
	return nil
}
