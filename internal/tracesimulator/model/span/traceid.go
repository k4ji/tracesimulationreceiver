package span

import "encoding/hex"

type TraceID struct {
	inner [16]byte
}

func NewTraceID(id [16]byte) TraceID {
	return TraceID{id}
}

func (t TraceID) String() string {
	return hex.EncodeToString(t.inner[:])
}

func (t TraceID) Bytes() []byte {
	return t.inner[:]
}
