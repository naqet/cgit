package objects

type Blob struct {
	_type []byte
	data  []byte
}

func NewBlob(data []byte) *Blob {
	return &Blob{[]byte("blob"), data}
}

func (b *Blob) GetType() []byte {
	return b._type
}

func (b *Blob) Serialize() []byte {
	return b.data
}

func (b *Blob) Deserialize(data []byte) {
	b.data = data
}
