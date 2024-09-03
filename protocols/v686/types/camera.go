package types

import "github.com/sandertv/gophertunnel/minecraft/protocol"

type CameraPreset struct {
	protocol.CameraPreset
}

// Marshal encodes/decodes a CameraPreset.
func (x *CameraPreset) Marshal(r protocol.IO) {
	r.String(&x.Name)
	r.String(&x.Parent)
	protocol.OptionalFunc(r, &x.PosX, r.Float32)
	protocol.OptionalFunc(r, &x.PosY, r.Float32)
	protocol.OptionalFunc(r, &x.PosZ, r.Float32)
	protocol.OptionalFunc(r, &x.RotX, r.Float32)
	protocol.OptionalFunc(r, &x.RotY, r.Float32)
	protocol.OptionalFunc(r, &x.AudioListener, r.Uint8)
	protocol.OptionalFunc(r, &x.PlayerEffects, r.Bool)
}
