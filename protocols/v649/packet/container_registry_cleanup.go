package packet

import (
	"github.com/oomph-ac/new-mv/protocols/v649/types"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// ContainerRegistryCleanup is sent by the server to trigger a client-side cleanup of the dynamic container
// registry.
type ContainerRegistryCleanup struct {
	// RemovedContainers is a list of protocol.FullContainerName's that should be removed from the client-side
	// container registry.
	RemovedContainers []types.FullContainerName
}

// ID ...
func (*ContainerRegistryCleanup) ID() uint32 {
	return packet.IDContainerRegistryCleanup
}

func (pk *ContainerRegistryCleanup) Marshal(io protocol.IO) {
	protocol.Slice(io, &pk.RemovedContainers)
}
