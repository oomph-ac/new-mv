package v662

import (
	_ "embed"

	"github.com/oomph-ac/new-mv/internal/chunk"
	"github.com/oomph-ac/new-mv/mapping"
	"github.com/oomph-ac/new-mv/protocols/latest"
	v662packet "github.com/oomph-ac/new-mv/protocols/v662/packet"
	"github.com/oomph-ac/new-mv/protocols/v662/types"
	v671 "github.com/oomph-ac/new-mv/protocols/v671"
	v671packet "github.com/oomph-ac/new-mv/protocols/v671/packet"
	v671types "github.com/oomph-ac/new-mv/protocols/v671/types"
	v686packet "github.com/oomph-ac/new-mv/protocols/v686/packet"
	"github.com/oomph-ac/new-mv/translator"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

const (
	// ItemVersion is the version of items of the game which use for downgrading and upgrading.
	ItemVersion = 181
	// BlockVersion is the version of blocks (states) of the game. This version is composed
	// of 4 bytes indicating a version, interpreted as a big endian int. The current version represents
	// 1.20.80.5
	BlockVersion int32 = (1 << 24) | (20 << 16) | (71 << 8) | 1
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
	return 662
}

func (Protocol) Ver() string {
	return "1.20.71"
}

func (Protocol) Packets(client bool) (pool packet.Pool) {
	pool = (v671.Protocol{}).Packets(client)
	if !client {
		pool[packet.IDCorrectPlayerMovePrediction] = func() packet.Packet { return &v662packet.CorrectPlayerMovePrediction{} }
		pool[packet.IDResourcePackStack] = func() packet.Packet { return &v662packet.ResourcePackStack{} }
		pool[packet.IDStartGame] = func() packet.Packet { return &v662packet.StartGame{} }
		pool[packet.IDUpdateBlockSynced] = func() packet.Packet { return &v662packet.UpdateBlockSynced{} }
		pool[packet.IDUpdatePlayerGameType] = func() packet.Packet { return &v662packet.UpdatePlayerGameType{} }
	}

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
	return v671.ProtoUpgrade(pks)
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
	pks = v671.ProtoDowngrade(pks)
	for index, pk := range pks {
		switch pk := pk.(type) {
		case *v686packet.CorrectPlayerMovePrediction:
			pks[index] = &v662packet.CorrectPlayerMovePrediction{
				PredictionType: pk.PredictionType,
				Position:       pk.Position,
				Delta:          pk.Delta,
				OnGround:       pk.OnGround,
				Tick:           pk.Tick,
			}
		case *packet.CraftingData:
			for index, recp := range pk.Recipes {
				switch recp := recp.(type) {
				case *v671types.ShapedRecipe:
					pk.Recipes[index] = &types.ShapedRecipe{ShapedRecipe: *recp}
				case *v671types.ShapedChemistryRecipe:
					pk.Recipes[index] = &types.ShapedRecipe{ShapedRecipe: recp.ShapedRecipe}
				}
			}
		case *packet.ResourcePackStack:
			pks[index] = &v662packet.ResourcePackStack{
				TexturePackRequired:          pk.TexturePackRequired,
				BehaviourPacks:               pk.BehaviourPacks,
				TexturePacks:                 pk.TexturePacks,
				BaseGameVersion:              pk.BaseGameVersion,
				Experiments:                  pk.Experiments,
				ExperimentsPreviouslyToggled: pk.ExperimentsPreviouslyToggled,
			}
		case *v671packet.StartGame:
			pks[index] = &v662packet.StartGame{StartGame: pk}
		case *packet.UpdateBlockSynced:
			pks[index] = &v662packet.UpdateBlockSynced{
				Position:          pk.Position,
				NewBlockRuntimeID: pk.NewBlockRuntimeID,
				Flags:             pk.Flags,
				Layer:             pk.Layer,
				EntityUniqueID:    int64(pk.EntityUniqueID),
				TransitionType:    pk.TransitionType,
			}
		case *packet.UpdatePlayerGameType:
			pks[index] = &v662packet.UpdatePlayerGameType{
				GameType:       pk.GameType,
				PlayerUniqueID: pk.PlayerUniqueID,
			}
		}
	}

	return pks
}
