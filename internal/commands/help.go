package commands

import (
	"cgit/internal/cargs"
	"fmt"
)

type HelpCmd struct{ }

func (c HelpCmd) Name() string {
	return "help"
}

func (c HelpCmd) Help() string {
	return "This is help message"
}

func (c HelpCmd) Process(_ *cargs.Args) error {
	fmt.Println(c.Help())
	return nil
}
