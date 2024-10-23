package utils_test

import (
	"cgit/internal/utils"
	"reflect"
	"testing"
)

func TestOrderedHashMapSet(t *testing.T) {
    omap := utils.NewOrderedHashMap();
    omap.Set("hello", []byte("world"))

    data, ok := omap.Get("hello")

    if !ok || !reflect.DeepEqual(data, []byte("world")) {
        t.Error("hello did not return world")
        return
    }

    if omap.Len() != 1 {
        t.Error("Length is not correct")
        return
    }

    omap.Set("one", []byte("two"))

    if omap.Len() != 2 {
        t.Error("Length is not correct")
        return
    }

    omap.Set("hello", []byte("changed"))

    if omap.Len() != 2 {
        t.Error("Length is not correct")
        return
    }

    data, ok = omap.Get("hello")

    if !ok || !reflect.DeepEqual(data, []byte("changed")) {
        t.Error("hello did not return another")
        return
    }

    data, ok = omap.Get("not-existent")

    if ok {
        t.Error("this key should not exist")
        return
    }
}
