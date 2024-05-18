package commands

import (
	"github.com/naqet/cgit/internal/repository"
)

type Arg struct {
	name         string
	defaultValue string
	help         string
}

var Commands =  map[string]*Command{
    "init": &initCmd,
}

type Command struct {
	name    string
	help    string
	args    []Arg
	Process func([]string) error
}

var initCmd Command = Command{
	help: "Initialize a git repository",
	args: []Arg{
		{
			name:         "path",
			defaultValue: ".",
			help:         "Path where to create a git repository",
		},
	},
	Process: func(args []string) error {
		path := "."

		if len(args) > 0 {
			path = args[0]
		}
		_, err := repository.InitRepository(path)

		return err
	},
}
