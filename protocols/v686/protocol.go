package v686

import (
	_ "embed"

	"github.com/oomph-ac/new-mv/internal/chunk"
	"github.com/oomph-ac/new-mv/mapping"
	"github.com/oomph-ac/new-mv/protocols/latest"
	v686packet "github.com/oomph-ac/new-mv/protocols/v686/packet"
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
	BlockVersion int32 = (1 << 24) | (21 << 16) | (2 << 8) | 2
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
	delete(packetPool_server, packet.IDCurrentStructureFeature)
	delete(packetPool_server, packet.IDJigsawStructureData)

	packetPool_server[packet.IDAddActor] = func() packet.Packet { return &v686packet.AddActor{} }
	packetPool_server[packet.IDAddPlayer] = func() packet.Packet { return &v686packet.AddPlayer{} }
	packetPool_server[packet.IDCameraInstruction] = func() packet.Packet { return &v686packet.CameraInstruction{} }
	packetPool_server[packet.IDCameraPresets] = func() packet.Packet { return &v686packet.CameraPresets{} }
	packetPool_server[packet.IDChangeDimension] = func() packet.Packet { return &v686packet.ChangeDimension{} }
	packetPool_server[packet.IDCorrectPlayerMovePrediction] = func() packet.Packet { return &v686packet.CorrectPlayerMovePrediction{} }
	packetPool_server[packet.IDDisconnect] = func() packet.Packet { return &v686packet.Disconnect{} }
	packetPool_server[packet.IDInventoryContent] = func() packet.Packet { return &v686packet.InventoryContent{} }
	packetPool_server[packet.IDInventorySlot] = func() packet.Packet { return &v686packet.InventorySlot{} }
	packetPool_server[packet.IDItemStackResponse] = func() packet.Packet { return &v686packet.ItemStackResponse{} }
	packetPool_server[packet.IDPlayerArmourDamage] = func() packet.Packet { return &v686packet.PlayerArmourDamage{} }
	packetPool_server[packet.IDResourcePacksInfo] = func() packet.Packet { return &v686packet.ResourcePacksInfo{} }
	packetPool_server[packet.IDSetTitle] = func() packet.Packet { return &v686packet.SetTitle{} }
	packetPool_server[packet.IDStopSound] = func() packet.Packet { return &v686packet.StopSound{} }
	packetPool_server[packet.IDSetActorLink] = func() packet.Packet { return &v686packet.SetActorLink{} }

	packetPool_client = packet.NewClientPool()
	delete(packetPool_client, packet.IDServerBoundLoadingScreen)
	delete(packetPool_client, packet.IDServerBoundDiagnostics)

	// packets used by both client and server...
	packetPool_server[packet.IDEditorNetwork] = func() packet.Packet { return &v686packet.EditorNetwork{} }
	packetPool_server[packet.IDMobArmourEquipment] = func() packet.Packet { return &v686packet.MobArmourEquipment{} }
	packetPool_client[packet.IDEditorNetwork] = func() packet.Packet { return &v686packet.EditorNetwork{} }
	packetPool_client[packet.IDMobArmourEquipment] = func() packet.Packet { return &v686packet.MobArmourEquipment{} }
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
	return 686
}

func (Protocol) Ver() string {
	return "1.21.2"
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
	return p.blockTranslator.UpgradeBlockPackets(
		p.itemTranslator.UpgradeItemPackets(ProtoUpgrade([]packet.Packet{pk}), conn),
		conn,
	)
}

