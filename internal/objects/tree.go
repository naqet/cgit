package objects

import (
	"bufio"
	"bytes"
	"encoding/binary"
)

type Leaf struct {
	Mode []byte
	SHA  []byte
	Path []byte
}

func parseLeaf(raw []byte, start int) (Leaf, int) {
	modeIdx := bytes.IndexByte(raw[start:], ' ')

	mode := raw[start:modeIdx]

	pathIdx := bytes.IndexByte(raw[modeIdx+1:], '\x00')

	path := raw[modeIdx+1 : pathIdx]

	sha := raw[pathIdx+1 : pathIdx+21]

	return Leaf{mode, sha, path}, pathIdx + 21
}

func parseTreeData(raw []byte) []Leaf {
	idx := 0
	maxIdx := len(raw)
	res := []Leaf{}

	for idx < maxIdx {
		data, pos := parseLeaf(raw, idx)
		res = append(res, data)
		idx = pos
	}

	return res
}

type Tree struct {
	Data []Leaf
}

func (b *Tree) GetType() []byte {
	return []byte("tree")
}

func (b *Tree) Serialize() []byte {
    writer := bufio.NewWriter(&bytes.Buffer{})
    for _, leaf := range b.Data {
        binary.Write(writer, binary.LittleEndian, leaf.Path)
        writer.WriteRune(' ')
        writer.Write(leaf.Path)
        writer.WriteRune('\x00')
        writer.Write(leaf.SHA)
    }
	return []byte{}
}

func (b *Tree) Deserialize(data []byte) {
    b.Data = parseTreeData(data)
}
