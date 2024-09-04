package v685

import "github.com/sandertv/gophertunnel/minecraft/protocol"

type Reader struct {
	*protocol.Reader
}

func NewReader(r *protocol.Reader) *Reader {
	return &Reader{r}
}

type Writer struct {
	*protocol.Writer
}

func NewWriter(w *protocol.Writer) *Writer {
	return &Writer{w}
}
