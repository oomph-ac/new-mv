package v671

import (
	"fmt"

	"github.com/oomph-ac/new-mv/protocols/v671/types"
	v685 "github.com/oomph-ac/new-mv/protocols/v685"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

type Reader struct {
	*v685.Reader
}

func NewReader(r *protocol.Reader) *Reader {
	return &Reader{v685.NewReader(r)}
}

// Recipe reads a Recipe from the reader.
func (r *Reader) Recipe(x *protocol.Recipe) {
	var recipeType int32
	r.Varint32(&recipeType)
	if !types.LookupRecipe(recipeType, x) {
		r.UnknownEnumOption(recipeType, "crafting data recipe type")
		return
	}
	fmt.Printf("MV recp read %T\n", *x)
	(*x).Unmarshal(r)
}

type Writer struct {
	*v685.Writer
}

func NewWriter(w *protocol.Writer) *Writer {
	return &Writer{v685.NewWriter(w)}
}

// Recipe writes a Recipe to the writer.
func (w *Writer) Recipe(x *protocol.Recipe) {
	var recipeType int32
	if !types.LookupRecipeType(*x, &recipeType) {
		w.UnknownEnumOption(fmt.Sprintf("%T", *x), "crafting data recipe type")
	}
	w.Varint32(&recipeType)
	fmt.Printf("MV recp write %T\n", *x)
	(*x).Marshal(w)
}
