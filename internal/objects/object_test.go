package objects_test

import (
	"cgit/internal/objects"
	"cgit/internal/utils"
	"reflect"
	"testing"
)

func TestKeyValueParser(t *testing.T) {
	raw := []byte("key value\n")
	dict := utils.NewOrderedHashMap()
	dict, err := objects.KeyValueParser(raw, 0, dict)

	if err != nil {
        t.Error(err)
        return
	}

	value, ok := dict.Get("key")

	if !ok {
		t.Error("Key should have some value")
		return
	}

	if !reflect.DeepEqual(value, []byte("value")) {
		t.Error("Key should have liternal 'value'")
		return
	}
}

func TestKeyValueSerializer (t *testing.T) {
	dict := utils.NewOrderedHashMap()

    dict.Set("hello", []byte("world"))
    dict.Set(objects.REMAINER, []byte("this is test message"))

    data := objects.KeyValueSerialize(dict)

    if !reflect.DeepEqual(data, []byte("hello world\n\nthis is test message\n")) {
        t.Error("Did not serialize properly: ", string(data))
        return
    }
}
