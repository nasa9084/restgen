package main

import (
	"bytes"
	"fmt"
	"go/format"
	"log"
	"os"
	"path/filepath"

	flags "github.com/jessevdk/go-flags"
	openapi "github.com/nasa9084/go-openapi"
	"github.com/nasa9084/restgen/internal/generator"
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

type GenerateCommand struct {
	Directory string `short:"d" long:"directory" default:"." description:"target directory"`

	Schema   GenerateSchemaCommand   `command:"schema"`
	Request  GenerateRequestCommand  `command:"request"`
	Response GenerateResponseCommand `command:"response"`
}

func (cmd GenerateCommand) Execute(args []string) error {
	if err := cmd.Schema.Execute(args); err != nil {
		return err
	}
	if err := cmd.Request.Execute(args); err != nil {
		return err
	}
	if err := cmd.Response.Execute(args); err != nil {
		return err
	}
	return nil
}

type GenerateSchemaCommand struct {
	Directory string `short:"d" long:"directory" default:"." description:"target directory"`
}

func (cmd GenerateSchemaCommand) Execute([]string) error {
	spec, err := openapi.LoadFile(filepath.Join(cmd.Directory, "api", "spec.yaml"))
	if err != nil {
		return err
	}
	src, err := generator.GenerateSchemaTypes(spec)
	if err != nil {
		return err
	}
	src, err = format.Source(src)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(filepath.Join(cmd.Directory, "internal", "pkg", "models", "schema_gen.go"), os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.Write(src); err != nil {
		return err
	}
	return nil
}

type GenerateRequestCommand struct {
	Directory string `short:"d" long:"directory" default:"." description:"target directory"`
}

func (cmd GenerateRequestCommand) Execute([]string) error {
	spec, err := openapi.LoadFile(filepath.Join(cmd.Directory, "api", "spec.yaml"))
	if err != nil {
		return err
	}
	src, err := generator.GenerateRequestTypes(spec)
	if err != nil {
		return err
	}
	src, err = format.Source(src)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(filepath.Join(cmd.Directory, "internal", "pkg", "models", "request_gen.go"), os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.Write(src); err != nil {
		return err
	}
	return nil
}

type GenerateResponseCommand struct {
	Directory string `short:"d" long:"directory" default:"." description:"target directory"`
}

func (cmd GenerateResponseCommand) Execute([]string) error {
	spec, err := openapi.LoadFile(filepath.Join(cmd.Directory, "api", "spec.yaml"))
	if err != nil {
		return err
	}
	src, err := generator.GenerateResponseTypes(spec)
	if err != nil {
		return err
	}
	src, err = format.Source(src)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(filepath.Join(cmd.Directory, "internal", "pkg", "models", "response_gen.go"), os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.Write(src); err != nil {
		return err
	}
	return nil
}
