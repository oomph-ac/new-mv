package v685

import (
	_ "embed"

	"github.com/oomph-ac/new-mv/internal/chunk"
	"github.com/oomph-ac/new-mv/mapping"
	"github.com/oomph-ac/new-mv/protocols/latest"
	v685packet "github.com/oomph-ac/new-mv/protocols/v686/packet"
	"github.com/oomph-ac/new-mv/protocols/v686/types"
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
	// 1.21.2.2
	BlockVersion int32 = (1 << 24) | (21 << 16) | (2 << 8)
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

	noPacketsAvailable = []packet.Packet{}
)

func init() {
	packetPool_server = packet.NewServerPool()
	packetPool_client = packet.NewClientPool()

	// ------------------------ 1.21.30 changes ------------------------
	delete(packetPool_server, packet.IDMovementEffect)
	delete(packetPool_server, packet.IDSetMovementAuthority)

	packetPool_server[packet.IDMobEffect] = func() packet.Packet { return &v685packet.MobEffect{} }
	packetPool_client[packet.IDPlayerAuthInput] = func() packet.Packet { return &v685packet.PlayerAuthInput{} }
	// ------------------------ 1.21.30 changes ------------------------

	// ------------------------ 1.21.20 changes ------------------------
	delete(packetPool_server, packet.IDCameraAimAssist)
	delete(packetPool_server, packet.IDContainerRegistryCleanup)

	packetPool_server[packet.IDEmote] = func() packet.Packet { return &v685packet.Emote{} }
	packetPool_client[packet.IDEmote] = func() packet.Packet { return &v685packet.Emote{} }

	packetPool_server[packet.IDCameraPresets] = func() packet.Packet { return &v685packet.CameraPresets{} }
	packetPool_server[packet.IDContainerRegistryCleanup] = func() packet.Packet { return &v685packet.ContainerRegistryCleanup{} }
	packetPool_server[packet.IDItemStackResponse] = func() packet.Packet { return &v685packet.ItemStackResponse{} }
	packetPool_server[packet.IDResourcePacksInfo] = func() packet.Packet { return &v685packet.ResourcePacksInfo{} }
	packetPool_server[packet.IDTransfer] = func() packet.Packet { return &v685packet.Transfer{} }
	packetPool_server[packet.IDUpdateAttributes] = func() packet.Packet { return &v685packet.UpdateAttributes{} }
	// ------------------------ 1.21.20 changes ------------------------

	// ------------------------ 1.21.2 changes ------------------------
	delete(packetPool_server, packet.IDCurrentStructureFeature)
	delete(packetPool_server, packet.IDJigsawStructureData)
	delete(packetPool_client, packet.IDServerBoundDiagnostics)
	delete(packetPool_client, packet.IDServerBoundLoadingScreen)

	packetPool_server[packet.IDMobArmourEquipment] = func() packet.Packet { return &v685packet.MobArmourEquipment{} }
	packetPool_client[packet.IDMobArmourEquipment] = func() packet.Packet { return &v685packet.MobArmourEquipment{} }

	packetPool_server[packet.IDEditorNetwork] = func() packet.Packet { return &v685packet.EditorNetwork{} }
	packetPool_client[packet.IDEditorNetwork] = func() packet.Packet { return &v685packet.EditorNetwork{} }

	packetPool_server[packet.IDAddActor] = func() packet.Packet { return &v685packet.AddActor{} }
	packetPool_server[packet.IDAddPlayer] = func() packet.Packet { return &v685packet.AddPlayer{} }
	packetPool_server[packet.IDCameraInstruction] = func() packet.Packet { return &v685packet.CameraInstruction{} }
	packetPool_server[packet.IDChangeDimension] = func() packet.Packet { return &v685packet.ChangeDimension{} }
	packetPool_server[packet.IDCompressedBiomeDefinitionList] = func() packet.Packet { return &v685packet.CompressedBiomeDefinitionList{} }
	packetPool_server[packet.IDCorrectPlayerMovePrediction] = func() packet.Packet { return &v685packet.CorrectPlayerMovePrediction{} }
	packetPool_server[packet.IDDisconnect] = func() packet.Packet { return &v685packet.Disconnect{} }
	packetPool_server[packet.IDInventoryContent] = func() packet.Packet { return &v685packet.InventoryContent{} }
	packetPool_server[packet.IDInventorySlot] = func() packet.Packet { return &v685packet.InventorySlot{} }
	packetPool_server[packet.IDPlayerArmourDamage] = func() packet.Packet { return &v685packet.PlayerArmourDamage{} }
	packetPool_server[packet.IDSetTitle] = func() packet.Packet { return &v685packet.SetTitle{} }
	packetPool_server[packet.IDStopSound] = func() packet.Packet { return &v685packet.StopSound{} }
	// ------------------------ 1.21.2 changes ------------------------
}

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
		case *v685packet.Emote:
			pks[index] = &packet.Emote{
				EntityRuntimeID: pk.EntityRuntimeID,
				EmoteID:         pk.EmoteID,
				EmoteLength:     100, // TODO: ???
				XUID:            pk.XUID,
				PlatformID:      pk.PlatformID,
				Flags:           pk.Flags,
			}
		case *v685packet.EditorNetwork:
			pks[index] = &packet.EditorNetwork{
				RouteToManager: false,
				Payload:        pk.Payload,
			}
		case *v685packet.MobArmourEquipment:
			pks[index] = &packet.MobArmourEquipment{
				EntityRuntimeID: pk.EntityRuntimeID,
				Helmet:          pk.Helmet,
				Chestplate:      pk.Chestplate,
				Leggings:        pk.Leggings,
				Boots:           pk.Boots,
			}
		case *packet.InventoryTransaction:
			var transactionData protocol.InventoryTransactionData = pk.TransactionData
			if t, ok := pk.TransactionData.(*types.UseItemTransactionData); ok {
				transactionData = &protocol.UseItemTransactionData{
					ActionType:      t.ActionType,
					BlockPosition:   t.BlockPosition,
					BlockFace:       t.BlockFace,
					HotBarSlot:      t.HotBarSlot,
					HeldItem:        t.HeldItem,
					Position:        t.Position,
					ClickedPosition: t.ClickedPosition,
					BlockRuntimeID:  t.BlockRuntimeID,
				}
			}

			pk.TransactionData = transactionData
			pks[index] = pk
		case *packet.ItemStackRequest:
			for i, req := range pk.Requests {
				pk.Requests[i] = protocol.ItemStackRequest{
					RequestID:     req.RequestID,
					Actions:       types.UpgradeItemStackActions(req.Actions),
					FilterStrings: req.FilterStrings,
					FilterCause:   req.FilterCause,
				}
			}
		case *v685packet.PlayerAuthInput:
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
		case *packet.AddActor:
			eLinks := make([]types.EntityLink, len(pk.EntityLinks))
			for index, link := range pk.EntityLinks {
				eLinks[index] = types.EntityLink{EntityLink: link}
			}

			pks[index] = &v685packet.AddActor{
				EntityUniqueID:   pk.EntityUniqueID,
				EntityRuntimeID:  pk.EntityRuntimeID,
				EntityType:       pk.EntityType,
				Position:         pk.Position,
				Velocity:         pk.Velocity,
				Pitch:            pk.Pitch,
				Yaw:              pk.Yaw,
				HeadYaw:          pk.HeadYaw,
				BodyYaw:          pk.BodyYaw,
				Attributes:       pk.Attributes,
				EntityMetadata:   pk.EntityMetadata,
				EntityProperties: pk.EntityProperties,
				EntityLinks:      eLinks,
			}
		case *packet.AddPlayer:
			eLinks := make([]types.EntityLink, len(pk.EntityLinks))
			for index, link := range pk.EntityLinks {
				eLinks[index] = types.EntityLink{EntityLink: link}
			}

			pks[index] = &v685packet.AddPlayer{
				UUID:             pk.UUID,
				Username:         pk.Username,
				EntityRuntimeID:  pk.EntityRuntimeID,
				PlatformChatID:   pk.PlatformChatID,
				Position:         pk.Position,
				Velocity:         pk.Velocity,
				Pitch:            pk.Pitch,
				Yaw:              pk.Yaw,
				HeadYaw:          pk.HeadYaw,
				HeldItem:         pk.HeldItem,
				GameType:         pk.GameType,
				EntityMetadata:   pk.EntityMetadata,
				EntityProperties: pk.EntityProperties,
				AbilityData:      pk.AbilityData,
				EntityLinks:      eLinks,
				DeviceID:         pk.DeviceID,
				BuildPlatform:    pk.BuildPlatform,
			}
		case *packet.CameraInstruction:
			pks[index] = &v685packet.CameraInstruction{
				Set:   pk.Set,
				Clear: pk.Clear,
				Fade:  pk.Fade,
			}
		case *packet.CameraPresets:
			presets := make([]types.CameraPreset, len(pk.Presets))
			for index, preset := range pk.Presets {
				presets[index] = types.CameraPreset{
					CameraPreset: preset,
				}
			}

			pks[index] = &v685packet.CameraPresets{
				Presets: presets,
			}
		case *packet.ChangeDimension:
			pks[index] = &v685packet.ChangeDimension{
				Dimension: pk.Dimension,
				Position:  pk.Position,
				Respawn:   pk.Respawn,
			}
		case *packet.ContainerRegistryCleanup:
			containers := make([]types.FullContainerName, len(pk.RemovedContainers))
			for index, container := range pk.RemovedContainers {
				containers[index] = types.DowngradeContainer(container)
			}

			pks[index] = &v685packet.ContainerRegistryCleanup{
				RemovedContainers: containers,
			}
		case *packet.CorrectPlayerMovePrediction:
			pks[index] = &v685packet.CorrectPlayerMovePrediction{
				PredictionType: pk.PredictionType,
				Position:       pk.Position,
				Delta:          pk.Delta,
				Rotation:       pk.Rotation,
				OnGround:       pk.OnGround,
				Tick:           pk.Tick,
			}
		case *packet.Disconnect:
			pks[index] = &v685packet.Disconnect{
				Reason:                  pk.Reason,
				HideDisconnectionScreen: pk.HideDisconnectionScreen,
				Message:                 pk.Message,
			}
		case *packet.EditorNetwork:
			pks[index] = &v685packet.EditorNetwork{
				Payload: pk.Payload,
			}
		case *packet.Emote:
			pks[index] = &v685packet.Emote{
				EntityRuntimeID: pk.EntityRuntimeID,
				EmoteID:         pk.EmoteID,
				XUID:            pk.XUID,
				PlatformID:      pk.PlatformID,
				Flags:           pk.Flags,
			}
		case *packet.InventoryContent:
			pks[index] = &v685packet.InventoryContent{
				WindowID: pk.WindowID,
				Content:  pk.Content,
			}
		case *packet.InventorySlot:
			pks[index] = &v685packet.InventorySlot{
				WindowID: pk.WindowID,
				Slot:     pk.Slot,
				NewItem:  pk.NewItem,
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

			pks[index] = &v685packet.ItemStackResponse{
				Responses: responses,
			}
		case *packet.MobArmourEquipment:
			pks[index] = &v685packet.MobArmourEquipment{
				EntityRuntimeID: pk.EntityRuntimeID,
				Helmet:          pk.Helmet,
				Chestplate:      pk.Chestplate,
				Leggings:        pk.Leggings,
				Boots:           pk.Boots,
			}
		case *packet.MobEffect:
			pks[index] = &v685packet.MobEffect{
				EntityRuntimeID: pk.EntityRuntimeID,
				Operation:       pk.Operation,
				EffectType:      pk.EffectType,
				Amplifier:       pk.Amplifier,
				Particles:       pk.Particles,
				Duration:        pk.Duration,
				Tick:            pk.Tick,
			}
		case *packet.PlayerArmourDamage:
			var bitset uint8
			if pk.Bitset&packet.PlayerArmourDamageFlagHelmet != 0 {
				bitset = 0b0001
			}
			if pk.Bitset&packet.PlayerArmourDamageFlagChestplate != 0 {
				bitset = bitset | 0b0010
			}
			if pk.Bitset&packet.PlayerArmourDamageFlagLeggings != 0 {
				bitset = bitset | 0b0100
			}
			if pk.Bitset&packet.PlayerArmourDamageFlagBoots != 0 {
				bitset = bitset | 0b1000
			}

			pks[index] = &v685packet.PlayerArmourDamage{
				Bitset:           bitset,
				HelmetDamage:     pk.HelmetDamage,
				ChestplateDamage: pk.ChestplateDamage,
				LeggingsDamage:   pk.LeggingsDamage,
				BootsDamage:      pk.BootsDamage,
			}
		case *packet.ResourcePacksInfo:
			tPacks := make([]types.TexturePackInfo, len(pk.TexturePacks))
			packURLs := []protocol.PackURL{}
			for index, pack := range pk.TexturePacks {
				tPacks[index] = types.TexturePackInfo{TexturePackInfo: pack}
				if pack.DownloadURL != "" {
					packURLs = append(packURLs, protocol.PackURL{
						UUIDVersion: pack.UUID + "_" + pack.Version,
						URL:         pack.DownloadURL,
					})
				}
			}

			pks[index] = &v685packet.ResourcePacksInfo{
				TexturePackRequired: pk.TexturePackRequired,
				HasAddons:           pk.HasAddons,
				HasScripts:          pk.HasScripts,
				BehaviourPacks:      []types.TexturePackInfo{},
				TexturePacks:        tPacks,
				ForcingServerPacks:  true,
				PackURLs:            packURLs,
			}
		case *packet.SetTitle:
			pks[index] = &v685packet.SetTitle{
				ActionType:       pk.ActionType,
				Text:             pk.Text,
				FadeInDuration:   pk.FadeInDuration,
				RemainDuration:   pk.RemainDuration,
				FadeOutDuration:  pk.FadeOutDuration,
				XUID:             pk.XUID,
				PlatformOnlineID: pk.PlatformOnlineID,
			}
		case *packet.StopSound:
			pks[index] = &v685packet.StopSound{
				SoundName: pk.SoundName,
				StopAll:   pk.StopAll,
			}
		case *packet.SetActorLink:
			pks[index] = &v685packet.SetActorLink{
				EntityLink: types.EntityLink{
					EntityLink: pk.EntityLink,
				},
			}
		case *packet.Transfer:
			pks[index] = &v685packet.Transfer{
				Address: pk.Address,
				Port:    pk.Port,
			}
		case *packet.UpdateAttributes:
			attributes := make([]types.Attribute, len(pk.Attributes))
			for index, a := range pk.Attributes {
				attributes[index] = types.Attribute{
					AttributeValue: protocol.AttributeValue{
						Name:  a.Name,
						Value: a.Value,
						Min:   a.Min,
						Max:   a.Max,
					},
					Default:   a.Default,
					Modifiers: a.Modifiers,
				}
			}

			pks[index] = &v685packet.UpdateAttributes{
				EntityRuntimeID: pk.EntityRuntimeID,
				Attributes:      attributes,
				Tick:            pk.Tick,
			}
		}
	}

	return pks
}
