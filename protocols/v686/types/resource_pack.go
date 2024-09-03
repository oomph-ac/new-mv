package types

import "github.com/sandertv/gophertunnel/minecraft/protocol"

type BehaviourPackInfo struct {
	protocol.BehaviourPackInfo
}

// Marshal encodes/decodes a BehaviourPackInfo.
func (x *BehaviourPackInfo) Marshal(r protocol.IO) {
	r.String(&x.UUID)
	r.String(&x.Version)
	r.Uint64(&x.Size)
	r.String(&x.ContentKey)
	r.String(&x.SubPackName)
	r.String(&x.ContentIdentity)
	r.Bool(&x.HasScripts)
}

type TexturePackInfo struct {
	protocol.TexturePackInfo
}

// Marshal encodes/decodes a TexturePackInfo.
func (x *TexturePackInfo) Marshal(r protocol.IO) {
	r.String(&x.UUID)
	r.String(&x.Version)
	r.Uint64(&x.Size)
	r.String(&x.ContentKey)
	r.String(&x.SubPackName)
	r.String(&x.ContentIdentity)
	r.Bool(&x.HasScripts)
	r.Bool(&x.RTXEnabled)
}
