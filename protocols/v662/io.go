package v662

import (
	"fmt"

	"github.com/oomph-ac/new-mv/protocols/v662/types"
	v671 "github.com/oomph-ac/new-mv/protocols/v671"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

type Reader struct {
	*v671.Reader
}

func NewReader(r *protocol.Reader) *Reader {
	return &Reader{v671.NewReader(r)}
}

// Recipe reads a Recipe from the reader.
func (r *Reader) Recipe(x *protocol.Recipe) {
	var recipeType int32
	r.Varint32(&recipeType)
	if !types.LookupRecipe(recipeType, x) {
		r.UnknownEnumOption(recipeType, "crafting data recipe type")
		return
	}
	(*x).Unmarshal(r)
}

type Writer struct {
	*v671.Writer
}

func NewWriter(w *protocol.Writer) *Writer {
	return &Writer{v671.NewWriter(w)}
}

// Recipe writes a Recipe to the writer.
func (w *Writer) Recipe(x *protocol.Recipe) {
	var recipeType int32
	if !types.LookupRecipeType(*x, &recipeType) {
		w.UnknownEnumOption(fmt.Sprintf("%T", *x), "crafting data recipe type")
	}
	w.Varint32(&recipeType)
	(*x).Marshal(w)
}
