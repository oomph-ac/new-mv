package v712

import (
	_ "embed"

	"github.com/oomph-ac/new-mv/internal/chunk"
	"github.com/oomph-ac/new-mv/mapping"
	"github.com/oomph-ac/new-mv/protocols/latest"
	v712packet "github.com/oomph-ac/new-mv/protocols/v712/packet"
	"github.com/oomph-ac/new-mv/protocols/v712/types"
	"github.com/oomph-ac/new-mv/translator"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

const (
	// ItemVersion is the version of items of the game which use for downgrading and upgrading.
	ItemVersion = 211
	// BlockVersion is the version of blocks (states) of the game. This version is composed
	// of 4 bytes indicating a version, interpreted as a big endian int. The current version represents
	// 1.21.20
	BlockVersion int32 = (1 << 24) | (21 << 16) | (20 << 8)
)

var (
	//go:embed required_item_list.json
	requiredItemList []byte
	//go:embed item_runtime_ids.nbt
	itemRuntimeIDData []byte
	//go:embed block_states.nbt
	blockStateData []byte

	packetPool_server packet.Pool
	packetPool_client packet.Pool
)

func init() {
	packetPool_server = packet.NewServerPool()
	packetPool_client = packet.NewClientPool()

	// ------------------------ 1.21.30 changes ------------------------
	delete(packetPool_server, packet.IDMovementEffect)
	delete(packetPool_server, packet.IDSetMovementAuthority)

	packetPool_server[packet.IDMobEffect] = func() packet.Packet { return &v712packet.MobEffect{} }
	packetPool_client[packet.IDPlayerAuthInput] = func() packet.Packet { return &v712packet.PlayerAuthInput{} }
	// ------------------------ 1.21.30 changes ------------------------

	// ------------------------ 1.21.20 changes ------------------------
	delete(packetPool_server, packet.IDCameraAimAssist)
	delete(packetPool_server, packet.IDContainerRegistryCleanup)

	packetPool_server[packet.IDEmote] = func() packet.Packet { return &v712packet.Emote{} }
	packetPool_client[packet.IDEmote] = func() packet.Packet { return &v712packet.Emote{} }

	packetPool_server[packet.IDCameraPresets] = func() packet.Packet { return &v712packet.CameraPresets{} }
	packetPool_server[packet.IDContainerRegistryCleanup] = func() packet.Packet { return &v712packet.ContainerRegistryCleanup{} }
	packetPool_server[packet.IDInventoryContent] = func() packet.Packet { return &v712packet.InventoryContent{} }
	packetPool_server[packet.IDInventorySlot] = func() packet.Packet { return &v712packet.InventorySlot{} }
	packetPool_server[packet.IDItemStackResponse] = func() packet.Packet { return &v712packet.ItemStackResponse{} }
	packetPool_server[packet.IDResourcePacksInfo] = func() packet.Packet { return &v712packet.ResourcePacksInfo{} }
	packetPool_server[packet.IDTransfer] = func() packet.Packet { return &v712packet.Transfer{} }
	packetPool_server[packet.IDUpdateAttributes] = func() packet.Packet { return &v712packet.UpdateAttributes{} }
	// ------------------------ 1.21.20 changes ------------------------
}

type Protocol struct {
	minecraft.Protocol
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
	return 712
}

func (Protocol) Ver() string {
	return "1.21.30"
}

func (Protocol) Packets(listener bool) packet.Pool {
	if listener {
		return packetPool_client
	}
	return packetPool_server
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
	return ProtoUpgrade(p.blockTranslator.UpgradeBlockPackets(
		p.itemTranslator.UpgradeItemPackets([]packet.Packet{pk}, conn),
		conn,
	))
}

func ProtoUpgrade(pks []packet.Packet) []packet.Packet {
	for index, pk := range pks {
		switch pk := pk.(type) {
		case *v712packet.Emote:
			pks[index] = &packet.Emote{
				EntityRuntimeID: pk.EntityRuntimeID,
				EmoteID:         pk.EmoteID,
				EmoteLength:     100, // TODO: ???
				XUID:            pk.XUID,
				PlatformID:      pk.PlatformID,
				Flags:           pk.Flags,
			}
		case *v712packet.PlayerAuthInput:
			pks[index] = &packet.PlayerAuthInput{
				Pitch:                  pk.Pitch,
				Yaw:                    pk.Yaw,
				Position:               pk.Position,
				MoveVector:             pk.MoveVector,
				HeadYaw:                pk.HeadYaw,
				InputData:              pk.InputData,
				InputMode:              pk.InputMode,
				PlayMode:               pk.PlayMode,
				InteractionModel:       pk.InteractionModel,
				InteractPitch:          pk.GazeDirection.X(),
				InteractYaw:            pk.GazeDirection.Y(),
				Tick:                   pk.Tick,
				Delta:                  pk.Delta,
				ItemInteractionData:    pk.ItemInteractionData,
				ItemStackRequest:       pk.ItemStackRequest,
				BlockActions:           pk.BlockActions,
				VehicleRotation:        pk.VehicleRotation,
				ClientPredictedVehicle: pk.ClientPredictedVehicle,
				AnalogueMoveVector:     pk.AnalogueMoveVector,
			}
		case *packet.ItemStackRequest:
			for index, request := range pk.Requests {
				pk.Requests[index] = protocol.ItemStackRequest{
					RequestID:     request.RequestID,
					Actions:       types.UpgradeItemStackActions(request.Actions),
					FilterStrings: request.FilterStrings,
					FilterCause:   request.FilterCause,
				}
			}
		}
	}

	return pks
}

func (p Protocol) ConvertFromLatest(pk packet.Packet, conn *minecraft.Conn) []packet.Packet {
	return ProtoDowngrade(p.blockTranslator.DowngradeBlockPackets(
		p.itemTranslator.DowngradeItemPackets([]packet.Packet{pk}, conn),
		conn,
	))
}

func ProtoDowngrade(pks []packet.Packet) []packet.Packet {
	for index, pk := range pks {
		switch pk := pk.(type) {
		case *packet.CameraPresets:
			presets := make([]types.CameraPreset, len(pk.Presets))
			for presetIndex, preset := range pk.Presets {
				presets[presetIndex] = types.CameraPreset{
					Name:          preset.Name,
					Parent:        preset.Parent,
					PosX:          preset.PosX,
					PosY:          preset.PosY,
					PosZ:          preset.PosZ,
					RotX:          preset.RotX,
					RotY:          preset.RotY,
					ViewOffset:    preset.ViewOffset,
					Radius:        preset.Radius,
					AudioListener: preset.AudioListener,
					PlayerEffects: preset.PlayerEffects,
				}
			}

			pks[index] = &v712packet.CameraPresets{
				Presets: presets,
			}
		case *packet.ContainerRegistryCleanup:
			containers := make([]types.FullContainerName, len(pk.RemovedContainers))
			for index, c := range pk.RemovedContainers {
				dynContainerID, _ := c.DynamicContainerID.Value()
				containers[index] = types.FullContainerName{
					ContainerID:        c.ContainerID,
					DynamicContainerID: dynContainerID,
				}
			}
		case *packet.Emote:
			pks[index] = &v712packet.Emote{
				EntityRuntimeID: pk.EntityRuntimeID,
				EmoteID:         pk.EmoteID,
				XUID:            pk.XUID,
				PlatformID:      pk.PlatformID,
				Flags:           pk.Flags,
			}
		case *packet.InventoryContent:
			dynID, _ := pk.Container.DynamicContainerID.Value()
			pks[index] = &v712packet.InventoryContent{
				WindowID:        pk.WindowID,
				Content:         pk.Content,
				DynamicWindowID: dynID,
			}
		case *packet.InventorySlot:
			dynID, _ := pk.Container.DynamicContainerID.Value()
			pks[index] = &v712packet.InventorySlot{
				WindowID:        pk.WindowID,
				Slot:            pk.Slot,
				DynamicWindowID: dynID,
				NewItem:         pk.NewItem,
			}
		case *packet.ItemStackResponse:
			responses := make([]types.ItemStackResponse, len(pk.Responses))
			for index, response := range pk.Responses {
				containerInfo := make([]types.StackResponseContainerInfo, len(response.ContainerInfo))
				for cIndex, info := range response.ContainerInfo {
					containerInfo[cIndex] = types.StackResponseContainerInfo{
						Container: types.DowngradeContainer(info.Container),
						SlotInfo:  info.SlotInfo,
					}
				}

				responses[index] = types.ItemStackResponse{
					Status:        response.Status,
					RequestID:     response.RequestID,
					ContainerInfo: containerInfo,
				}
			}

			pks[index] = &v712packet.ItemStackResponse{
				Responses: responses,
			}
		case *packet.MobEffect:
			pks[index] = &v712packet.MobEffect{
				EntityRuntimeID: pk.EntityRuntimeID,
				Operation:       pk.Operation,
				EffectType:      pk.EffectType,
				Amplifier:       pk.Amplifier,
				Particles:       pk.Particles,
				Duration:        pk.Duration,
				Tick:            pk.Tick,
			}
		case *packet.ResourcePacksInfo:
			packs := make([]types.TexturePackInfo, len(pk.TexturePacks))
			packURLs := []protocol.PackURL{}
			for packIndex, pack := range pk.TexturePacks {
				packs[packIndex] = types.TexturePackInfo{
					UUID:            pack.UUID,
					Version:         pack.Version,
					Size:            pack.Size,
					ContentKey:      pack.ContentKey,
					SubPackName:     pack.SubPackName,
					ContentIdentity: pack.ContentIdentity,
					HasScripts:      pack.HasScripts,
					AddonPack:       pack.AddonPack,
					RTXEnabled:      pack.RTXEnabled,
				}

				if pack.DownloadURL != "" {
					packURLs = append(packURLs, protocol.PackURL{
						UUIDVersion: pack.UUID + "_" + pack.Version,
						URL:         pack.DownloadURL,
					})
				}
			}

			pks[index] = &v712packet.ResourcePacksInfo{
				TexturePackRequired: pk.TexturePackRequired,
				HasAddons:           pk.HasAddons,
				HasScripts:          pk.HasScripts,
				TexturePacks:        packs,
				PackURLs:            packURLs,
			}
		case *packet.Transfer:
			pks[index] = &v712packet.Transfer{
				Address: pk.Address,
				Port:    pk.Port,
			}
		case *packet.UpdateAttributes:
			attributes := make([]types.Attribute, len(pk.Attributes))
			for index, a := range pk.Attributes {
				attributes[index] = types.Attribute{
					AttributeValue: a.AttributeValue,
					Default:        a.Default,
					Modifiers:      a.Modifiers,
				}
			}

			pks[index] = &v712packet.UpdateAttributes{
				EntityRuntimeID: pk.EntityRuntimeID,
				Attributes:      attributes,
				Tick:            pk.Tick,
			}
		}
	}

	return pks
}
