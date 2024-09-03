package types

import "github.com/sandertv/gophertunnel/minecraft/protocol"

type UseItemTransactionData struct {
	protocol.UseItemTransactionData
}

func (data *UseItemTransactionData) Marshal(r protocol.IO) {
	r.Varuint32(&data.ActionType)
	r.UBlockPos(&data.BlockPosition)
	r.Varint32(&data.BlockFace)
	r.Varint32(&data.HotBarSlot)
	r.ItemInstance(&data.HeldItem)
	r.Vec3(&data.Position)
	r.Vec3(&data.ClickedPosition)
	r.Varuint32(&data.BlockRuntimeID)
}
