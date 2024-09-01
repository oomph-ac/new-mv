package latest

import (
	_ "embed"

	"github.com/oomph-ac/new-mv/mapping"
)

const ItemVersion = 191

var (
	//go:embed required_item_list.json
	requiredItemList []byte
	//go:embed item_runtime_ids.nbt
	itemRuntimeIDData []byte
)

func NewItemMapping(direct bool) mapping.Item {
	return mapping.NewItemMapping(itemRuntimeIDData, requiredItemList, ItemVersion, direct)
}
