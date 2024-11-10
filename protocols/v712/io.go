package v712

import (
	"fmt"

	"github.com/oomph-ac/new-mv/protocols/v712/types"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

type Reader struct {
	*protocol.Reader
}

func NewReader(r *protocol.Reader) *Reader {
	return &Reader{r}
}

func (*Reader) Reads() bool {
	return true
}

// StackRequestAction reads a StackRequestAction from the reader.
func (r *Reader) StackRequestAction(x *protocol.StackRequestAction) {
	var id uint8
	r.Uint8(&id)
	if !types.LookupStackRequestAction(id, x) {
		r.UnknownEnumOption(id, "stack request action type")
		return
	}
	(*x).Marshal(r)
}

type Writer struct {
	*protocol.Writer
}

func NewWriter(w *protocol.Writer) *Writer {
	return &Writer{w}
}

func (*Writer) Writes() bool {
	return true
}

// StackRequestAction writes a StackRequestAction to the writer.
func (w *Writer) StackRequestAction(x *protocol.StackRequestAction) {
	var id byte
	if !types.LookupStackRequestActionType(*x, &id) {
		w.UnknownEnumOption(fmt.Sprintf("%T", *x), "stack request action type")
	}
	w.Uint8(&id)
	(*x).Marshal(w)
}
