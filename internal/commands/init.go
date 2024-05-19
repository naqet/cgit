package commands

import "github.com/naqet/cgit/internal/repository"

type InitCommand struct {
	command Command
}

func (c *InitCommand) Process(args []string) error {
	path := "."

	if len(args) > 0 {
		path = args[0]
	}
	_, err := repository.InitRepository(path)

	return err
}

var initCmd = InitCommand{
	Command{
		help: "Initialize a git repository",
		args: []Arg{
			{
				name:         "path",
				defaultValue: ".",
				help:         "Path where to create a git repository",
			},
		},
	},
}
