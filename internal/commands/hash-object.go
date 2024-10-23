package commands

import (
	"cgit/internal/cargs"
	"cgit/internal/objects"
	"fmt"
	"os"
	"slices"
)

type HashObjectCmd struct{}

func (c HashObjectCmd) Name() string {
	return "hash-object"
}

func (c HashObjectCmd) Help() string {
	return "This is help message for hash object"
}

func (c HashObjectCmd) Process(args *cargs.Args) error {
    isWrite := false;
    objType := "blob"
    path := ""

    for i, val := range args.Values {
        if i > 0 && args.Values[i-1] == "-t" {
            types := []string{"blob", "commit", "tag", "tree"}
            idx := slices.Index(types, val)

            if idx == -1 {
                return fmt.Errorf("Invalid type")
            }

            objType = val
            continue
        }

        if val == "-w" {
            isWrite = true
            continue
        }

        path = val
    }

    data, err := os.ReadFile(path)

    if os.IsNotExist(err) {
        return fmt.Errorf("Invalid path")
    } else if err != nil {
        return err
    }

    var repo *objects.Repository

    if isWrite {
        repo, err = objects.FindRepository()

        if err != nil {
            return err
        }
    }

    var obj objects.Object

    switch objType {
    case "blob":
        obj = &objects.Blob{Data: data}
    }

    sha, err := objects.WriteObject(obj, repo)

    if err != nil {
        return err
    }

    fmt.Println(sha)

	return nil
}
