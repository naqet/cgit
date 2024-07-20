package commands

import (
	"cgit/internal/cargs"
	"cgit/internal/objects"
	"fmt"
)

type CatFileCmd struct{}

func (c CatFileCmd) Name() string {
	return "cat-file"
}

func (c CatFileCmd) Help() string {
	return "This is help message of cat-file"
}

func (c CatFileCmd) Process(args *cargs.Args) error {
    if len(args.Values) != 2 {
        return fmt.Errorf("Required 2 arguments")
    }

    // Type will be used in the future
    objType := ""
    for _, t := range []string{"blob", "commit", "tag", "tree"} {
        if args.Values[0] == t {
            objType = t;
            break;
        }
    }

    if objType == "" {
        return fmt.Errorf("Invalid type")
    }

    hash := args.Values[1];

    repo, err := objects.FindRepository(".")

    if err != nil {
        return err
    }

    obj, err := repo.ReadObject([]byte(hash))

    if err != nil {
        return err
    }

    fmt.Print(string(obj.Serialize()))

    return nil
}
