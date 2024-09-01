package latest

import (
	_ "embed"

	"github.com/oomph-ac/new-mv/internal/chunk"
	"github.com/oomph-ac/new-mv/mapping"
)

const (
	// BlockVersion is the version of blocks (states) of the game. This version is composed
	// of 4 bytes indicating a version, interpreted as a big endian int. The current version represents
	// 1.21.20.3 {1, 21, 20, 3}.
	BlockVersion int32 = (1 << 24) | (21 << 16) | (0 << 8) | 3
)

var (
	//go:embed block_states.nbt
	blockStateData []byte

	// blockMapping is the BlockMapping used for translating blocks between versions.
	blockMapping = mapping.NewBlockMapping(blockStateData)
	// NetworkPersistentEncoding is the Encoding used for sending a Chunk over network. It uses NBT, unlike NetworkEncoding.
	NetworkPersistentEncoding = chunk.NewNetworkPersistentEncoding(blockMapping, BlockVersion)
	// BlockPaletteEncoding is the paletteEncoding used for encoding a palette of block states encoded as NBT.
	BlockPaletteEncoding = chunk.NewBlockPaletteEncoding(blockMapping, BlockVersion)
)

func NewBlockMapping() *mapping.DefaultBlockMapping {
	return blockMapping
}
