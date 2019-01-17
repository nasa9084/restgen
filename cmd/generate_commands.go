package main

import (
	"os"
	"path/filepath"

	openapi "github.com/nasa9084/go-openapi"
	"github.com/nasa9084/restgen/internal/generator"
)

type GenerateCommand struct {
	Directory string `short:"d" long:"directory" default:"." description:"target directory"`

	Route    GenerateRouteCommand    `command:"route"`
	Schema   GenerateSchemaCommand   `command:"schema"`
	Request  GenerateRequestCommand  `command:"request"`
	Response GenerateResponseCommand `command:"response"`
}

func (cmd GenerateCommand) Execute(args []string) error {
	cmd.Route.Directory = cmd.Directory
	cmd.Schema.Directory = cmd.Directory
	cmd.Request.Directory = cmd.Directory
	cmd.Response.Directory = cmd.Directory
	if err := cmd.Route.Execute(args); err != nil {
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

type GenerateRouteCommand struct {
	Directory string `short:"d" long:"directory" default:"." description:"target directory"`
}

func (cmd GenerateRouteCommand) Execute([]string) error {
	spec, err := openapi.LoadFile(filepath.Join(cmd.Directory, "api", "spec.yaml"))
	if err != nil {
		return err
	}
	src, err := generator.GenerateRoutes(spec)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(filepath.Join(cmd.Directory, "cmd", "server", "route_gen.go"), os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
	if err != nil {
		return err
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
	src, err := generator.GenerateSchemaTypes(spec)
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
