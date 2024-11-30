package types

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

// ItemStackRequest represents a single request present in an ItemStackRequest packet sent by the client to
// change an item in an inventory.
// Item stack requests are either approved or rejected by the server using the ItemStackResponse packet.
type ItemStackRequest struct {
	// RequestID is a unique ID for the request. This ID is used by the server to send a response for this
	// specific request in the ItemStackResponse packet.
	RequestID int32
	// Actions is a list of actions performed by the client. The actual type of the actions depends on which
	// ID was present, and is one of the concrete types below.
	Actions []protocol.StackRequestAction
	// FilterStrings is a list of filter strings involved in the request. This is typically filled with one string
	// when an anvil or cartography is used.
	FilterStrings []string
	// FilterCause represents the cause of any potential filtering. This is one of the constants above.
	FilterCause int32
}

// Marshal encodes/decodes an ItemStackRequest.
func (x *ItemStackRequest) Marshal(r protocol.IO) {
	r.Varint32(&x.RequestID)
	protocol.FuncSlice(r, &x.Actions, r.StackRequestAction)
	protocol.FuncSlice(r, &x.FilterStrings, r.String)
	r.Int32(&x.FilterCause)
}

// ItemStackResponse is a response to an individual ItemStackRequest.
type ItemStackResponse struct {
	// Status specifies if the request with the RequestID below was successful. If this is the case, the
	// ContainerInfo below will have information on what slots ended up changing. If not, the container info
	// will be empty.
	// A non-0 status means an error occurred and will result in the action being reverted.
	Status uint8
	// RequestID is the unique ID of the request that this response is in reaction to. If rejected, the client
	// will undo the actions from the request with this ID.
	RequestID int32
	// ContainerInfo holds information on the containers that had their contents changed as a result of the
	// request.
	ContainerInfo []StackResponseContainerInfo
}

// Marshal encodes/decodes an ItemStackResponse.
func (x *ItemStackResponse) Marshal(r protocol.IO) {
	r.Uint8(&x.Status)
	r.Varint32(&x.RequestID)
	if x.Status == protocol.ItemStackResponseStatusOK {
		protocol.Slice(r, &x.ContainerInfo)
	}
}

// StackResponseContainerInfo holds information on what slots in a container have what item stack in them.
type StackResponseContainerInfo struct {
	// Container is the FullContainerName that describes the container that the slots that follow are in. For
	// the main inventory, the ContainerID seems to be 0x1b. Fur the cursor, this value seems to be 0x3a. For
	// the crafting grid, this value seems to be 0x0d.
	Container FullContainerName
	// SlotInfo holds information on what item stack should be present in specific slots in the container.
	SlotInfo []protocol.StackResponseSlotInfo
}

// Marshal encodes/decodes a StackResponseContainerInfo.
func (x *StackResponseContainerInfo) Marshal(r protocol.IO) {
	r.Uint8(&x.Container.ContainerID)
	protocol.Slice(r, &x.SlotInfo)
}

// StackRequestSlotInfo holds information on a specific slot client-side.
type StackRequestSlotInfo struct {
	// Container is the FullContainerName that describes the container that the slot is in.
	Container FullContainerName
	// Slot is the index of the slot within the container with the ContainerID above.
	Slot byte
	// StackNetworkID is the unique stack ID that the client assumes to be present in this slot. The server
	// must check if these IDs match. If they do not match, servers should reject the stack request that the
	// action holding this info was in.
	StackNetworkID int32
}

// StackReqSlotInfo reads/writes a StackRequestSlotInfo x using IO r.
func StackReqSlotInfo(r protocol.IO, x *StackRequestSlotInfo) {
	r.Uint8(&x.Container.ContainerID)
	r.Uint8(&x.Slot)
	r.Varint32(&x.StackNetworkID)
}

// slot into another.
type transferStackRequestAction struct {
	// Count is the count of the item in the source slot that was taken towards the destination slot.
	Count byte
	// Source and Destination point to the source slot from which Count of the item stack were taken and the
	// destination slot to which this item was moved.
	Source, Destination StackRequestSlotInfo
}

