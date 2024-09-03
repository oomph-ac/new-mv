package types

import "github.com/sandertv/gophertunnel/minecraft/protocol"

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
	if x.Status == protocol.ItemStackResponseStatusOK { // 0
		protocol.Slice(r, &x.ContainerInfo)
	}
}

// StackResponseContainerInfo holds information on what slots in a container have what item stack in them.
type StackResponseContainerInfo struct {
	protocol.StackResponseContainerInfo
}

// Marshal encodes/decodes a StackResponseContainerInfo.
func (x *StackResponseContainerInfo) Marshal(r protocol.IO) {
	r.Uint8(&x.Container.ContainerID)
	protocol.Slice(r, &x.SlotInfo)
}

// CraftRecipeStackRequestAction is sent by the client the moment it begins crafting an item. This is the
// first action sent, before the Consume and Create item stack request actions.
// This action is also sent when an item is enchanted. Enchanting should be treated mostly the same way as
// crafting, where the old item is consumed.
type CraftRecipeStackRequestAction struct {
	protocol.CraftRecipeStackRequestAction
}

// Marshal ...
func (a *CraftRecipeStackRequestAction) Marshal(r protocol.IO) {
	r.Varuint32(&a.RecipeNetworkID)
}

// AutoCraftRecipeStackRequestAction is sent by the client similarly to the CraftRecipeStackRequestAction. The
// only difference is that the recipe is automatically created and crafted by shift clicking the recipe book.
type AutoCraftRecipeStackRequestAction struct {
	protocol.AutoCraftRecipeStackRequestAction
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
	protocol.CraftCreativeStackRequestAction
}

// Marshal ...
func (a *CraftCreativeStackRequestAction) Marshal(r protocol.IO) {
	r.Varuint32(&a.CreativeItemNetworkID)
}

// CraftRecipeOptionalStackRequestAction is sent when using an anvil. When this action is sent, the
// FilterStrings field in the respective stack request is non-empty and contains the name of the item created
// using the anvil or cartography table.
type CraftRecipeOptionalStackRequestAction struct {
	protocol.CraftRecipeOptionalStackRequestAction
}

// Marshal ...
func (c *CraftRecipeOptionalStackRequestAction) Marshal(r protocol.IO) {
	r.Varuint32(&c.RecipeNetworkID)
	r.Int32(&c.FilterStringIndex)
}

// CraftGrindstoneRecipeStackRequestAction is sent when a grindstone recipe is crafted. It contains the RecipeNetworkID
// to identify the recipe crafted, and the cost for crafting the recipe.
type CraftGrindstoneRecipeStackRequestAction struct {
	protocol.CraftGrindstoneRecipeStackRequestAction
}

// Marshal ...
func (c *CraftGrindstoneRecipeStackRequestAction) Marshal(r protocol.IO) {
	r.Varuint32(&c.RecipeNetworkID)
	r.Varint32(&c.Cost)
}

// StackRequestSlotInfo holds information on a specific slot client-side.
type StackRequestSlotInfo struct {
	ContainerID byte
	// Slot is the index of the slot within the container with the ContainerID above.
	Slot byte
	// StackNetworkID is the unique stack ID that the client assumes to be present in this slot. The server
	// must check if these IDs match. If they do not match, servers should reject the stack request that the
	// action holding this info was in.
	StackNetworkID int32
}

// StackReqSlotInfo reads/writes a StackRequestSlotInfo x using IO r.
func StackReqSlotInfo(r protocol.IO, x *StackRequestSlotInfo) {
	r.Uint8(&x.ContainerID)
	r.Uint8(&x.Slot)
	r.Varint32(&x.StackNetworkID)
}

// transferStackRequestAction is the structure shared by StackRequestActions that transfer items from one
// slot into another.
type transferStackRequestAction struct {
	Count               byte
	Source, Destination StackRequestSlotInfo
}

// Marshal ...
func (a *transferStackRequestAction) Marshal(r protocol.IO) {
	r.Uint8(&a.Count)
	StackReqSlotInfo(r, &a.Source)
	StackReqSlotInfo(r, &a.Destination)
}

// TakeStackRequestAction is sent by the client to the server to take x amount of items from one slot in a
// container to the cursor.
type TakeStackRequestAction struct {
	transferStackRequestAction
}

// PlaceStackRequestAction is sent by the client to the server to place x amount of items from one slot into
// another slot, such as when shift clicking an item in the inventory to move it around or when moving an item
// in the cursor into a slot.
type PlaceStackRequestAction struct {
	transferStackRequestAction
}

// SwapStackRequestAction is sent by the client to swap the item in its cursor with an item present in another
// container. The two item stacks swap places.
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

// DropStackRequestAction is sent by the client when it drops an item out of the inventory when it has its
// inventory opened. This action is not sent when a player drops an item out of the hotbar using the Q button
// (or the equivalent on mobile). The InventoryTransaction packet is still used for that action, regardless of
// whether the item stack network IDs are used or not.
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

// DestroyStackRequestAction is sent by the client when it destroys an item in creative mode by moving it
// back into the creative inventory.
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

// ConsumeStackRequestAction is sent by the client when it uses an item to craft another item. The original
// item is 'consumed'.
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

// LookupStackRequestActionType looks up the ID of a StackRequestAction.
func LookupStackRequestActionType(x protocol.StackRequestAction, id *uint8) bool {
	switch x.(type) {
	case *TakeStackRequestAction:
		*id = protocol.StackRequestActionTake
	case *PlaceStackRequestAction:
		*id = protocol.StackRequestActionPlace
	case *SwapStackRequestAction:
		*id = protocol.StackRequestActionSwap
	case *DropStackRequestAction:
		*id = protocol.StackRequestActionDrop
	case *DestroyStackRequestAction:
		*id = protocol.StackRequestActionDestroy
	case *ConsumeStackRequestAction:
		*id = protocol.StackRequestActionConsume
	case *protocol.CreateStackRequestAction:
		*id = protocol.StackRequestActionCreate
	case *PlaceInContainerStackRequestAction:
		*id = protocol.StackRequestActionPlaceInContainer
	case *TakeOutContainerStackRequestAction:
		*id = protocol.StackRequestActionTakeOutContainer
	case *LabTableCombineStackRequestAction:
		*id = protocol.StackRequestActionLabTableCombine
	case *protocol.BeaconPaymentStackRequestAction:
		*id = protocol.StackRequestActionBeaconPayment
	case *protocol.MineBlockStackRequestAction:
		*id = protocol.StackRequestActionMineBlock
	case *CraftRecipeStackRequestAction:
		*id = protocol.StackRequestActionCraftRecipe
	case *AutoCraftRecipeStackRequestAction:
		*id = protocol.StackRequestActionCraftRecipeAuto
	case *CraftCreativeStackRequestAction:
		*id = protocol.StackRequestActionCraftCreative
	case *CraftRecipeOptionalStackRequestAction:
		*id = protocol.StackRequestActionCraftRecipeOptional
	case *CraftGrindstoneRecipeStackRequestAction:
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

// LookupStackRequestAction looks up the StackRequestAction matching an ID
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
		*x = &CraftRecipeStackRequestAction{}
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
