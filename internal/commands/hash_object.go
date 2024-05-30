package commands

import (
	//"errors"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/naqet/cgit/internal/objects"
	"github.com/naqet/cgit/internal/repository"
)

type HashObjectCommand struct {
	Command
}

var hashObjectCmd = HashObjectCommand{
	Command{
		help: "Compute object ID and optionally create a blob from a file",
		args: []Arg{
			{
				name: "path",
				help: "Read object from <file> at path",
			},
		},
	},
}

func (c *HashObjectCommand) Process(args []string) error {
	//TODO implement utils args validator
	if len(args) < 1 {
		return errors.New("Path to file not defined")
    }
    //TODO: implement proper flag parsing
	objType := flag.String("t", "blob", "Specify the type")
	write := flag.Bool("w", false, "Write the object into the database")

	flag.Parse()

	var repo *repository.Repository
	if *write {
		r, err := repository.FindRepository(".")

		if err != nil {
			return err
		}
		repo = r
	}

	content, err := os.ReadFile(args[0])

	if err != nil {
		return err
	}

	var obj objects.Object
    //TODO create proper utils function for object creation
	switch *objType {
	case "blob":
		obj = objects.NewBlob(content)
	}

    sha, err := objects.WriteObject(obj, repo)

    if err != nil {
        return err
    }

    fmt.Println(sha)
    
	return nil
}
