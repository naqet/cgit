package utils;

type entry struct {
	value []byte
	order int
}

type OrderedHashMap struct {
	data map[string]entry
    absCount int
}

func NewOrderedHashMap() *OrderedHashMap {
    return &OrderedHashMap{map[string]entry{}, 0}
}

func (o *OrderedHashMap) Keys() []string {
    keys := make([]string, o.Len());

    // TODO: this works only when there is not delete method on orderedHashMap
    for key, entry := range o.data {
        keys[entry.order] = key
    }

    return keys;
}

func (o *OrderedHashMap) Len() int {
    return len(o.data)
}

func (o *OrderedHashMap) Get(key string) ([]byte, bool) {
    value, ok := o.data[key]

    if !ok {
        return []byte{}, false
    }

    return value.value, ok
}

func (o *OrderedHashMap) Set(key string, value []byte) {
    ent, ok := o.data[key]

    if ok {
        o.data[key] = entry{value, ent.order}
        return
    }
    o.data[key] = entry{value, o.absCount}
    o.absCount++;
}
