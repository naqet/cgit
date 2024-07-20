package objects

import (
	"bufio"
	"compress/zlib"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"

	"gopkg.in/ini.v1"
)

type Repository struct {
	Worktree string
	Gitdir   string
	Conf     *ini.File
}

const mode = 0777

func InitRepository(path string) (*Repository, error) {
    gitDirPath := filepath.Join(path, ".cgit")
	_, err := os.ReadDir(gitDirPath)

	if err == nil {
		// TODO: allow force reinit
		return nil, fmt.Errorf("cgit already initialized")
	} else if !os.IsNotExist(err) {
		return nil, err
	}

	repo := &Repository{
		Worktree: path,
		Gitdir:   gitDirPath,
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

	defer file.Close()
	_, err = config.WriteTo(file)

	return repo, err
}

func GetRepository(path string) (*Repository, error) {
    gitDirPath := filepath.Join(path, ".cgit")
	_, err := os.ReadDir(gitDirPath)

	if os.IsNotExist(err) {
        return nil, fmt.Errorf("Invalid path")
    } else if err != nil {
        return nil, err
    }

	repo := &Repository{
		Worktree: path,
		Gitdir:   gitDirPath,
	}

    config, err := ini.Load(repo.GetPath("config"))

    if err != nil {
        return nil, fmt.Errorf("Invalid config file: %s", err)
    }

    repo.Conf = config

    return repo, nil
}

func FindRepository(args ...string) (*Repository, error) {
    path := "."

    if len(args) == 1 {
        path = args[0]
    }
    _, err := os.Stat(filepath.Join(path, ".cgit"))

    if err == nil {
        return GetRepository(path)
    }

    parentPath := filepath.Join(path, "..")

    // Base case "/" path
    if parentPath == path {
        return nil, fmt.Errorf("No .cgit dir")
    }

    return FindRepository(parentPath)
}

func (r *Repository) ReadObject(sha []byte) (Object, error) {
	if len(sha) < 2 {
		return nil, fmt.Errorf("Invalid object hash")
	}

	path := r.GetPath("objects", string(sha[0:2]), string(sha[2:]))

	file, err := os.Open(path)

	if os.IsNotExist(err) {
		return nil, fmt.Errorf("Invalid object hash")
	} else if err != nil {
		return nil, err
	}

	defer file.Close()
	rd, err := zlib.NewReader(file)

	if err != nil {
		return nil, err
	}
	defer rd.Close()

	read := bufio.NewReader(rd)

	objType, err := read.ReadBytes(' ')

	if err != nil {
		return nil, err
	}

	sizeBytes, err := read.ReadBytes('\x00')

	if err != nil {
		return nil, err
	}

    size, err := strconv.Atoi(string(sizeBytes[:len(sizeBytes)-1]))

	if err != nil {
		return nil, err
	}

    data, err := io.ReadAll(read)

	if err != nil {
		return nil, err
	}

	if size != len(data) {
		return nil, fmt.Errorf("Malformed object: bad length")
	}

	var obj Object

    switch string(objType[:len(objType)-1]) {
	case "blob":
        obj = &Blob{Data: data}
	//case "commit":
	//	obj = Blob{}
	//case "tag":
	//	obj = Blob{}
	//case "tree":
	//	obj = Blob{}
	}

	return obj, nil
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