func ProtoUpgrade(pks []packet.Packet) []packet.Packet {
	for index, pk := range pks {
		switch pk := pk.(type) {
		case *v686packet.EditorNetwork:
			pks[index] = &packet.EditorNetwork{
				RouteToManager: false,
				Payload:        pk.Payload,
			}
		case *v686packet.MobArmourEquipment:
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
				transactionData = &t.UseItemTransactionData
			}

			pk.TransactionData = transactionData
			pks[index] = pk
		case *packet.ItemStackRequest:
			for i, req := range pk.Requests {
				for aIndex, action := range req.Actions {
					var newAction protocol.StackRequestAction = action

					switch oldAction := action.(type) {
					case *types.TakeStackRequestAction:
						a := &protocol.TakeStackRequestAction{}
						a.Count = oldAction.Count
						a.Source = protocol.StackRequestSlotInfo{
							Container:      protocol.FullContainerName{ContainerID: oldAction.Source.ContainerID},
							Slot:           oldAction.Source.Slot,
							StackNetworkID: oldAction.Source.StackNetworkID,
						}
						a.Destination = protocol.StackRequestSlotInfo{
							Container:      protocol.FullContainerName{ContainerID: oldAction.Destination.ContainerID},
							Slot:           oldAction.Destination.Slot,
							StackNetworkID: oldAction.Destination.StackNetworkID,
						}
						newAction = a
					case *types.PlaceStackRequestAction:
						a := &protocol.PlaceStackRequestAction{}
						a.Count = oldAction.Count
						a.Source = protocol.StackRequestSlotInfo{
							Container:      protocol.FullContainerName{ContainerID: oldAction.Source.ContainerID},
							Slot:           oldAction.Source.Slot,
							StackNetworkID: oldAction.Source.StackNetworkID,
						}
						a.Destination = protocol.StackRequestSlotInfo{
							Container:      protocol.FullContainerName{ContainerID: oldAction.Destination.ContainerID},
							Slot:           oldAction.Destination.Slot,
							StackNetworkID: oldAction.Destination.StackNetworkID,
						}
						newAction = a
					case *types.SwapStackRequestAction:
						a := &protocol.SwapStackRequestAction{}
						a.Source = protocol.StackRequestSlotInfo{
							Container:      protocol.FullContainerName{ContainerID: oldAction.Source.ContainerID},
							Slot:           oldAction.Source.Slot,
							StackNetworkID: oldAction.Source.StackNetworkID,
						}
						a.Destination = protocol.StackRequestSlotInfo{
							Container:      protocol.FullContainerName{ContainerID: oldAction.Destination.ContainerID},
							Slot:           oldAction.Destination.Slot,
							StackNetworkID: oldAction.Destination.StackNetworkID,
						}
						newAction = a
					case *types.DropStackRequestAction:
						newAction = &protocol.DropStackRequestAction{
							Count: oldAction.Count,
							Source: protocol.StackRequestSlotInfo{
								Container:      protocol.FullContainerName{ContainerID: oldAction.Source.ContainerID},
								Slot:           oldAction.Source.Slot,
								StackNetworkID: oldAction.Source.StackNetworkID,
							},
							Randomly: oldAction.Randomly,
						}
					case *types.DestroyStackRequestAction:
						newAction = &protocol.DestroyStackRequestAction{
							Count: oldAction.Count,
							Source: protocol.StackRequestSlotInfo{
								Container:      protocol.FullContainerName{ContainerID: oldAction.Source.ContainerID},
								Slot:           oldAction.Source.Slot,
								StackNetworkID: oldAction.Source.StackNetworkID,
							},
						}
					case *types.ConsumeStackRequestAction:
						newAction = &protocol.DestroyStackRequestAction{
							Count: oldAction.Count,
							Source: protocol.StackRequestSlotInfo{
								Container:      protocol.FullContainerName{ContainerID: oldAction.Source.ContainerID},
								Slot:           oldAction.Source.Slot,
								StackNetworkID: oldAction.Source.StackNetworkID,
							},
						}
					case *types.PlaceInContainerStackRequestAction:
						a := &protocol.TakeStackRequestAction{}
						a.Count = oldAction.Count
						a.Source = protocol.StackRequestSlotInfo{
							Container:      protocol.FullContainerName{ContainerID: oldAction.Source.ContainerID},
							Slot:           oldAction.Source.Slot,
							StackNetworkID: oldAction.Source.StackNetworkID,
						}
						a.Destination = protocol.StackRequestSlotInfo{
							Container:      protocol.FullContainerName{ContainerID: oldAction.Destination.ContainerID},
							Slot:           oldAction.Destination.Slot,
							StackNetworkID: oldAction.Destination.StackNetworkID,
						}
						newAction = a
					case *types.TakeOutContainerStackRequestAction:
						a := &protocol.TakeStackRequestAction{}
						a.Count = oldAction.Count
						a.Source = protocol.StackRequestSlotInfo{
							Container:      protocol.FullContainerName{ContainerID: oldAction.Source.ContainerID},
							Slot:           oldAction.Source.Slot,
							StackNetworkID: oldAction.Source.StackNetworkID,
						}
						a.Destination = protocol.StackRequestSlotInfo{
							Container:      protocol.FullContainerName{ContainerID: oldAction.Destination.ContainerID},
							Slot:           oldAction.Destination.Slot,
							StackNetworkID: oldAction.Destination.StackNetworkID,
						}
						newAction = a
					case *types.LabTableCombineStackRequestAction:
						newAction = &protocol.LabTableCombineStackRequestAction{}
					case *types.CraftRecipeStackRequestAction:
						newAction = &oldAction.CraftRecipeStackRequestAction
					case *types.AutoCraftRecipeStackRequestAction:
						newAction = &oldAction.AutoCraftRecipeStackRequestAction
					case *types.CraftCreativeStackRequestAction:
						newAction = &oldAction.CraftCreativeStackRequestAction
					case *types.CraftRecipeOptionalStackRequestAction:
						newAction = &oldAction.CraftRecipeOptionalStackRequestAction
					case *types.CraftGrindstoneRecipeStackRequestAction:
						newAction = &oldAction.CraftGrindstoneRecipeStackRequestAction
					}

					req.Actions[aIndex] = newAction
				}

				pk.Requests[i] = req
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

			pks[index] = &v686packet.AddActor{
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

			pks[index] = &v686packet.AddPlayer{
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
			pks[index] = &v686packet.CameraInstruction{
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

			pks[index] = &v686packet.CameraPresets{
				Presets: presets,
			}
		case *packet.ChangeDimension:
			pks[index] = &v686packet.ChangeDimension{
				Dimension: pk.Dimension,
				Position:  pk.Position,
				Respawn:   pk.Respawn,
			}
		case *packet.CorrectPlayerMovePrediction:
			pks[index] = &v686packet.CorrectPlayerMovePrediction{
				PredictionType: pk.PredictionType,
				Position:       pk.Position,
				Delta:          pk.Delta,
				Rotation:       pk.Rotation,
				OnGround:       pk.OnGround,
				Tick:           pk.Tick,
			}
		case *packet.Disconnect:
			pks[index] = &v686packet.Disconnect{
				Reason:                  pk.Reason,
				HideDisconnectionScreen: pk.HideDisconnectionScreen,
				Message:                 pk.Message,
			}
		case *packet.EditorNetwork:
			pks[index] = &v686packet.EditorNetwork{
				Payload: pk.Payload,
			}
		case *packet.InventoryContent:
			pks[index] = &v686packet.InventoryContent{
				WindowID: pk.WindowID,
				Content:  pk.Content,
			}
		case *packet.InventorySlot:
			pks[index] = &v686packet.InventorySlot{
				WindowID: pk.WindowID,
				Slot:     pk.Slot,
				NewItem:  pk.NewItem,
			}
		case *packet.ItemStackResponse:
			responses := make([]types.ItemStackResponse, len(pk.Responses))
			for index, response := range pk.Responses {
				tResponse := types.ItemStackResponse{}
				tResponse.Status = response.Status
				tResponse.RequestID = response.RequestID

				containerInfo := make([]types.StackResponseContainerInfo, len(response.ContainerInfo))
				for cIndex, info := range response.ContainerInfo {
					containerInfo[cIndex] = types.StackResponseContainerInfo{StackResponseContainerInfo: info}
				}
				tResponse.ContainerInfo = containerInfo

				responses[index] = tResponse
			}

			pks[index] = &v686packet.ItemStackResponse{
				Responses: responses,
			}
		case *packet.MobArmourEquipment:
			pks[index] = &v686packet.MobArmourEquipment{
				EntityRuntimeID: pk.EntityRuntimeID,
				Helmet:          pk.Helmet,
				Chestplate:      pk.Chestplate,
				Leggings:        pk.Leggings,
				Boots:           pk.Boots,
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

			pks[index] = &v686packet.PlayerArmourDamage{
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

			pks[index] = &v686packet.ResourcePacksInfo{
				TexturePackRequired: pk.TexturePackRequired,
				HasAddons:           pk.HasAddons,
				HasScripts:          pk.HasScripts,
				BehaviourPacks:      []types.TexturePackInfo{},
				TexturePacks:        tPacks,
				ForcingServerPacks:  true,
				PackURLs:            packURLs,
			}
		case *packet.SetTitle:
			pks[index] = &v686packet.SetTitle{
				ActionType:       pk.ActionType,
				Text:             pk.Text,
				FadeInDuration:   pk.FadeInDuration,
				RemainDuration:   pk.RemainDuration,
				FadeOutDuration:  pk.FadeOutDuration,
				XUID:             pk.XUID,
				PlatformOnlineID: pk.PlatformOnlineID,
			}
		case *packet.StopSound:
			pks[index] = &v686packet.StopSound{
				SoundName: pk.SoundName,
				StopAll:   pk.StopAll,
			}
		case *packet.SetActorLink:
			pks[index] = &v686packet.SetActorLink{
				EntityLink: types.EntityLink{
					EntityLink: pk.EntityLink,
				},
			}
		}
	}

	return pks
}
