package span

import "encoding/hex"

type ID struct {
	inner [8]byte
}

func NewSpanID(id [8]byte) ID {
	return ID{id}
}

func (t ID) String() string {
	return hex.EncodeToString(t.inner[:])
}

func (t ID) Bytes() []byte {
	return t.inner[:]
}
