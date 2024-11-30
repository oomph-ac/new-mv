package types

import "github.com/sandertv/gophertunnel/minecraft/protocol"

// FullContainerName contains information required to identify a container in a StackRequestSlotInfo.
type FullContainerName struct {
	// ContainerID is the ID of the container that the slot was in.
	ContainerID byte
	// DynamicContainerID is the ID of the container if it is dynamic. If the container is not dynamic, this
	// field should be left empty. A non-optional value of 0 is assumed to be non-empty.
	DynamicContainerID uint32
}

func (x *FullContainerName) Marshal(r protocol.IO) {
	r.Uint8(&x.ContainerID)
	r.Uint32(&x.DynamicContainerID)
}
