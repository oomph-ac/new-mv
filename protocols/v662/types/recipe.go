package types

import (
	v671types "github.com/oomph-ac/new-mv/protocols/v671/types"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

type ShapedRecipe struct {
	v671types.ShapedRecipe
}

func (recipe *ShapedRecipe) Marshal(w protocol.IO) {
	marshalShaped(w, recipe)
}

func (recipe *ShapedRecipe) Unmarshal(r protocol.IO) {
	marshalShaped(r, recipe)
}

type ShapedChemistryRecipe struct {
	ShapedRecipe
}

func marshalShaped(r protocol.IO, recipe *ShapedRecipe) {
	r.String(&recipe.RecipeID)
	r.Varint32(&recipe.Width)
	r.Varint32(&recipe.Height)
	protocol.FuncSliceOfLen(r, uint32(recipe.Width*recipe.Height), &recipe.Input, r.ItemDescriptorCount)
	protocol.FuncSlice(r, &recipe.Output, r.Item)
	r.UUID(&recipe.UUID)
	r.String(&recipe.Block)
	r.Varint32(&recipe.Priority)
	r.Varuint32(&recipe.RecipeNetworkID)
}

// LookupRecipe looks up the Recipe for a recipe type. False is returned if not
// found.
func LookupRecipe(recipeType int32, x *protocol.Recipe) bool {
	switch recipeType {
	case protocol.RecipeShapeless:
		*x = &v671types.ShapelessRecipe{}
	case protocol.RecipeShaped:
		*x = &ShapedRecipe{}
	case protocol.RecipeFurnace:
		*x = &protocol.FurnaceRecipe{}
	case protocol.RecipeFurnaceData:
		*x = &protocol.FurnaceDataRecipe{}
	case protocol.RecipeMulti:
		*x = &protocol.MultiRecipe{}
	case protocol.RecipeShulkerBox:
		*x = &v671types.ShulkerBoxRecipe{}
	case protocol.RecipeShapelessChemistry:
		*x = &v671types.ShapelessChemistryRecipe{}
	case protocol.RecipeShapedChemistry:
		*x = &ShapedChemistryRecipe{}
	case protocol.RecipeSmithingTransform:
		*x = &protocol.SmithingTransformRecipe{}
	case protocol.RecipeSmithingTrim:
		*x = &protocol.SmithingTrimRecipe{}
	default:
		return false
	}
	return true
}

// LookupRecipeType looks up the recipe type for a Recipe. False is returned if
// none was found.
func LookupRecipeType(x protocol.Recipe, recipeType *int32) bool {
	switch x.(type) {
	case *v671types.ShapelessRecipe:
		*recipeType = protocol.RecipeShapeless
	case *ShapedRecipe:
		*recipeType = protocol.RecipeShaped
	case *protocol.FurnaceRecipe:
		*recipeType = protocol.RecipeFurnace
	case *protocol.FurnaceDataRecipe:
		*recipeType = protocol.RecipeFurnaceData
	case *protocol.MultiRecipe:
		*recipeType = protocol.RecipeMulti
	case *v671types.ShulkerBoxRecipe:
		*recipeType = protocol.RecipeShulkerBox
	case *v671types.ShapelessChemistryRecipe:
		*recipeType = protocol.RecipeShapelessChemistry
	case *ShapedChemistryRecipe:
		*recipeType = protocol.RecipeShapedChemistry
	case *protocol.SmithingTransformRecipe:
		*recipeType = protocol.RecipeSmithingTransform
	case *protocol.SmithingTrimRecipe:
		*recipeType = protocol.RecipeSmithingTrim
	default:
		return false
	}
	return true
}
