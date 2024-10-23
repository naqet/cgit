package objects

import (
	"bytes"
	"cgit/internal/utils"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

const REMAINER string = "remainer"

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

func KeyValueParser(raw []byte, start int, dict *utils.OrderedHashMap) (*utils.OrderedHashMap, error) {
	space := bytes.IndexByte(raw[start:], ' ')
	newLine := bytes.IndexByte(raw[start:], '\n')

	if len(raw) <= start {
		return dict, nil
	}

	if newLine == -1 {
		return nil, fmt.Errorf("Invalid raw. No new line char found.")
	}

	if space < 0 || newLine < space {
		value := []byte{}
		if start+1 <= len(raw) {
			value = raw[start+1:]
		}
		dict.Set(REMAINER, value)
		return dict, nil
	}

	key := raw[start:space]
	value := raw[space+1 : newLine]

	val, ok := dict.Get(string(key))

	if ok {
		dict.Set(string(key), append(append(val, '\n'), value...))
	} else {
		dict.Set(string(key), value)
	}

	return KeyValueParser(raw, newLine+1, dict)
}

func KeyValueSerialize(dict *utils.OrderedHashMap) []byte {
	buf := bytes.NewBuffer([]byte{})
	keys := dict.Keys()

	for _, key := range keys {
		if key == REMAINER {
			continue
		}

		unprocessedValue, ok := dict.Get(key)

		if !ok {
			fmt.Println("Invalid key: ", key)
			continue
		}

		values := bytes.Split(unprocessedValue, []byte("\n"))

		for _, val := range values {
			value := bytes.ReplaceAll(val, []byte("\n"), []byte("\n "))
			buf.WriteString(key)
			buf.WriteByte(' ')
			buf.Write(value)
			buf.WriteByte('\n')
		}
	}

    remainer, ok := dict.Get(REMAINER)

    if !ok {
        return buf.Bytes()
    }
    buf.WriteByte('\n')
    buf.Write(remainer)
    buf.WriteByte('\n')

    return buf.Bytes();
}
