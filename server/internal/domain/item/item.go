// Package item defines the core domain entities for in-game items and inventory.
// This package is PURE and must NOT import any infrastructure packages.
package item

// ItemType represents the kind of item.
type ItemType string

const (
	ItemRice        ItemType = "RICE"         // Free, basic sustenance
	ItemWater       ItemType = "WATER"        // Free, basic hydration
	ItemSushi       ItemType = "SUSHI"        // Premium food, costs Hype
	ItemCigarette   ItemType = "CIGARETTE"    // Contraband, restores sanity
	ItemPhone       ItemType = "PHONE"        // Contraband, massive advantage
	ItemElixir      ItemType = "ELIXIR"       // Mystic item, placebo effect
	ItemDragonBlood ItemType = "DRAGON_BLOOD" // Mystic item, tradeable
)

// ItemStack represents a quantity of a specific item type.
type ItemStack struct {
	Type     ItemType `json:"type"`
	Quantity int      `json:"quantity"`
}

// ItemDefinition provides metadata about an item type.
type ItemDefinition struct {
	Name        string
	Description string
	BaseValue   float64 // Value in Hype for trading
	IsFood      bool
	Nutrition   int // Hunger restoration
	Hydration   int // Thirst restoration
	SanityMod   int // Sanity restoration/drain
}

// Registry contains all known items and their properties.
var Registry = map[ItemType]ItemDefinition{
	ItemRice: {
		Name:        "Ración de Arroz",
		Description: "Comida básica e insípida. Quita el hambre, pero deprime.",
		BaseValue:   0,
		IsFood:      true,
		Nutrition:   20,
		Hydration:   0,
		SanityMod:   -2,
	},
	ItemWater: {
		Name:        "Agua del Grifo",
		Description: "Agua tibia de la celda.",
		BaseValue:   0,
		IsFood:      false,
		Nutrition:   0,
		Hydration:   25,
		SanityMod:   0,
	},
	ItemSushi: {
		Name:        "Bandeja de Sushi",
		Description: "Un lujo comprado con el bote. Sienta de maravilla.",
		BaseValue:   50.0,
		IsFood:      true,
		Nutrition:   40,
		Hydration:   0,
		SanityMod:   15,
	},
	ItemCigarette: {
		Name:        "Cigarrillo de Contrabando",
		Description: "Relaja los nervios, pero está prohibido.",
		BaseValue:   20.0,
		IsFood:      false,
		Nutrition:   0,
		Hydration:   -5,
		SanityMod:   25,
	},
	ItemPhone: {
		Name:        "Teléfono Móvil Oculto",
		Description: "Comunicación con el exterior. Objeto de máximo contrabando.",
		BaseValue:   500.0,
		IsFood:      false,
		SanityMod:   50,
	},
	ItemElixir: {
		Name:        "Elixir de Tartaria",
		Description: "Agua mezclada con tierra. El místico jura que cura todo.",
		BaseValue:   100.0, // High perceived value
		IsFood:      false,
		Hydration:   5,
		SanityMod:   5, // Placebo
	},
}

// GetItem returns the definition for an item type.
func GetItem(t ItemType) (ItemDefinition, bool) {
	def, ok := Registry[t]
	return def, ok
}