// Marshal ...
func (a *transferStackRequestAction) Marshal(r protocol.IO) {
	r.Uint8(&a.Count)
	StackReqSlotInfo(r, &a.Source)
	StackReqSlotInfo(r, &a.Destination)
}

type TakeStackRequestAction struct {
	transferStackRequestAction
}

type PlaceStackRequestAction struct {
	transferStackRequestAction
}

type SwapStackRequestAction struct {
	// Source and Destination point to the source slot from which Count of the item stack were taken and the
	// destination slot to which this item was moved.
	Source, Destination StackRequestSlotInfo
}

// Marshal ...
func (a *SwapStackRequestAction) Marshal(r protocol.IO) {
	StackReqSlotInfo(r, &a.Source)
	StackReqSlotInfo(r, &a.Destination)
}

type DropStackRequestAction struct {
	// Count is the count of the item in the source slot that was taken towards the destination slot.
	Count byte
	// Source is the source slot from which items were dropped to the ground.
	Source StackRequestSlotInfo
	// Randomly seems to be set to false in most cases. I'm not entirely sure what this does, but this is what
	// vanilla calls this field.
	Randomly bool
}

// Marshal ...
func (a *DropStackRequestAction) Marshal(r protocol.IO) {
	r.Uint8(&a.Count)
	StackReqSlotInfo(r, &a.Source)
	r.Bool(&a.Randomly)
}

type DestroyStackRequestAction struct {
	// Count is the count of the item in the source slot that was destroyed.
	Count byte
	// Source is the source slot from which items came that were destroyed by moving them into the creative
	// inventory.
	Source StackRequestSlotInfo
}

// Marshal ...
func (a *DestroyStackRequestAction) Marshal(r protocol.IO) {
	r.Uint8(&a.Count)
	StackReqSlotInfo(r, &a.Source)
}

type ConsumeStackRequestAction struct {
	DestroyStackRequestAction
}

// PlaceInContainerStackRequestAction currently has no known purpose.
type PlaceInContainerStackRequestAction struct {
	transferStackRequestAction
}

// TakeOutContainerStackRequestAction currently has no known purpose.
type TakeOutContainerStackRequestAction struct {
	transferStackRequestAction
}

// LabTableCombineStackRequestAction is sent by the client when it uses a lab table to combine item stacks.
type LabTableCombineStackRequestAction struct{}

// Marshal ...
func (a *LabTableCombineStackRequestAction) Marshal(protocol.IO) {}

// CraftRecipeStackRequestAction is sent by the client the moment it begins crafting an item. This is the
// first action sent, before the Consume and Create item stack request actions.
// This action is also sent when an item is enchanted. Enchanting should be treated mostly the same way as
// crafting, where the old item is consumed.
type CraftRecipeStackRequestAction struct {
	// RecipeNetworkID is the network ID of the recipe that is about to be crafted. This network ID matches
	// one of the recipes sent in the CraftingData packet, where each of the recipes have a RecipeNetworkID as
	// of 1.16.
	RecipeNetworkID uint32
}

// Marshal ...
func (a *CraftRecipeStackRequestAction) Marshal(r protocol.IO) {
	r.Varuint32(&a.RecipeNetworkID)
}

// AutoCraftRecipeStackRequestAction is sent by the client similarly to the CraftRecipeStackRequestAction. The
// only difference is that the recipe is automatically created and crafted by shift clicking the recipe book.
type AutoCraftRecipeStackRequestAction struct {
	// RecipeNetworkID is the network ID of the recipe that is about to be crafted. This network ID matches
	// one of the recipes sent in the CraftingData packet, where each of the recipes have a RecipeNetworkID as
	// of 1.16.
	RecipeNetworkID uint32
	// TimesCrafted is how many times the recipe was crafted.
	TimesCrafted byte
	// Ingredients is a slice of ItemDescriptorCount that contains the ingredients that were used to craft the recipe.
	// It is not exactly clear what this is used for, but it is sent by the vanilla client.
	Ingredients []protocol.ItemDescriptorCount
}

