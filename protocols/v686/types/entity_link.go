package types

import "github.com/sandertv/gophertunnel/minecraft/protocol"

// EntityLink is a link between two entities, typically being one entity riding another.
type EntityLink struct {
	protocol.EntityLink
}

// Marshal encodes/decodes a single entity link.
func (x *EntityLink) Marshal(r protocol.IO) {
	r.Varint64(&x.RiddenEntityUniqueID)
	r.Varint64(&x.RiderEntityUniqueID)
	r.Uint8(&x.Type)
	r.Bool(&x.Immediate)
	r.Bool(&x.RiderInitiated)
}
