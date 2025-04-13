package iox

import (
	"bytes"
	"io"
)

// StripBOM removes UTF-8 BOM from the reader if present.
func StripBOM(r io.Reader) io.Reader {
	bom := make([]byte, 3)
	n, err := r.Read(bom)
	if err != nil || n < 3 || !bytes.Equal(bom, []byte{0xEF, 0xBB, 0xBF}) {
		return io.MultiReader(bytes.NewReader(bom[:n]), r)
	}
	return r
}
