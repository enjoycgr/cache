package app

type ByteView struct {
	B []byte
}

func (v ByteView) Len() int {
	return len(v.B)
}

func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.B)
}

func (v ByteView) String() string {
	return string(v.B)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
