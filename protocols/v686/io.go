package v686

import (
	"fmt"

	"github.com/oomph-ac/new-mv/protocols/v686/types"
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

func (r *Reader) StackRequestAction(x *protocol.StackRequestAction) {
	var id uint8
	r.Uint8(&id)
	if !types.LookupStackRequestAction(id, x) {
		r.UnknownEnumOption(id, "stack request action type")
		return
	}

	(*x).Marshal(r)
}

func (r *Reader) PlayerInventoryAction(x *protocol.UseItemTransactionData) {
	r.Varint32(&x.LegacyRequestID)
	if x.LegacyRequestID < -1 && (x.LegacyRequestID&1) == 0 {
		protocol.Slice(r, &x.LegacySetItemSlots)
	}
	protocol.Slice(r, &x.Actions)
	r.Varuint32(&x.ActionType)
	r.BlockPos(&x.BlockPosition)
	r.Varint32(&x.BlockFace)
	r.Varint32(&x.HotBarSlot)
	r.ItemInstance(&x.HeldItem)
	r.Vec3(&x.Position)
	r.Vec3(&x.ClickedPosition)
	r.Varuint32(&x.BlockRuntimeID)
}

func (r *Reader) TransactionDataType(x *protocol.InventoryTransactionData) {
	var transactionType uint32
	r.Varuint32(&transactionType)
	switch transactionType {
	case protocol.InventoryTransactionTypeNormal: // 0
		*x = &protocol.NormalTransactionData{}
	case protocol.InventoryTransactionTypeMismatch: // 1
		*x = &protocol.MismatchTransactionData{}
	case protocol.InventoryTransactionTypeUseItem: // 2
		*x = &types.UseItemTransactionData{}
	case protocol.InventoryTransactionTypeUseItemOnEntity: // 3
		*x = &protocol.UseItemOnEntityTransactionData{}
	case protocol.InventoryTransactionTypeReleaseItem: // 4
		*x = &protocol.ReleaseItemTransactionData{}
	default:
		r.UnknownEnumOption(transactionType, "inventory transaction data type read")
	}
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

func (w *Writer) StackRequestAction(x *protocol.StackRequestAction) {
	var id uint8
	if !types.LookupStackRequestActionType(*x, &id) {
		w.UnknownEnumOption(id, "stack request action type")
		return
	}

	w.Uint8(&id)
	(*x).Marshal(w)
}

// PlayerInventoryAction writes a PlayerInventoryAction.
func (w *Writer) PlayerInventoryAction(x *protocol.UseItemTransactionData) {
	w.Varint32(&x.LegacyRequestID)
	if x.LegacyRequestID < -1 && (x.LegacyRequestID&1) == 0 {
		protocol.Slice(w, &x.LegacySetItemSlots)
	}
	protocol.Slice(w, &x.Actions)
	w.Varuint32(&x.ActionType)
	w.BlockPos(&x.BlockPosition)
	w.Varint32(&x.BlockFace)
	w.Varint32(&x.HotBarSlot)
	w.ItemInstance(&x.HeldItem)
	w.Vec3(&x.Position)
	w.Vec3(&x.ClickedPosition)
	w.Varuint32(&x.BlockRuntimeID)
}

func (w *Writer) TransactionDataType(x *protocol.InventoryTransactionData) {
	var id uint32
	if !lookupTransactionDataType(*x, &id) {
		w.UnknownEnumOption(fmt.Sprintf("%T", x), "inventory transaction data type write")
	}
	w.Varuint32(&id)
}

func lookupTransactionDataType(x protocol.InventoryTransactionData, id *uint32) bool {
	switch x.(type) {
	case *protocol.NormalTransactionData:
		*id = protocol.InventoryTransactionTypeNormal // 0
	case *protocol.MismatchTransactionData:
		*id = protocol.InventoryTransactionTypeMismatch // 1
	case *protocol.UseItemTransactionData, *types.UseItemTransactionData:
		*id = protocol.InventoryTransactionTypeUseItem // 2
	case *protocol.UseItemOnEntityTransactionData:
		*id = protocol.InventoryTransactionTypeUseItemOnEntity // 3
	case *protocol.ReleaseItemTransactionData:
		*id = protocol.InventoryTransactionTypeReleaseItem // 4
	default:
		return false
	}
	return true
}
