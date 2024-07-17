package objects

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/ini.v1"
)

type Repository struct {
	Worktree string
	Gitdir   string
	Conf     *ini.File
}

const mode = 0777

func InitRepository(path string) (*Repository, error) {
	_, err := os.ReadDir(fmt.Sprintf("%s/.cgit", path))

	if err == nil {
		// TODO: allow force reinit
		return nil, fmt.Errorf("cgit already initialized")
	} else if !os.IsNotExist(err) {
		return nil, err
	}

	repo := &Repository{
		Worktree: path,
		Gitdir:   filepath.Join(path, ".cgit"),
	}

	err = os.MkdirAll(repo.GetPath("branches"), mode)
	if err != nil {
		return nil, err
	}

	err = os.MkdirAll(repo.GetPath("objects"), mode)
	if err != nil {
		return nil, err
	}

	err = os.MkdirAll(repo.GetPath("refs", "tags"), mode)
	if err != nil {
		return nil, err
	}

	err = os.MkdirAll(repo.GetPath("refs", "heads"), mode)
	if err != nil {
		return nil, err
	}

	err = os.WriteFile(
		repo.GetPath("description"),
		[]byte("Unnamed repository; edit this file 'description' to name the repository.\n"),
		0666,
	)
	if err != nil {
		return nil, err
	}

	err = os.WriteFile(
		repo.GetPath("HEAD"),
		[]byte("ref: refs/heads/master\n"),
		0666,
	)
	if err != nil {
		return nil, err
	}

    config, err := getDefaultConfig()

    if err != nil {
        return nil, err
    }

    file, err := os.Create(repo.GetPath("config"))

    if err != nil {
        return nil, err
    }
    _, err = config.WriteTo(file)

	return repo, err
}

func (r *Repository) GetPath(paths ...string) string {
	return filepath.Join(r.Gitdir, filepath.Join(paths...))
}

func getDefaultConfig() (file *ini.File, err error) {
    file = ini.Empty()

    section, err := file.NewSection("core")

    if err != nil {
        return nil, err
    }

    _, err = section.NewKey("repositoryformatversion", "0")
    if err != nil {
        return nil, err
    }

    _, err = section.NewKey("filemode", "false")
    if err != nil {
        return nil, err
    }

    _, err = section.NewKey("bare", "false")
    if err != nil {
        return nil, err
    }

    return file, nil
}
