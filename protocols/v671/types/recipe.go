package types

import "github.com/sandertv/gophertunnel/minecraft/protocol"

// ShapelessRecipe is a recipe that has no particular shape. Its functionality is shared with the
// RecipeShulkerBox and RecipeShapelessChemistry types.
type ShapelessRecipe struct {
	protocol.ShapelessRecipe
}

func (recipe *ShapelessRecipe) Marshal(w protocol.IO) {
	marshalShapeless(w, recipe)
}

func (recipe *ShapelessRecipe) Unmarshal(r protocol.IO) {
	marshalShapeless(r, recipe)
}

// ShulkerBoxRecipe is a shapeless recipe made specifically for shulker box crafting, so that they don't lose
// their user data when dyeing a shulker box.
type ShulkerBoxRecipe struct {
	ShapelessRecipe
}

// ShapelessChemistryRecipe is a recipe specifically made for chemistry related features, which exist only in
// the Education Edition. They function the same as shapeless recipes do.
type ShapelessChemistryRecipe struct {
	ShapelessRecipe
}

// ShapedRecipe is a recipe that has a specific shape that must be used to craft the output of the recipe.
// Trying to craft the item in any other shape will not work. The ShapedRecipe is of the same structure as the
// ShapedChemistryRecipe.
type ShapedRecipe struct {
	protocol.ShapedRecipe
}

type ShapedChemistryRecipe struct {
	ShapedRecipe
}

func (recipe *ShapedRecipe) Marshal(w protocol.IO) {
	marshalShaped(w, recipe)
}

func (recipe *ShapedRecipe) Unmarshal(r protocol.IO) {
	marshalShaped(r, recipe)
}

func marshalShapeless(r protocol.IO, recipe *ShapelessRecipe) {
	r.String(&recipe.RecipeID)
	protocol.FuncSlice(r, &recipe.Input, r.ItemDescriptorCount)
	protocol.FuncSlice(r, &recipe.Output, r.Item)
	r.UUID(&recipe.UUID)
	r.String(&recipe.Block)
	r.Varint32(&recipe.Priority)
	r.Varuint32(&recipe.RecipeNetworkID)
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
	r.Bool(&recipe.AssumeSymmetry)
	r.Varuint32(&recipe.RecipeNetworkID)
}

// LookupRecipe looks up the Recipe for a recipe type. False is returned if not
// found.
func LookupRecipe(recipeType int32, x *protocol.Recipe) bool {
	switch recipeType {
	case protocol.RecipeShapeless:
		*x = &ShapelessRecipe{}
	case protocol.RecipeShaped:
		*x = &ShapedRecipe{}
	case protocol.RecipeFurnace:
		*x = &protocol.FurnaceRecipe{}
	case protocol.RecipeFurnaceData:
		*x = &protocol.FurnaceDataRecipe{}
	case protocol.RecipeMulti:
		*x = &protocol.MultiRecipe{}
	case protocol.RecipeShulkerBox:
		*x = &ShulkerBoxRecipe{}
	case protocol.RecipeShapelessChemistry:
		*x = &ShapelessChemistryRecipe{}
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
	case *ShapelessRecipe:
		*recipeType = protocol.RecipeShapeless
	case *ShapedRecipe:
		*recipeType = protocol.RecipeShaped
	case *protocol.FurnaceRecipe:
		*recipeType = protocol.RecipeFurnace
	case *protocol.FurnaceDataRecipe:
		*recipeType = protocol.RecipeFurnaceData
	case *protocol.MultiRecipe:
		*recipeType = protocol.RecipeMulti
	case *ShulkerBoxRecipe:
		*recipeType = protocol.RecipeShulkerBox
	case *ShapelessChemistryRecipe:
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

func DowngradeRecipes(recipes []protocol.Recipe) []protocol.Recipe {
	for index, r := range recipes {
		switch r := r.(type) {
		case *protocol.ShapelessRecipe:
			recipes[index] = &ShapelessRecipe{ShapelessRecipe: *r}
		case *protocol.ShapedRecipe:
			recipes[index] = &ShapedRecipe{ShapedRecipe: *r}
		case *protocol.ShulkerBoxRecipe:
			recipes[index] = &ShulkerBoxRecipe{ShapelessRecipe: ShapelessRecipe{ShapelessRecipe: r.ShapelessRecipe}}
		case *protocol.ShapelessChemistryRecipe:
			recipes[index] = &ShulkerBoxRecipe{ShapelessRecipe: ShapelessRecipe{ShapelessRecipe: r.ShapelessRecipe}}
		case *protocol.ShapedChemistryRecipe:
			recipes[index] = &ShapedChemistryRecipe{ShapedRecipe: ShapedRecipe{ShapedRecipe: r.ShapedRecipe}}
		}
	}

	return recipes
}
