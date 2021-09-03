package main

import (
	"git.jobtome.io/auxiliary/docker-registry-cleaner/pkg/cmd"
	"github.com/jessevdk/go-flags"
	"os"
)

var GitVersion string

func main() {
	if GitVersion == "" {
		GitVersion = "develop"
	}

	app := &cmd.AppCommand{}
	parser := flags.NewParser(app, flags.Default)
	cmd.RegisterCleanCommand(parser)

	_, err := parser.ParseArgs(os.Args[1:])
	if err != nil {
		os.Exit(-1)
	}
}
