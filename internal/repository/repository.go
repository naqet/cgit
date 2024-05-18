package repository

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"gopkg.in/ini.v1"
)

type Repository struct {
	root   string
	path   string
	config *ini.File
}

func InitRepository(path string) (*Repository, error) {
	root, gitPath, err := getPaths(path)
	if err != nil {
		return nil, err
	}

	_, err = os.ReadDir(gitPath)

	if err == nil {
		return nil, errors.New("Git repository already exists")
	}

	if !os.IsNotExist(err) {
		return nil, err
	}

	repo := Repository{root, gitPath, nil}

	err = repo.initConfig()

	if err != nil {
		return nil, err
	}

	err = repo.initRequiredDirs()
	if err != nil {
		return nil, err
	}

	err = repo.initRequiredFiles()
	if err != nil {
		return nil, err
	}

	return &repo, nil
}

func FindRepository(path string) (*Repository, error) {
	root, gitPath, err := getPaths(path)

	if err != nil {
		return nil, err
	}

	dir, err := os.ReadDir(gitPath)

	if err == nil {
		configFile, err := getConfigFile(&dir, gitPath)

		if err != nil {
			return nil, err
		}

		return &Repository{root, gitPath, configFile}, nil
	}

	parent := filepath.Dir(root)

	if parent == path {
		return nil, errors.New("No git repository")
	}

	return FindRepository(parent)
}

func (r *Repository) initConfig() error {
	var configFile *ini.File

	err := os.Mkdir(r.path, 0777)

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

	if err != nil {
		return err
	}

	r.config = configFile

	return err
}

func (r *Repository) initRequiredDirs() error {
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

func (r *Repository) initRequiredFiles() error {
	file, err := os.Create(r.path + "/description")
	defer file.Close()

	if err != nil {
		return err
	}

	_, err = file.Write([]byte("Unnamed repository; edit this file 'description' to name the repository.\n"))

	if err != nil {
		return err
	}
	file, err = os.Create(r.path + "/HEAD")
	defer file.Close()

	if err != nil {
		return err
	}

	_, err = file.Write([]byte("ref: refs/heads/master\n"))

	return err
}

func getPaths(path string) (root string, gitPath string, err error) {
	root, err = filepath.Abs(path)
	if err != nil {
		return "", "", err
	}

	gitPath = root + "/.cgit"

	return root, gitPath, nil
}

func getConfigFile(dir *[]fs.DirEntry, gitPath string) (configFile *ini.File, err error) {
	for _, entry := range *dir {
		if entry.Name() == "config" {
			configFile, err = ini.Load(gitPath + "/" + entry.Name())
			break
		}
	}

	if err != nil {
		return nil, err
	}

	if configFile == nil {
		return nil, errors.New("Config file not found")
	}

	return configFile, nil
}
