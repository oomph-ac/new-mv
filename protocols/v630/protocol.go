package v630

import (
	_ "embed"

	"github.com/oomph-ac/new-mv/internal/chunk"
	"github.com/oomph-ac/new-mv/mapping"
	"github.com/oomph-ac/new-mv/protocols/latest"
	v630packet "github.com/oomph-ac/new-mv/protocols/v630/packet"
	"github.com/oomph-ac/new-mv/protocols/v630/types"
	"github.com/oomph-ac/new-mv/translator"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

const (
	// ItemVersion is the version of items of the game which use for downgrading and upgrading.
	ItemVersion = 161
	// BlockVersion is the version of blocks (states) of the game. This version is composed
	// of 4 bytes indicating a version, interpreted as a big endian int. The current version represents
	// 1.20.70.0
	BlockVersion int32 = (1 << 24) | (20 << 16) | (60 << 8)
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
	// ------------------------ 1.21.30 changes ------------------------

	// ------------------------ 1.21.20 changes ------------------------
	delete(packetPool_server, packet.IDCameraAimAssist)
	delete(packetPool_server, packet.IDContainerRegistryCleanup)

	packetPool_server[packet.IDEmote] = func() packet.Packet { return &v630packet.Emote{} }
	packetPool_client[packet.IDEmote] = func() packet.Packet { return &v630packet.Emote{} }

	packetPool_server[packet.IDCameraPresets] = func() packet.Packet { return &v630packet.CameraPresets{} }
	packetPool_server[packet.IDContainerRegistryCleanup] = func() packet.Packet { return &v630packet.ContainerRegistryCleanup{} }
	packetPool_server[packet.IDItemStackResponse] = func() packet.Packet { return &v630packet.ItemStackResponse{} }
	packetPool_server[packet.IDResourcePacksInfo] = func() packet.Packet { return &v630packet.ResourcePacksInfo{} }
	packetPool_server[packet.IDTransfer] = func() packet.Packet { return &v630packet.Transfer{} }
	packetPool_server[packet.IDUpdateAttributes] = func() packet.Packet { return &v630packet.UpdateAttributes{} }
	// ------------------------ 1.21.20 changes ------------------------

	// ------------------------ 1.21.2 changes ------------------------
	delete(packetPool_server, packet.IDCurrentStructureFeature)
	delete(packetPool_server, packet.IDJigsawStructureData)
	delete(packetPool_client, packet.IDServerBoundDiagnostics)
	delete(packetPool_client, packet.IDServerBoundLoadingScreen)

	packetPool_server[packet.IDMobArmourEquipment] = func() packet.Packet { return &v630packet.MobArmourEquipment{} }
	packetPool_client[packet.IDMobArmourEquipment] = func() packet.Packet { return &v630packet.MobArmourEquipment{} }

	packetPool_server[packet.IDEditorNetwork] = func() packet.Packet { return &v630packet.EditorNetwork{} }
	packetPool_client[packet.IDEditorNetwork] = func() packet.Packet { return &v630packet.EditorNetwork{} }

	packetPool_server[packet.IDAddActor] = func() packet.Packet { return &v630packet.AddActor{} }
	packetPool_server[packet.IDAddPlayer] = func() packet.Packet { return &v630packet.AddPlayer{} }
	packetPool_server[packet.IDCameraInstruction] = func() packet.Packet { return &v630packet.CameraInstruction{} }
	packetPool_server[packet.IDChangeDimension] = func() packet.Packet { return &v630packet.ChangeDimension{} }
	packetPool_server[packet.IDCompressedBiomeDefinitionList] = func() packet.Packet { return &v630packet.CompressedBiomeDefinitionList{} }
	packetPool_server[packet.IDDisconnect] = func() packet.Packet { return &v630packet.Disconnect{} }
	packetPool_server[packet.IDInventoryContent] = func() packet.Packet { return &v630packet.InventoryContent{} }
	packetPool_server[packet.IDInventorySlot] = func() packet.Packet { return &v630packet.InventorySlot{} }
	packetPool_server[packet.IDPlayerArmourDamage] = func() packet.Packet { return &v630packet.PlayerArmourDamage{} }
	packetPool_server[packet.IDSetTitle] = func() packet.Packet { return &v630packet.SetTitle{} }
	packetPool_server[packet.IDStopSound] = func() packet.Packet { return &v630packet.StopSound{} }
	// ------------------------ 1.21.2 changes ------------------------

	// ------------------------ 1.20.80 changes ------------------------
	delete(packetPool_server, packet.IDAwardAchievement)
	packetPool_server[packet.IDContainerClose] = func() packet.Packet { return &v630packet.ContainerClose{} }
	packetPool_client[packet.IDContainerClose] = func() packet.Packet { return &v630packet.ContainerClose{} }
	packetPool_server[packet.IDText] = func() packet.Packet { return &v630packet.Text{} }
	packetPool_client[packet.IDText] = func() packet.Packet { return &v630packet.Text{} }
	// ------------------------ 1.20.80 changes ------------------------

	// ------------------------ 1.20.70 changes ------------------------
	packetPool_server[packet.IDCorrectPlayerMovePrediction] = func() packet.Packet { return &v630packet.CorrectPlayerMovePrediction{} }
	packetPool_server[packet.IDResourcePackStack] = func() packet.Packet { return &v630packet.ResourcePackStack{} }
	packetPool_server[packet.IDStartGame] = func() packet.Packet { return &v630packet.StartGame{} }
	packetPool_server[packet.IDUpdateBlockSynced] = func() packet.Packet { return &v630packet.UpdateBlockSynced{} }
	packetPool_server[packet.IDUpdatePlayerGameType] = func() packet.Packet { return &v630packet.UpdatePlayerGameType{} }
	// ------------------------ 1.20.70 changes ------------------------

	// ------------------------ 1.20.60 changes ------------------------
	// NOTE: Packet with ID 71 (ItemFrameDropItem is server bound and should be cancelled).
	packetPool_client[packet.IDLecternUpdate] = func() packet.Packet { return &v630packet.LecternUpdate{} }
	packetPool_client[packet.IDPlayerAuthInput] = func() packet.Packet { return &v630packet.PlayerAuthInput{} }

	packetPool_server[packet.IDMobEffect] = func() packet.Packet { return &v630packet.MobEffect{} }
	packetPool_server[packet.IDSetActorMotion] = func() packet.Packet { return &v630packet.SetActorMotion{} }
	// ------------------------ 1.20.60 changes ------------------------

	// ------------------------ 1.20.50 changes ------------------------
	packetPool_server[packet.IDPlayerList] = func() packet.Packet { return &v630packet.PlayerList{} }
	packetPool_server[packet.IDLevelChunk] = func() packet.Packet { return &v630packet.LevelChunk{} }
	delete(packetPool_server, packet.IDSetHud)
	// ------------------------ 1.20.50 changes ------------------------
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
	return 630
}

func (Protocol) Ver() string {
	return "1.20.50"
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
	return p.blockTranslator.UpgradeBlockPackets(p.itemTranslator.UpgradeItemPackets(
		ProtoUpgrade([]packet.Packet{pk}), conn),
		conn,
	)
}

func ProtoUpgrade(pks []packet.Packet) []packet.Packet {
	for index, pk := range pks {
		switch pk := pk.(type) {
		case *v630packet.ContainerClose:
			pks[index] = &packet.ContainerClose{
				WindowID:   pk.WindowID,
				ServerSide: pk.ServerSide,
			}
		case *v630packet.Emote:
			pks[index] = &packet.Emote{
				EntityRuntimeID: pk.EntityRuntimeID,
				EmoteID:         pk.EmoteID,
				EmoteLength:     100, // TODO: ???
				XUID:            pk.XUID,
				PlatformID:      pk.PlatformID,
				Flags:           pk.Flags,
			}
		case *v630packet.EditorNetwork:
			pks[index] = &packet.EditorNetwork{
				RouteToManager: false,
				Payload:        pk.Payload,
			}
		case *v630packet.LecternUpdate:
			pks[index] = &packet.LecternUpdate{
				Page:      pk.Page,
				PageCount: pk.PageCount,
				Position:  pk.Position,
			}
		case *v630packet.MobArmourEquipment:
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
		case *v630packet.PlayerAuthInput:
			pks[index] = &packet.PlayerAuthInput{
				Pitch:               pk.Pitch,
				Yaw:                 pk.Yaw,
				Position:            pk.Position,
				MoveVector:          pk.MoveVector,
				HeadYaw:             pk.HeadYaw,
				InputData:           pk.InputData,
				InputMode:           pk.InputMode,
				PlayMode:            pk.PlayMode,
				InteractionModel:    pk.InteractionModel,
				InteractPitch:       pk.GazeDirection.X(),
				InteractYaw:         pk.GazeDirection.Y(),
				Tick:                pk.Tick,
				Delta:               pk.Delta,
				ItemInteractionData: pk.ItemInteractionData,
				ItemStackRequest:    pk.ItemStackRequest,
				BlockActions:        pk.BlockActions,
				AnalogueMoveVector:  pk.AnalogueMoveVector,
			}
		case *v630packet.Text:
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

			pks[index] = &v630packet.AddActor{
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

			pks[index] = &v630packet.AddPlayer{
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
		case *packet.AvailableCommands:
			for cmdIndex, cmd := range pk.Commands {
				for overloadIndex, overload := range cmd.Overloads {
					for paramIndex, param := range overload.Parameters {
						var newT uint32 = protocol.CommandArgValid

						switch param.Type ^ protocol.CommandArgValid {
						case protocol.CommandArgTypeEquipmentSlots:
							newT |= types.CommandArgTypeEquipmentSlots
						case protocol.CommandArgTypeString:
							newT |= types.CommandArgTypeString
						case protocol.CommandArgTypeBlockPosition:
							newT |= types.CommandArgTypeBlockPosition
						case protocol.CommandArgTypePosition:
							newT |= types.CommandArgTypePosition
						case protocol.CommandArgTypeMessage:
							newT |= types.CommandArgTypeMessage
						case protocol.CommandArgTypeRawText:
							newT |= types.CommandArgTypeRawText
						case protocol.CommandArgTypeJSON:
							newT |= types.CommandArgTypeJSON
						case protocol.CommandArgTypeBlockStates:
							newT |= types.CommandArgTypeBlockStates
						case protocol.CommandArgTypeCommand:
							newT |= types.CommandArgTypeCommand
						default:
							newT = param.Type
						}
						param.Type = newT

						overload.Parameters[paramIndex] = param
					}

					cmd.Overloads[overloadIndex] = overload
				}

				pk.Commands[cmdIndex] = cmd
			}
		case *packet.CameraInstruction:
			pks[index] = &v630packet.CameraInstruction{
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

			pks[index] = &v630packet.CameraPresets{
				Presets: presets,
			}
		case *packet.ChangeDimension:
			pks[index] = &v630packet.ChangeDimension{
				Dimension: pk.Dimension,
				Position:  pk.Position,
				Respawn:   pk.Respawn,
			}
		case *packet.ContainerClose:
			pks[index] = &v630packet.ContainerClose{
				WindowID:   pk.WindowID,
				ServerSide: pk.ServerSide,
			}
		case *packet.ContainerRegistryCleanup:
			containers := make([]types.FullContainerName, len(pk.RemovedContainers))
			for index, container := range pk.RemovedContainers {
				containers[index] = types.DowngradeContainer(container)
			}

			pks[index] = &v630packet.ContainerRegistryCleanup{
				RemovedContainers: containers,
			}
		case *packet.CorrectPlayerMovePrediction:
			pks[index] = &v630packet.CorrectPlayerMovePrediction{
				Position: pk.Position,
				Delta:    pk.Delta,
				OnGround: pk.OnGround,
				Tick:     pk.Tick,
			}
		case *packet.CraftingData:
			pk.Recipes = types.DowngradeRecipes(pk.Recipes)
		case *packet.Disconnect:
			pks[index] = &v630packet.Disconnect{
				Reason:                  pk.Reason,
				HideDisconnectionScreen: pk.HideDisconnectionScreen,
				Message:                 pk.Message,
			}
		case *packet.EditorNetwork:
			pks[index] = &v630packet.EditorNetwork{
				Payload: pk.Payload,
			}
		case *packet.Emote:
			pks[index] = &v630packet.Emote{
				EntityRuntimeID: pk.EntityRuntimeID,
				EmoteID:         pk.EmoteID,
				XUID:            pk.XUID,
				PlatformID:      pk.PlatformID,
				Flags:           pk.Flags,
			}
		case *packet.InventoryContent:
			pks[index] = &v630packet.InventoryContent{
				WindowID: pk.WindowID,
				Content:  pk.Content,
			}
		case *packet.InventorySlot:
			pks[index] = &v630packet.InventorySlot{
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

			pks[index] = &v630packet.ItemStackResponse{
				Responses: responses,
			}
		case *packet.LevelChunk:
			pks[index] = &v630packet.LevelChunk{
				Position:        pk.Position,
				HighestSubChunk: pk.HighestSubChunk,
				SubChunkCount:   pk.SubChunkCount,
				CacheEnabled:    pk.CacheEnabled,
				BlobHashes:      pk.BlobHashes,
				RawPayload:      pk.RawPayload,
			}
		case *packet.MobArmourEquipment:
			pks[index] = &v630packet.MobArmourEquipment{
				EntityRuntimeID: pk.EntityRuntimeID,
				Helmet:          pk.Helmet,
				Chestplate:      pk.Chestplate,
				Leggings:        pk.Leggings,
				Boots:           pk.Boots,
			}
		case *packet.MobEffect:
			pks[index] = &v630packet.MobEffect{
				EntityRuntimeID: pk.EntityRuntimeID,
				Operation:       pk.Operation,
				EffectType:      pk.EffectType,
				Amplifier:       pk.Amplifier,
				Particles:       pk.Particles,
				Duration:        pk.Duration,
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

			pks[index] = &v630packet.PlayerArmourDamage{
				Bitset:           bitset,
				HelmetDamage:     pk.HelmetDamage,
				ChestplateDamage: pk.ChestplateDamage,
				LeggingsDamage:   pk.LeggingsDamage,
				BootsDamage:      pk.BootsDamage,
			}
		case *packet.PlayerList:
			entries := make([]types.PlayerListEntry, len(pk.Entries))
			for index, entry := range pk.Entries {
				entries[index] = types.PlayerListEntry{
					UUID:           entry.UUID,
					EntityUniqueID: entry.EntityUniqueID,
					Username:       entry.Username,
					XUID:           entry.XUID,
					PlatformChatID: entry.PlatformChatID,
					BuildPlatform:  entry.BuildPlatform,
					Skin:           entry.Skin,
					Teacher:        entry.Teacher,
					Host:           entry.Host,
				}
			}

			pks[index] = &v630packet.PlayerList{
				ActionType: pk.ActionType,
				Entries:    entries,
			}
		case *packet.ResourcePackStack:
			pks[index] = &v630packet.ResourcePackStack{
				TexturePackRequired:          pk.TexturePackRequired,
				BehaviourPacks:               pk.BehaviourPacks,
				TexturePacks:                 pk.TexturePacks,
				BaseGameVersion:              pk.BaseGameVersion,
				Experiments:                  pk.Experiments,
				ExperimentsPreviouslyToggled: pk.ExperimentsPreviouslyToggled,
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

			pks[index] = &v630packet.ResourcePacksInfo{
				TexturePackRequired: pk.TexturePackRequired,
				HasScripts:          pk.HasScripts,
				BehaviourPacks:      []types.TexturePackInfo{},
				TexturePacks:        tPacks,
				ForcingServerPacks:  true,
				PackURLs:            packURLs,
			}
		case *packet.SetTitle:
			pks[index] = &v630packet.SetTitle{
				ActionType:       pk.ActionType,
				Text:             pk.Text,
				FadeInDuration:   pk.FadeInDuration,
				RemainDuration:   pk.RemainDuration,
				FadeOutDuration:  pk.FadeOutDuration,
				XUID:             pk.XUID,
				PlatformOnlineID: pk.PlatformOnlineID,
			}
		case *packet.StartGame:
			pks[index] = &v630packet.StartGame{
				EntityUniqueID:                 pk.EntityUniqueID,
				EntityRuntimeID:                pk.EntityRuntimeID,
				PlayerGameMode:                 pk.PlayerGameMode,
				PlayerPosition:                 pk.PlayerPosition,
				Pitch:                          pk.Pitch,
				Yaw:                            pk.Yaw,
				WorldSeed:                      pk.WorldSeed,
				SpawnBiomeType:                 pk.SpawnBiomeType,
				UserDefinedBiomeName:           pk.UserDefinedBiomeName,
				Dimension:                      pk.Dimension,
				Generator:                      pk.Generator,
				WorldGameMode:                  pk.WorldGameMode,
				Hardcore:                       pk.Hardcore,
				Difficulty:                     pk.Difficulty,
				WorldSpawn:                     pk.WorldSpawn,
				AchievementsDisabled:           pk.AchievementsDisabled,
				EditorWorldType:                pk.EditorWorldType,
				CreatedInEditor:                pk.CreatedInEditor,
				ExportedFromEditor:             pk.ExportedFromEditor,
				DayCycleLockTime:               pk.DayCycleLockTime,
				EducationEditionOffer:          pk.EducationEditionOffer,
				EducationFeaturesEnabled:       pk.EducationFeaturesEnabled,
				EducationProductID:             pk.EducationProductID,
				RainLevel:                      pk.RainLevel,
				LightningLevel:                 pk.LightningLevel,
				ConfirmedPlatformLockedContent: pk.ConfirmedPlatformLockedContent,
				MultiPlayerGame:                pk.MultiPlayerGame,
				LANBroadcastEnabled:            pk.LANBroadcastEnabled,
				XBLBroadcastMode:               pk.XBLBroadcastMode,
				PlatformBroadcastMode:          pk.PlatformBroadcastMode,
				CommandsEnabled:                pk.CommandsEnabled,
				TexturePackRequired:            pk.TexturePackRequired,
				GameRules:                      pk.GameRules,
				Experiments:                    pk.Experiments,
				ExperimentsPreviouslyToggled:   pk.ExperimentsPreviouslyToggled,
				BonusChestEnabled:              pk.BonusChestEnabled,
				StartWithMapEnabled:            pk.StartWithMapEnabled,
				PlayerPermissions:              pk.PlayerPermissions,
				ServerChunkTickRadius:          pk.ServerChunkTickRadius,
				HasLockedBehaviourPack:         pk.HasLockedBehaviourPack,
				HasLockedTexturePack:           pk.HasLockedTexturePack,
				FromLockedWorldTemplate:        pk.FromLockedWorldTemplate,
				MSAGamerTagsOnly:               pk.MSAGamerTagsOnly,
				FromWorldTemplate:              pk.FromWorldTemplate,
				WorldTemplateSettingsLocked:    pk.WorldTemplateSettingsLocked,
				OnlySpawnV1Villagers:           pk.OnlySpawnV1Villagers,
				PersonaDisabled:                pk.PersonaDisabled,
				CustomSkinsDisabled:            pk.CustomSkinsDisabled,
				EmoteChatMuted:                 pk.EmoteChatMuted,
				BaseGameVersion:                pk.BaseGameVersion,
				LimitedWorldWidth:              pk.LimitedWorldWidth,
				LimitedWorldDepth:              pk.LimitedWorldDepth,
				NewNether:                      pk.NewNether,
				EducationSharedResourceURI:     pk.EducationSharedResourceURI,
				ForceExperimentalGameplay:      pk.ForceExperimentalGameplay,
				LevelID:                        pk.LevelID,
				WorldName:                      pk.WorldName,
				TemplateContentIdentity:        pk.TemplateContentIdentity,
				Trial:                          pk.Trial,
				PlayerMovementSettings:         pk.PlayerMovementSettings,
				Time:                           pk.Time,
				EnchantmentSeed:                pk.EnchantmentSeed,
				Blocks:                         pk.Blocks,
				Items:                          pk.Items,
				MultiPlayerCorrelationID:       pk.MultiPlayerCorrelationID,
				ServerAuthoritativeInventory:   pk.ServerAuthoritativeInventory,
				GameVersion:                    pk.GameVersion,
				PropertyData:                   pk.PropertyData,
				ServerBlockStateChecksum:       pk.ServerBlockStateChecksum,
				ClientSideGeneration:           pk.ClientSideGeneration,
				WorldTemplateID:                pk.WorldTemplateID,
				ChatRestrictionLevel:           pk.ChatRestrictionLevel,
				DisablePlayerInteractions:      pk.DisablePlayerInteractions,
				UseBlockNetworkIDHashes:        pk.UseBlockNetworkIDHashes,
				ServerAuthoritativeSound:       pk.ServerAuthoritativeSound,
			}
		case *packet.StopSound:
			pks[index] = &v630packet.StopSound{
				SoundName: pk.SoundName,
				StopAll:   pk.StopAll,
			}
		case *packet.SetActorLink:
			pks[index] = &v630packet.SetActorLink{
				EntityLink: types.EntityLink{
					EntityLink: pk.EntityLink,
				},
			}
		case *packet.SetActorMotion:
			pks[index] = &v630packet.SetActorMotion{
				EntityRuntimeID: pk.EntityRuntimeID,
				Velocity:        pk.Velocity,
			}
		case *packet.Text:
			pks[index] = &v630packet.Text{
				TextType:         pk.TextType,
				NeedsTranslation: pk.NeedsTranslation,
				SourceName:       pk.SourceName,
				Message:          pk.Message,
				Parameters:       pk.Parameters,
				XUID:             pk.XUID,
				PlatformChatID:   pk.PlatformChatID,
			}
		case *packet.Transfer:
			pks[index] = &v630packet.Transfer{
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

			pks[index] = &v630packet.UpdateAttributes{
				EntityRuntimeID: pk.EntityRuntimeID,
				Attributes:      attributes,
				Tick:            pk.Tick,
			}
		case *packet.UpdateBlockSynced:
			pks[index] = &v630packet.UpdateBlockSynced{
				Position:          pk.Position,
				NewBlockRuntimeID: pk.NewBlockRuntimeID,
				Flags:             pk.Flags,
				Layer:             pk.Layer,
				EntityUniqueID:    int64(pk.EntityUniqueID),
				TransitionType:    pk.TransitionType,
			}
		case *packet.UpdatePlayerGameType:
			pks[index] = &v630packet.UpdatePlayerGameType{
				GameType:       pk.GameType,
				PlayerUniqueID: pk.PlayerUniqueID,
			}
		}
	}

	return pks
}
