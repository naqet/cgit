package objects

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"os"
	"path/filepath"
	"strconv"
)

type Object interface {
	GetType() []byte
	Serialize() []byte
	Deserialize([]byte)
}

func WriteObject(obj Object, repo *Repository) (string, error) {
	data := obj.Serialize()
	objType := obj.GetType()

	res := bytes.NewBuffer(append(objType, ' '))

	length := []byte(strconv.Itoa(len(data)))

	_, err := res.Write(append(length, '\x00'))
	if err != nil {
		return "", err
	}

	_, err = res.Write(data)

	if err != nil {
		return "", err
	}

	hasher := sha1.New()
	_, err = hasher.Write(res.Bytes())

	if err != nil {
		return "", err
	}

	sha := hex.EncodeToString(hasher.Sum(nil))

	if repo != nil {
		path := repo.GetPath("objects", sha[0:2], sha[2:])

		_, err := os.Stat(path)

		if os.IsNotExist(err) {
            if err = os.MkdirAll(filepath.Dir(path), mode); err != nil {
                return "", err
            }

            file, err := os.Create(path)

			if err != nil {
				return "", err
			}
			defer file.Close()

			writer := zlib.NewWriter(file)
			defer writer.Close()

			writer.Write(res.Bytes())
		} else if err != nil {
			return "", err
		}
	}

	return sha, nil
}
