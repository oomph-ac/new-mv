package raknet

import (
	"context"
	"log/slog"
	"net"

	"github.com/sandertv/go-raknet"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// MultiRakNet is an implementation of a RakNet v9/10 Network.
type MultiRakNet struct {
	l *slog.Logger
}

// legacyRakNet represents the legacy version of RakNet, necessary for versions higher or equal to v1.16.0.
const legacyRakNet = 10

// DialContext ...
func (r MultiRakNet) DialContext(ctx context.Context, address string) (net.Conn, error) {
	return raknet.Dialer{ErrorLog: r.l.With("net origin", "raknet")}.DialContext(ctx, address)
}

// PingContext ...
func (r MultiRakNet) PingContext(ctx context.Context, address string) (response []byte, err error) {
	return raknet.Dialer{ErrorLog: r.l.With("net origin", "raknet")}.PingContext(ctx, address)
}

// Listen ...
func (r MultiRakNet) Listen(address string) (minecraft.NetworkListener, error) {
	return raknet.ListenConfig{
		ErrorLog:         r.l.With("net origin", "raknet"),
		ProtocolVersions: []byte{legacyRakNet}, // Version 10 is required for legacy versions.
	}.Listen(address)
}

// Compression ...
func (MultiRakNet) Compression(net.Conn) packet.Compression {
	return packet.FlateCompression
}

// init registers the MultiRakNet network. It overrides the existing minecraft.RakNet network.
func init() {
	minecraft.RegisterNetwork("raknet", func(l *slog.Logger) minecraft.Network {
		return MultiRakNet{l: l}
	})
}
