package objects;

type Blob struct {
    Data []byte
}

func (b *Blob) GetType() []byte {
    return []byte("blob")
}

func (b *Blob) Serialize() []byte {
    return b.Data
}

func (b *Blob) Deserialize(data []byte) {
    b.Data = data
}
