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
	Arg       struct {
		ApplicationName string `positional-arg-name:"APPLICATION_NAME"`
	} `positional-args:"yes"`
}

var dirs = []string{
	"api",
	"internal/pkg/models",
	"cmd",
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
	buf.WriteString("\nopenapi: 3.0.2")
	buf.WriteString("\ninfo: ")
	fmt.Fprintf(&buf, "\n  title: %s", cmd.Arg.ApplicationName)
	buf.WriteString("\n  version: v0.0.1")
	buf.WriteString("\npaths:")
	buf.WriteString("\n  /:")
	buf.WriteString("\n    get:")
	buf.WriteString("\n      description: health check")
	buf.WriteString("\n      operationId: HealthCheck")
	buf.WriteString("\n      responses:")
	buf.WriteString("\n        '200':")
	buf.WriteString("\n          $ref: \"#/components/responses/HealthCheckResponse\"")
	buf.WriteString("\ncomponents:")
	buf.WriteString("\n  responses:")
	buf.WriteString("\n    HealthCheckResponse:")
	buf.WriteString("\n      description: response for HealthCheck")
	buf.WriteString("\n      content:")
	buf.WriteString("\n        application/json:")
	buf.WriteString("\n          schema:")
	buf.WriteString("\n            type: object")
	buf.WriteString("\n            properties:")
	buf.WriteString("\n              status:")
	buf.WriteString("\n                type: string")
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