// Marshal ...
func (a *AutoCraftRecipeStackRequestAction) Marshal(r protocol.IO) {
	r.Varuint32(&a.RecipeNetworkID)
	r.Uint8(&a.TimesCrafted)
	protocol.FuncSlice(r, &a.Ingredients, r.ItemDescriptorCount)
}

// CraftCreativeStackRequestAction is sent by the client when it takes an item out fo the creative inventory.
// The item is thus not really crafted, but instantly created.
type CraftCreativeStackRequestAction struct {
	// CreativeItemNetworkID is the network ID of the creative item that is being created. This is one of the
	// creative item network IDs sent in the CreativeContent packet.
	CreativeItemNetworkID uint32
}

// Marshal ...
func (a *CraftCreativeStackRequestAction) Marshal(r protocol.IO) {
	r.Varuint32(&a.CreativeItemNetworkID)
}

// CraftRecipeOptionalStackRequestAction is sent when using an anvil. When this action is sent, the
// FilterStrings field in the respective stack request is non-empty and contains the name of the item created
// using the anvil or cartography table.
type CraftRecipeOptionalStackRequestAction struct {
	// RecipeNetworkID is the network ID of the multi-recipe that is about to be crafted. This network ID matches
	// one of the multi-recipes sent in the CraftingData packet, where each of the recipes have a RecipeNetworkID as
	// of 1.16.
	RecipeNetworkID uint32
	// FilterStringIndex is the index of a filter string sent in a ItemStackRequest.
	FilterStringIndex int32
}

// Marshal ...
func (c *CraftRecipeOptionalStackRequestAction) Marshal(r protocol.IO) {
	r.Varuint32(&c.RecipeNetworkID)
	r.Int32(&c.FilterStringIndex)
}

// CraftGrindstoneRecipeStackRequestAction is sent when a grindstone recipe is crafted. It contains the RecipeNetworkID
// to identify the recipe crafted, and the cost for crafting the recipe.
type CraftGrindstoneRecipeStackRequestAction struct {
	// RecipeNetworkID is the network ID of the recipe that is about to be crafted. This network ID matches
	// one of the recipes sent in the CraftingData packet, where each of the recipes have a RecipeNetworkID as
	// of 1.16.
	RecipeNetworkID uint32
	// Cost is the cost of the recipe that was crafted.
	Cost int32
}

// Marshal ...
func (c *CraftGrindstoneRecipeStackRequestAction) Marshal(r protocol.IO) {
	r.Varuint32(&c.RecipeNetworkID)
	r.Varint32(&c.Cost)
}

func UpgradeContainer(c FullContainerName) protocol.FullContainerName {
	var optionalDynID protocol.Optional[uint32]
	if c.DynamicContainerID != 0 {
		optionalDynID = protocol.Option(c.DynamicContainerID)
	}

	return protocol.FullContainerName{
		ContainerID:        c.ContainerID,
		DynamicContainerID: optionalDynID,
	}
}

func DowngradeContainer(c protocol.FullContainerName) FullContainerName {
	dynID, _ := c.DynamicContainerID.Value()
	return FullContainerName{
		ContainerID:        c.ContainerID,
		DynamicContainerID: dynID,
	}
}

func DowngradeStackRequestSlotInfo(s protocol.StackRequestSlotInfo) StackRequestSlotInfo {
	return StackRequestSlotInfo{
		Container:      DowngradeContainer(s.Container),
		Slot:           s.Slot,
		StackNetworkID: s.StackNetworkID,
	}
}

func UpgradeStackRequestSlotInfo(s StackRequestSlotInfo) protocol.StackRequestSlotInfo {
	return protocol.StackRequestSlotInfo{
		Container:      UpgradeContainer(s.Container),
		Slot:           s.Slot,
		StackNetworkID: s.StackNetworkID,
	}
}

