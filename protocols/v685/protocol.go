package v685

import (
	_ "embed"

	"github.com/oomph-ac/new-mv/internal/chunk"
	"github.com/oomph-ac/new-mv/mapping"
	"github.com/oomph-ac/new-mv/protocols/latest"
	v686 "github.com/oomph-ac/new-mv/protocols/v686"
	"github.com/oomph-ac/new-mv/translator"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

const (
	// ItemVersion is the version of items of the game which use for downgrading and upgrading.
	ItemVersion = 201
	// BlockVersion is the version of blocks (states) of the game. This version is composed
	// of 4 bytes indicating a version, interpreted as a big endian int. The current version represents
	// 1.21.0.3
	BlockVersion int32 = (1 << 24) | (21 << 16) | (0 << 8) | 3
)

var (
	//go:embed required_item_list.json
	requiredItemList []byte
	//go:embed item_runtime_ids.nbt
	itemRuntimeIDData []byte
	//go:embed block_states.nbt
	blockStateData []byte
)

type Protocol struct {
	itemMapping     mapping.Item
	blockMapping    mapping.Block
	itemTranslator  translator.ItemTranslator
	blockTranslator translator.BlockTranslator
}

func New(direct bool) *Protocol {
	itemMapping := mapping.NewItemMapping(itemRuntimeIDData, requiredItemList, ItemVersion, false)
	blockMapping := mapping.NewBlockMapping(blockStateData)
	latestBlockMapping := latest.NewBlockMapping()
	return &Protocol{
		itemMapping:     itemMapping,
		blockMapping:    blockMapping,
		itemTranslator:  translator.NewItemTranslator(itemMapping, latest.NewItemMapping(false), blockMapping, latestBlockMapping),
		blockTranslator: translator.NewBlockTranslator(blockMapping, latestBlockMapping, chunk.NewNetworkPersistentEncoding(blockMapping, BlockVersion), chunk.NewBlockPaletteEncoding(blockMapping, BlockVersion), false),
	}
}

func (Protocol) ID() int32 {
	return 685
}

func (Protocol) Ver() string {
	return "1.21.0"
}

func (Protocol) Packets(client bool) (pool packet.Pool) {
	pool = (v686.Protocol{}).Packets(client)
	delete(pool, packet.IDClientBoundCloseForm)
	return pool
}

func (Protocol) Encryption(key [32]byte) packet.Encryption {
	return packet.NewCTREncryption(key[:])
}

func (Protocol) NewReader(r minecraft.ByteReader, shieldID int32, enableLimits bool) protocol.IO {
	return NewReader(protocol.NewReader(r, shieldID, enableLimits))
}

func (Protocol) NewWriter(w minecraft.ByteWriter, shieldID int32) protocol.IO {
	return NewWriter(protocol.NewWriter(w, shieldID))
}

func (p Protocol) ConvertToLatest(pk packet.Packet, conn *minecraft.Conn) []packet.Packet {
	return p.blockTranslator.UpgradeBlockPackets(
		p.itemTranslator.UpgradeItemPackets(ProtoUpgrade([]packet.Packet{pk}), conn),
		conn,
	)
}

func ProtoUpgrade(pks []packet.Packet) []packet.Packet {
	return v686.ProtoUpgrade(pks)
}

func (p Protocol) ConvertFromLatest(pk packet.Packet, conn *minecraft.Conn) []packet.Packet {
	return ProtoDowngrade(
		p.blockTranslator.DowngradeBlockPackets(
			p.itemTranslator.DowngradeItemPackets([]packet.Packet{pk}, conn),
			conn,
		),
	)
}

func ProtoDowngrade(pks []packet.Packet) []packet.Packet {
	return v686.ProtoDowngrade(pks)
}
