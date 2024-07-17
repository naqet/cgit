package cargs

import (
	"fmt"
	"os"
)

type Command interface {
	Name() string
    Help() string
	Process(*Args) error
}

type Args struct {
	Commands map[string]Command
    Values []string
}

func InitArgs() *Args {
	return &Args{
		Commands: map[string]Command{},
	}
}

func (a *Args) Process() (Command, error) {
	args := os.Args[1:]

	if len(args) == 0 {
		if cmd, ok := a.Commands["help"]; ok {
			return cmd, nil
		}
		return nil, fmt.Errorf("Invalid command")
	}

	cmdName := args[0]

	cmd, ok := a.Commands[cmdName]

	if !ok {
		return nil, fmt.Errorf("Invalid command")
	}
    a.Values = args[1:]

	return cmd, nil
}

func (a *Args) AddCommand(cmd Command) {
	a.Commands[cmd.Name()] = cmd
}
