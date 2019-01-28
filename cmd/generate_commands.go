package main

import (
	"log"
	"os"
	"path/filepath"

	openapi "github.com/nasa9084/go-openapi"
	"github.com/nasa9084/restgen/internal/pkg/generator"
	"github.com/pkg/errors"
)

type GenerateCommand struct {
	Directory string `short:"d" long:"directory" default:"." description:"target directory"`

	Handler  GenerateHandlerCommand  `command:"handler"`
	Schema   GenerateSchemaCommand   `command:"schema"`
	Request  GenerateRequestCommand  `command:"request"`
	Response GenerateResponseCommand `command:"response"`
}

func (cmd GenerateCommand) Execute(args []string) error {
	if err := cmd.createServerMain(); err != nil {
		return err
	}
	if err := cmd.createHTTPErr(); err != nil {
		return err
	}
	cmd.Handler.Directory = cmd.Directory
	cmd.Schema.Directory = cmd.Directory
	cmd.Request.Directory = cmd.Directory
	cmd.Response.Directory = cmd.Directory
	if err := cmd.Handler.Execute(args); err != nil {
		return err
	}
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

func (cmd GenerateCommand) createServerMain() error {
	path := filepath.Join(cmd.Directory, "cmd", "server", "main.go")
	if _, err := generator.Template("/server_main.go.tmpl", path, "%s"); err != nil {
		return errors.Wrap(err, path)
	}
	if opts.Verbose {
		log.Printf("generated %s", path)
	}
	return nil
}

func (cmd GenerateCommand) createHTTPErr() error {
	path := filepath.Join(cmd.Directory, "internal", "pkg", "httperr", "httperr_gen.go")
	if _, err := generator.Template("/httperr_httperr.go.tmpl", path); err != nil {
		return errors.Wrap(err, path)
	}
	if opts.Verbose {
		log.Printf("generated %s", path)
	}
	return nil
}

type GenerateHandlerCommand struct {
	Directory string `short:"d" long:"directory" default:"." description:"target directory"`
}

func (cmd GenerateHandlerCommand) Execute([]string) error {
	specPath := filepath.Join(cmd.Directory, "api", "spec.yaml")
	spec, err := openapi.LoadFile(specPath)
	if err != nil {
		return errors.Wrap(err, specPath)
	}
	if err := spec.Validate(); err != nil {
		return err
	}
	src, err := generator.GenerateHandlers(spec)
	if err != nil {
		return err
	}
	path := filepath.Join(cmd.Directory, "cmd", "server", "route_gen.go")
	f, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
	if err != nil {
		return errors.Wrap(err, path)
	}
	defer f.Close()
	if _, err := f.Write(src); err != nil {
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
	if err := spec.Validate(); err != nil {
		return err
	}
	src, err := generator.GenerateSchemaTypes(spec)
	if err != nil {
		return err
	}
	if len(src) == 0 {
		return nil
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
	if err := spec.Validate(); err != nil {
		return err
	}
	if err := cmd.generateRequests(spec); err != nil {
		return err
	}
	if err := cmd.generateRequestBodies(spec); err != nil {
		return err
	}
	return nil
}

func (cmd GenerateRequestCommand) generateRequests(spec *openapi.Document) error {
	src, err := generator.GenerateRequestTypes(spec)
	if err != nil {
		return err
	}
	if len(src) == 0 {
		return nil
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

func (cmd GenerateRequestCommand) generateRequestBodies(spec *openapi.Document) error {
	src, err := generator.GenerateRequestBodyTypes(spec)
	if err != nil {
		return err
	}
	if len(src) == 0 {
		return nil
	}
	f, err := os.OpenFile(filepath.Join(cmd.Directory, "internal", "pkg", "models", "request_body_gen.go"), os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
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
	if len(src) == 0 {
		return nil
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
