package commands

import (
	"cgit/internal/cargs"
	"cgit/internal/objects"
	"fmt"
	"path/filepath"
)

type LsTreeCmd struct{}

func (c LsTreeCmd) Name() string {
	return "ls-tree"
}

func (c LsTreeCmd) Help() string {
	return "This is help message for ls-tree"
}

func (c LsTreeCmd) Process(args *cargs.Args) error {
	if len(args.Values) < 1 {
		return fmt.Errorf("Required at least 1 argument")
	}

	isRecursive := false

	if args.Values[0] == "-r" {
		return fmt.Errorf("Specify tree")
	}

	if args.Values[0] == "-r" {
		isRecursive = true
	}
    sha := ""

    if isRecursive {
        sha = args.Values[1]
    } else {
        sha = args.Values[0]
    }

	repo, err := objects.FindRepository(".")

	if err != nil {
		return err
	}

	return lsTree(repo, isRecursive, sha, "")
}

func lsTree(repo *objects.Repository, isRecursive bool, sha string, prefix string) error {
	obj, err := repo.ReadObject([]byte(sha))

	if err != nil {
		return err
	}
	tree, ok := obj.(*objects.Tree)

	if !ok {
		return fmt.Errorf("Object is not a tree")
	}

	for _, leaf := range tree.Data {
		entryType := ""
		if len(leaf.Mode) == 5 {
			entryType = string(leaf.Mode[:1])
		} else {
			entryType = string(leaf.Mode[:2])
		}

		switch entryType {
		case "04":
			entryType = "tree"
		case "10":
		case "12":
            entryType = "blob"
		case "16":
            entryType = "commit"
        default: 
            panic("Invalid leaf")
		}

        if !(isRecursive && entryType == "tree") {
            fmt.Printf("%s %s %s\t%s", leaf.Mode, entryType, leaf.SHA, filepath.Join(prefix, string(leaf.Path)))
        } else {
            return lsTree(repo, isRecursive, sha, filepath.Join(prefix, string(leaf.Path)))
        }
	}
    return nil
}
 
