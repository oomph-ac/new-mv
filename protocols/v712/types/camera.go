package types

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

// CameraPreset represents a basic preset that can be extended upon by more complex instructions.
type CameraPreset struct {
	// Name is the name of the preset. Each preset must have their own unique name.
	Name string
	// Parent is the name of the preset that this preset extends upon. This can be left empty.
	Parent string
	// PosX is the default X position of the camera.
	PosX protocol.Optional[float32]
	// PosY is the default Y position of the camera.
	PosY protocol.Optional[float32]
	// PosZ is the default Z position of the camera.
	PosZ protocol.Optional[float32]
	// RotX is the default pitch of the camera.
	RotX protocol.Optional[float32]
	// RotY is the default yaw of the camera.
	RotY protocol.Optional[float32]
	// ViewOffset is only used in a follow_orbit camera and controls an offset based on a pivot point to the
	// player, causing it to be shifted in a certain direction.
	ViewOffset protocol.Optional[mgl32.Vec2]
	// Radius is only used in a follow_orbit camera and controls how far away from the player the camera should
	// be rendered.
	Radius protocol.Optional[float32]
	// AudioListener defines where the audio should be played from when using this preset. This is one of the
	// constants above.
	AudioListener protocol.Optional[byte]
	// PlayerEffects is currently unknown.
	PlayerEffects protocol.Optional[bool]
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
	protocol.OptionalFunc(r, &x.ViewOffset, r.Vec2)
	protocol.OptionalFunc(r, &x.Radius, r.Float32)
	protocol.OptionalFunc(r, &x.AudioListener, r.Uint8)
	protocol.OptionalFunc(r, &x.PlayerEffects, r.Bool)
}
