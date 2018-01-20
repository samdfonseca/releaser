package flags

import (
	"os"
	"gopkg.in/urfave/cli.v2"
)

var (
	Config = &cli.StringFlag{
		Name: "config",
		Aliases: []string{"c"},
		Value: os.ExpandEnv("${HOME}/.releaser/config.json"),
		Usage: "path to the config.json file",
		EnvVars: []string{"RELEASER_CONFIG"},
	}
)

