package commands

import (
	"cgit/internal/cargs"
	"cgit/internal/objects"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type CheckoutCmd struct{}

func (c CheckoutCmd) Name() string {
	return "checkout"
}

func (c CheckoutCmd) Help() string {
	return "This is help message for checkout"
}

func (c CheckoutCmd) Process(args *cargs.Args) error {
	if len(args.Values) != 2 {
		return fmt.Errorf("Requires 2 arguments. Commit or tree hash and directory")
	}

	commit := args.Values[0]
	path := args.Values[1]

	repo, err := objects.FindRepository(".")

	if err != nil {
		return err
	}

	obj, err := repo.ReadObject([]byte(commit))

	if err != nil {
		return err
	}

	data, ok := obj.(*objects.Commit)

	if !ok {
		return fmt.Errorf("Invalid commit sha")
	}

    treeSha, ok := data.Data.Get("tree")

	if !ok {
		return fmt.Errorf("Commit does not have tree")
	}

	info, err := os.Stat(path)

	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("Error while checking path: %s", err.Error())
	}

	if errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(path, 0644)

		if err != nil {
			return err
		}
	} else if !info.IsDir() {
        return fmt.Errorf("Path does not lead to dir")
	} else {
        dir, err := os.ReadDir(path)

        if err != nil {
            return err
        }

        if len(dir) != 0 {
            return fmt.Errorf("Directory is not empty")
        }
    }

	return treeCheckout(repo, treeSha, path)
}

func treeCheckout(repo *objects.Repository, treeSha []byte, dir string) error {
    obj, err := repo.ReadObject(treeSha)

    if err != nil {
        return err
    }

    tree, ok := obj.(*objects.Tree)

    if !ok {
        return fmt.Errorf("Object is not a tree")
    }

    for _, leaf := range tree.Data {
        leafObj, err := repo.ReadObject([]byte(leaf.Path))

        if err != nil {
            fmt.Println("Error while accessing one of the tree leaves at: ", leaf.Path)
            continue;
        }

        destination := filepath.Join(dir, string(leaf.Path))

        if string(leafObj.GetType()) == "tree" {
            os.MkdirAll(destination, 0644);
            treeCheckout(repo, leaf.SHA, destination)
        } else if string(leafObj.GetType()) == "blob" {
            data := leafObj.Serialize()
            fmt.Println(string(data))
        }
    }

    return nil
}
