package objects

import (
	"cgit/internal/utils"
	"fmt"
)

type Commit struct {
	Data *utils.OrderedHashMap
}

func NewCommit() *Commit {
    return &Commit{utils.NewOrderedHashMap()}
}

func (b *Commit) GetType() []byte {
	return []byte("commit")
}

func (b *Commit) Serialize() []byte {
	return KeyValueSerialize(b.Data)
}

func (b *Commit) Deserialize(data []byte) {
    _, err := KeyValueParser(data, 0, b.Data)
    fmt.Println("Error while parsing the data in commit object", err)
}
