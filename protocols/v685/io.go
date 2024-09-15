package v685

import (
	v686 "github.com/oomph-ac/new-mv/protocols/v686"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

type Reader struct {
	*v686.Reader
}

func NewReader(r *protocol.Reader) *Reader {
	return &Reader{v686.NewReader(r)}
}

type Writer struct {
	*v686.Writer
}

func NewWriter(w *protocol.Writer) *Writer {
	return &Writer{v686.NewWriter(w)}
}
