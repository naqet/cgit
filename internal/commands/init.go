package commands

import (
	"cgit/internal/cargs"
	"cgit/internal/objects"
)

type InitCmd struct{}

func (c InitCmd) Name() string {
	return "init"
}

func (c InitCmd) Help() string {
    //TODO: update help message
	return "Init help"
}

func (c InitCmd) Process(args *cargs.Args) error {
    path := "."
	if len(args.Values) > 0 {
		path = args.Values[0]
	}
    _, err := objects.InitRepository(path)
    if err != nil {
        return err
    }

	return nil
}
