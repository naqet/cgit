package repository

import (
	"errors"
	"io/fs"
	"os"

	"gopkg.in/ini.v1"
)

type Repository struct {
	root   string
	path   string
	dir    *[]fs.DirEntry
	config *ini.File
}

func InitRepository(path string) (*Repository, error) {
	gitPath := path + "/.cgit"
	dir, err := os.ReadDir(gitPath)

	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}

	repo := Repository{path, gitPath, &dir, nil}

	err = repo.initConfig(errors.Is(err, os.ErrNotExist))

	if err != nil {
		return nil, err
	}

    err = repo.createDirs()
    if err != nil {
        return nil, err
    }

    file, err := os.Create(repo.path + "/description")
    defer file.Close()

    if err != nil {
        return nil, err
    }

    file.Write([]byte("Unnamed repository; edit this file 'description' to name the repository.\n"))

    file, err = os.Create(repo.path + "/HEAD")
    defer file.Close()

    if err != nil {
        return nil, err
    }

    file.Write([]byte("ref: refs/heads/master\n"))

	return &repo, nil
}

func (r *Repository) initConfig(newRepo bool) error {
	var configFile *ini.File
	var err error

	if newRepo {
		err = os.Mkdir(r.path, 0777)

		if err != nil {
			return err
		}

		configFile = ini.Empty()
		section, err := configFile.NewSection("core")

		if err != nil {
			return err
		}

		_, err = section.NewKey("repositoryformatversion", "0")
		if err != nil {
			return err
		}

		_, err = section.NewKey("filemode", "false")
		if err != nil {
			return err
		}

		_, err = section.NewKey("bare", "false")
		if err != nil {
			return err
		}

		err = configFile.SaveTo(r.path + "/config")
	} else {
		if *r.dir == nil {
			return errors.New("Error reading .cgit dir")
		}

		for _, entry := range *r.dir {
			if entry.Name() == "config" {
				configFile, err = ini.Load(r.path + "/" + entry.Name())
				break
			}
		}

		if err != nil {
			return err
		}

		if configFile == nil {
			return errors.New("Config file not found")
		}
	}

	r.config = configFile

	return err
}

func (r *Repository) createDirs() error {
	dirs := []string{"/refs/heads", "/refs/tags", "/objects", "/branches"}
	var err error

	for _, dir := range dirs {
		err = os.MkdirAll(r.path+dir, 0777)
		if err != nil {
			break
		}
	}

	return err
}