func DowngradeItemStackActions(actions []protocol.StackRequestAction) []protocol.StackRequestAction {
	for index, action := range actions {
		switch action := action.(type) {
		case *protocol.TakeStackRequestAction:
			actions[index] = &TakeStackRequestAction{
				transferStackRequestAction: transferStackRequestAction{
					Count:       action.Count,
					Source:      DowngradeStackRequestSlotInfo(action.Source),
					Destination: DowngradeStackRequestSlotInfo(action.Destination),
				},
			}
		case *protocol.PlaceStackRequestAction:
			actions[index] = &PlaceStackRequestAction{
				transferStackRequestAction: transferStackRequestAction{
					Count:       action.Count,
					Source:      DowngradeStackRequestSlotInfo(action.Source),
					Destination: DowngradeStackRequestSlotInfo(action.Destination),
				},
			}
		case *protocol.SwapStackRequestAction:
			actions[index] = &SwapStackRequestAction{
				Source:      DowngradeStackRequestSlotInfo(action.Source),
				Destination: DowngradeStackRequestSlotInfo(action.Destination),
			}
		case *protocol.DropStackRequestAction:
			actions[index] = &DropStackRequestAction{
				Count:    action.Count,
				Source:   DowngradeStackRequestSlotInfo(action.Source),
				Randomly: action.Randomly,
			}
		case *protocol.DestroyStackRequestAction:
			actions[index] = &DestroyStackRequestAction{
				Count:  action.Count,
				Source: DowngradeStackRequestSlotInfo(action.Source),
			}
		case *protocol.ConsumeStackRequestAction:
			actions[index] = &ConsumeStackRequestAction{
				DestroyStackRequestAction: DestroyStackRequestAction{
					Count:  action.Count,
					Source: DowngradeStackRequestSlotInfo(action.Source),
				},
			}
		case *protocol.PlaceInContainerStackRequestAction:
			actions[index] = &PlaceInContainerStackRequestAction{
				transferStackRequestAction: transferStackRequestAction{
					Count:       action.Count,
					Source:      DowngradeStackRequestSlotInfo(action.Source),
					Destination: DowngradeStackRequestSlotInfo(action.Destination),
				},
			}
		case *protocol.TakeOutContainerStackRequestAction:
			actions[index] = &PlaceInContainerStackRequestAction{
				transferStackRequestAction: transferStackRequestAction{
					Count:       action.Count,
					Source:      DowngradeStackRequestSlotInfo(action.Source),
					Destination: DowngradeStackRequestSlotInfo(action.Destination),
				},
			}
		case *protocol.LabTableCombineStackRequestAction:
			actions[index] = &LabTableCombineStackRequestAction{}
		case *protocol.AutoCraftRecipeStackRequestAction:
			actions[index] = &AutoCraftRecipeStackRequestAction{
				RecipeNetworkID: action.RecipeNetworkID,
				TimesCrafted:    action.TimesCrafted,
				Ingredients:     action.Ingredients,
			}
		case *protocol.CraftCreativeStackRequestAction:
			actions[index] = &CraftCreativeStackRequestAction{
				CreativeItemNetworkID: action.CreativeItemNetworkID,
			}
		case *protocol.CraftRecipeOptionalStackRequestAction:
			actions[index] = &CraftRecipeOptionalStackRequestAction{
				RecipeNetworkID:   action.RecipeNetworkID,
				FilterStringIndex: action.FilterStringIndex,
			}
		case *protocol.CraftGrindstoneRecipeStackRequestAction:
			actions[index] = &CraftGrindstoneRecipeStackRequestAction{
				RecipeNetworkID: action.RecipeNetworkID,
				Cost:            action.Cost,
			}
		default:
			actions[index] = action
		}
	}

	return actions
}

