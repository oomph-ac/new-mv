package types

import "github.com/sandertv/gophertunnel/minecraft/protocol"

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
