package v686

import "github.com/sandertv/gophertunnel/minecraft/protocol/packet"

type Protocol struct {
}

func (Protocol) ID() int32 {
	return 686
}

func (Protocol) Ver() string {
	return "1.20.2"
}

func (Protocol) Packets(server bool) packet.Pool {
	if server {
		return packet.NewServerPool()
	}
	return packet.NewClientPool()
}