func UpgradeItemStackActions(actions []protocol.StackRequestAction) []protocol.StackRequestAction {
	for index, action := range actions {
		var a protocol.StackRequestAction

		switch action := action.(type) {
		case *CraftCreativeStackRequestAction:
			a = &protocol.CraftCreativeStackRequestAction{
				CreativeItemNetworkID: action.CreativeItemNetworkID,
				NumberOfCrafts:        0,
			}
		case *TakeStackRequestAction:
			ra := &protocol.TakeStackRequestAction{}
			ra.Count = action.Count
			ra.Source = UpgradeStackRequestSlotInfo(action.Source)
			ra.Destination = UpgradeStackRequestSlotInfo(action.Destination)
			a = ra
		case *PlaceStackRequestAction:
			ra := &protocol.PlaceStackRequestAction{}
			ra.Count = action.Count
			ra.Source = UpgradeStackRequestSlotInfo(action.Source)
			ra.Destination = UpgradeStackRequestSlotInfo(action.Destination)
			a = ra
		case *SwapStackRequestAction:
			a = &protocol.SwapStackRequestAction{
				Source:      UpgradeStackRequestSlotInfo(action.Source),
				Destination: UpgradeStackRequestSlotInfo(action.Destination),
			}
		case *DropStackRequestAction:
			a = &protocol.DropStackRequestAction{
				Count:    action.Count,
				Source:   UpgradeStackRequestSlotInfo(action.Source),
				Randomly: action.Randomly,
			}
		case *DestroyStackRequestAction:
			a = &protocol.DestroyStackRequestAction{
				Count:  action.Count,
				Source: UpgradeStackRequestSlotInfo(action.Source),
			}
		case *ConsumeStackRequestAction:
			a = &protocol.ConsumeStackRequestAction{
				DestroyStackRequestAction: protocol.DestroyStackRequestAction{
					Count:  action.Count,
					Source: UpgradeStackRequestSlotInfo(action.Source),
				},
			}
		case *PlaceInContainerStackRequestAction:
			ra := &protocol.PlaceInContainerStackRequestAction{}
			ra.Count = action.Count
			ra.Source = UpgradeStackRequestSlotInfo(action.Source)
			ra.Destination = UpgradeStackRequestSlotInfo(action.Destination)
			a = ra
		case *TakeOutContainerStackRequestAction:
			ra := &protocol.PlaceInContainerStackRequestAction{}
			ra.Count = action.Count
			ra.Source = UpgradeStackRequestSlotInfo(action.Source)
			ra.Destination = UpgradeStackRequestSlotInfo(action.Destination)
			a = ra
		case *LabTableCombineStackRequestAction:
			a = &protocol.LabTableCombineStackRequestAction{}
		case *AutoCraftRecipeStackRequestAction:
			a = &protocol.AutoCraftRecipeStackRequestAction{
				RecipeNetworkID: action.RecipeNetworkID,
				NumberOfCrafts:  action.TimesCrafted,
				TimesCrafted:    action.TimesCrafted,
				Ingredients:     action.Ingredients,
			}
		case *CraftRecipeOptionalStackRequestAction:
			a = &protocol.CraftRecipeOptionalStackRequestAction{
				RecipeNetworkID:   action.RecipeNetworkID,
				FilterStringIndex: action.FilterStringIndex,
			}
		case *CraftGrindstoneRecipeStackRequestAction:
			a = &protocol.CraftGrindstoneRecipeStackRequestAction{
				RecipeNetworkID: action.RecipeNetworkID,
				NumberOfCrafts:  0,
				Cost:            action.Cost,
			}
		default:
			a = action
		}

		actions[index] = a
	}

	return actions
}

