package commands

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/naqet/cgit/internal/objects"
	"github.com/naqet/cgit/internal/repository"
)

type CatFileCommand struct {
	Command
}

var catFileCmd = CatFileCommand{
	Command{
		help: "Provide content of repository objects",
		args: []Arg{
			{
				name:    "type",
				help:    "Path where to create a git repository",
				choices: []string{"blob", "commit", "tag", "tree"},
			},
			{
				name: "object",
				help: "The object to display",
			},
		},
	},
}

func (c *CatFileCommand) Process(args []string) error {
	if len(args) < 2 {
		return errors.New("Command needs 2 arguments")
	}

    types := c.args[0].choices
	if slices.Index(types, args[0]) == -1 {
		return errors.New("Type needs to be one of:\n" + strings.Join(types, "\n"))
	}

	repo, err := repository.FindRepository(".")

	if err != nil {
		return err
	}

	sha := args[1]

	data, err := objects.ReadObject(sha, repo)

	if err != nil {
		return err
	}

	fmt.Println(string(data.Serialize()))

	return nil
}
