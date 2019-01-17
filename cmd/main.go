package main

import (
	"fmt"
	"log"
	"os"

	flags "github.com/jessevdk/go-flags"
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