// LookupStackRequestAction looks up the StackRequestAction matching an ID.
func LookupStackRequestAction(id uint8, x *protocol.StackRequestAction) bool {
	switch id {
	case protocol.StackRequestActionTake:
		*x = &TakeStackRequestAction{}
	case protocol.StackRequestActionPlace:
		*x = &PlaceStackRequestAction{}
	case protocol.StackRequestActionSwap:
		*x = &SwapStackRequestAction{}
	case protocol.StackRequestActionDrop:
		*x = &DropStackRequestAction{}
	case protocol.StackRequestActionDestroy:
		*x = &DestroyStackRequestAction{}
	case protocol.StackRequestActionConsume:
		*x = &ConsumeStackRequestAction{}
	case protocol.StackRequestActionCreate:
		*x = &protocol.CreateStackRequestAction{}
	case protocol.StackRequestActionPlaceInContainer:
		*x = &PlaceInContainerStackRequestAction{}
	case protocol.StackRequestActionTakeOutContainer:
		*x = &TakeOutContainerStackRequestAction{}
	case protocol.StackRequestActionLabTableCombine:
		*x = &LabTableCombineStackRequestAction{}
	case protocol.StackRequestActionBeaconPayment:
		*x = &protocol.BeaconPaymentStackRequestAction{}
	case protocol.StackRequestActionMineBlock:
		*x = &protocol.MineBlockStackRequestAction{}
	case protocol.StackRequestActionCraftRecipe:
		*x = &protocol.CraftRecipeStackRequestAction{}
	case protocol.StackRequestActionCraftRecipeAuto:
		*x = &AutoCraftRecipeStackRequestAction{}
	case protocol.StackRequestActionCraftCreative:
		*x = &CraftCreativeStackRequestAction{}
	case protocol.StackRequestActionCraftRecipeOptional:
		*x = &CraftRecipeOptionalStackRequestAction{}
	case protocol.StackRequestActionCraftGrindstone:
		*x = &CraftGrindstoneRecipeStackRequestAction{}
	case protocol.StackRequestActionCraftLoom:
		*x = &protocol.CraftLoomRecipeStackRequestAction{}
	case protocol.StackRequestActionCraftNonImplementedDeprecated:
		*x = &protocol.CraftNonImplementedStackRequestAction{}
	case protocol.StackRequestActionCraftResultsDeprecated:
		*x = &protocol.CraftResultsDeprecatedStackRequestAction{}
	default:
		return false
	}
	return true
}

// LookupStackRequestActionType looks up the ID of a StackRequestAction.
func LookupStackRequestActionType(x protocol.StackRequestAction, id *uint8) bool {
	switch x.(type) {
	case *protocol.TakeStackRequestAction, *TakeStackRequestAction:
		*id = protocol.StackRequestActionTake
	case *protocol.PlaceStackRequestAction, *PlaceStackRequestAction:
		*id = protocol.StackRequestActionPlace
	case *protocol.SwapStackRequestAction, *SwapStackRequestAction:
		*id = protocol.StackRequestActionSwap
	case *protocol.DropStackRequestAction, *DropStackRequestAction:
		*id = protocol.StackRequestActionDrop
	case *protocol.DestroyStackRequestAction, *DestroyStackRequestAction:
		*id = protocol.StackRequestActionDestroy
	case *protocol.ConsumeStackRequestAction, *ConsumeStackRequestAction:
		*id = protocol.StackRequestActionConsume
	case *protocol.CreateStackRequestAction:
		*id = protocol.StackRequestActionCreate
	case *protocol.LabTableCombineStackRequestAction, *LabTableCombineStackRequestAction:
		*id = protocol.StackRequestActionLabTableCombine
	case *protocol.BeaconPaymentStackRequestAction:
		*id = protocol.StackRequestActionBeaconPayment
	case *protocol.MineBlockStackRequestAction:
		*id = protocol.StackRequestActionMineBlock
	case *protocol.CraftRecipeStackRequestAction:
		*id = protocol.StackRequestActionCraftRecipe
	case *protocol.AutoCraftRecipeStackRequestAction, *AutoCraftRecipeStackRequestAction:
		*id = protocol.StackRequestActionCraftRecipeAuto
	case *protocol.CraftCreativeStackRequestAction, *CraftCreativeStackRequestAction:
		*id = protocol.StackRequestActionCraftCreative
	case *protocol.CraftRecipeOptionalStackRequestAction, *CraftRecipeOptionalStackRequestAction:
		*id = protocol.StackRequestActionCraftRecipeOptional
	case *protocol.CraftGrindstoneRecipeStackRequestAction, *CraftGrindstoneRecipeStackRequestAction:
		*id = protocol.StackRequestActionCraftGrindstone
	case *protocol.CraftLoomRecipeStackRequestAction:
		*id = protocol.StackRequestActionCraftLoom
	case *protocol.CraftNonImplementedStackRequestAction:
		*id = protocol.StackRequestActionCraftNonImplementedDeprecated
	case *protocol.CraftResultsDeprecatedStackRequestAction:
		*id = protocol.StackRequestActionCraftResultsDeprecated
	default:
		return false
	}
	return true
}
