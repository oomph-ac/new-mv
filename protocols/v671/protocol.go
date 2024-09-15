package v671

import (
	_ "embed"
	"fmt"

	"github.com/oomph-ac/new-mv/internal/chunk"
	"github.com/oomph-ac/new-mv/mapping"
	"github.com/oomph-ac/new-mv/protocols/latest"
	v671packet "github.com/oomph-ac/new-mv/protocols/v671/packet"
	"github.com/oomph-ac/new-mv/protocols/v671/types"
	v685 "github.com/oomph-ac/new-mv/protocols/v685"
	"github.com/oomph-ac/new-mv/translator"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

const (
	// ItemVersion is the version of items of the game which use for downgrading and upgrading.
	ItemVersion = 191
	// BlockVersion is the version of blocks (states) of the game. This version is composed
	// of 4 bytes indicating a version, interpreted as a big endian int. The current version represents
	// 1.20.80.5
	BlockVersion int32 = (1 << 24) | (20 << 16) | (80 << 8) | 5
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
	return 671
}

func (Protocol) Ver() string {
	return "1.20.80"
}

func (Protocol) Packets(client bool) (pool packet.Pool) {
	pool = (v685.Protocol{}).Packets(client)
	if !client {
		pool[packet.IDContainerClose] = func() packet.Packet { return &v671packet.StartGame{} }
	}

	pool[packet.IDContainerClose] = func() packet.Packet { return &v671packet.ContainerClose{} }
	pool[packet.IDText] = func() packet.Packet { return &v671packet.Text{} }
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
	//fmt.Printf("client -> server %T\n", pk)
	return p.blockTranslator.UpgradeBlockPackets(
		p.itemTranslator.UpgradeItemPackets(ProtoUpgrade([]packet.Packet{pk}), conn),
		conn,
	)
}

func ProtoUpgrade(pks []packet.Packet) []packet.Packet {
	for index, pk := range pks {
		switch pk := pk.(type) {
		case *v671packet.ContainerClose:
			pks[index] = &packet.ContainerClose{
				WindowID:   pk.WindowID,
				ServerSide: pk.ServerSide,
			}
		case *v671packet.Text:
			pks[index] = &packet.Text{
				TextType:         pk.TextType,
				NeedsTranslation: pk.NeedsTranslation,
				SourceName:       pk.SourceName,
				Message:          pk.Message,
				Parameters:       pk.Parameters,
				XUID:             pk.XUID,
				PlatformChatID:   pk.PlatformChatID,
			}
		}
	}

	return v685.ProtoUpgrade(pks)
}

func (p Protocol) ConvertFromLatest(pk packet.Packet, conn *minecraft.Conn) []packet.Packet {
	fmt.Printf("server -> client %T\n", pk)
	return ProtoDowngrade(
		p.blockTranslator.DowngradeBlockPackets(
			p.itemTranslator.DowngradeItemPackets([]packet.Packet{pk}, conn),
			conn,
		),
	)
}

func ProtoDowngrade(pks []packet.Packet) []packet.Packet {
	pks = v685.ProtoDowngrade(pks)
	for index, pk := range pks {
		switch pk := pk.(type) {
		case *packet.CraftingData:
			for index, recp := range pk.Recipes {
				switch recp := recp.(type) {
				case *protocol.ShapelessRecipe:
					pk.Recipes[index] = &types.ShapelessRecipe{ShapelessRecipe: *recp}
				case *protocol.ShapedRecipe:
					pk.Recipes[index] = &types.ShapedRecipe{ShapedRecipe: *recp}
				case *protocol.ShulkerBoxRecipe:
					pk.Recipes[index] = &types.ShulkerBoxRecipe{ShapelessRecipe: types.ShapelessRecipe{ShapelessRecipe: recp.ShapelessRecipe}}
				case *protocol.ShapelessChemistryRecipe:
					pk.Recipes[index] = &types.ShulkerBoxRecipe{ShapelessRecipe: types.ShapelessRecipe{ShapelessRecipe: recp.ShapelessRecipe}}
				case *protocol.ShapedChemistryRecipe:
					pk.Recipes[index] = &types.ShapedChemistryRecipe{ShapedRecipe: types.ShapedRecipe{ShapedRecipe: recp.ShapedRecipe}}
				}
			}

			pks[index] = pk
		case *packet.ContainerClose:
			pks[index] = &v671packet.ContainerClose{
				WindowID:   pk.WindowID,
				ServerSide: pk.ServerSide,
			}
		case *packet.Text:
			pks[index] = &v671packet.Text{
				TextType:         pk.TextType,
				NeedsTranslation: pk.NeedsTranslation,
				SourceName:       pk.SourceName,
				Message:          pk.Message,
				Parameters:       pk.Parameters,
				XUID:             pk.XUID,
				PlatformChatID:   pk.PlatformChatID,
			}
		case *packet.StartGame:
			pks[index] = &v671packet.StartGame{StartGame: pk}
		}
	}

	return pks
}
