package objects

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"io"
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
func ReadObject(sha string, repo *repository.Repository) (Object, error) {
	if repo == nil {
		return nil, errors.New("Repository is not available")
	}

	if len(sha) < 3 {
		return nil, errors.New("Invalid hash")
	}

	dir := sha[:2]
	file := sha[2:]

	f, err := os.Open(path.Join(repo.Path, "objects", dir, file))

	if err != nil {
		return nil, err
	}

	defer f.Close()

	r, err := zlib.NewReader(f)

	if err != nil {
		return nil, err
	}

	hashData := bytes.Buffer{}

	_, err = io.Copy(&hashData, r)

	r.Close()

	if err != nil {
		return nil, err
	}

	objType, err := hashData.ReadBytes(' ')

	if err != nil {
		return nil, err
	}

	_, err = hashData.ReadBytes('\x00')

	if err != nil {
		return nil, err
	}

    data := bytes.Buffer{}

    var b byte

    for {
        b, err = hashData.ReadByte()
        
        if err != nil {
            break;
        }

        err = data.WriteByte(b)

        if err != nil {
            break
        }
    }

    if err != nil && !errors.Is(err, io.EOF) {
        return nil, err
    }

	if len(objType) < 2 {
		return nil, errors.New("Invalid object type")
	}

	var obj Object

	switch string(objType[:len(objType)-1]) {
	case "blob":
		obj = NewBlob(data.Bytes())
	}

	return obj, nil
}
