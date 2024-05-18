package objects

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"os"
	"path"
	"strconv"

	"github.com/naqet/cgit/internal/repository"
)

type Object interface {
	GetType() []byte
	Serialize() []byte
	Deserialize([]byte)
}

func WriteObject(obj Object, repo *repository.Repository) (string, error) {
	data := obj.Serialize()

	var buf bytes.Buffer

	_, err := buf.Write(obj.GetType())
	if err != nil {
		return "", err
	}

	_, err = buf.WriteRune(' ')
	if err != nil {
		return "", err
	}
	_, err = buf.Write([]byte(strconv.Itoa(len(data))))
	if err != nil {
		return "", err
	}
	_, err = buf.WriteRune('\x00')
	if err != nil {
		return "", err
	}

	_, err = buf.Write(data)
	if err != nil {
		return "", err
	}

	if err != nil {
		return "", err
	}

	hash := sha1.New()

	_, err = hash.Write(buf.Bytes())

	if err != nil {
		return "", err
	}
	hashBytes := hash.Sum(nil)

	sha := hex.EncodeToString(hashBytes)

	if repo != nil {
		r := *repo
		dir := sha[:2]
		file := sha[2:]

		err = os.Mkdir(path.Join(r.Path, "objects", dir), 0777)

		if err != nil {
			return "", err
		}

		f, err := os.Create(path.Join(r.Path, "objects", dir, file))

        if err != nil {
            return "", err
        }

        defer f.Close()

        w := zlib.NewWriter(f)

        _, err = w.Write(buf.Bytes())

        if err != nil {
            return "", err
        }

        err = w.Close()

        if err != nil {
            return "", err
        }
	}

	return sha, nil
}
