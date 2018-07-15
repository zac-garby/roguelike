package lib

import (
	"github.com/nsf/termbox-go"
)

type (
	// An Item is an object which the player can have in his inventory. An
	// item has a 'quality' which is calculated by analysing its attributes
	// like type and rarity.
	Item interface {
		Type() ItemType
		Rarity() ItemRarity
		Modifiers() []ItemModifier
		Quality() int
	}

	// ItemKind is the type of a type, e.g. "weapon", "food"
	ItemKind struct {
		Name    string
		Quality int
	}

	// ItemType is the type of an item, e.g. "sword"
	ItemType struct {
		Name    string
		Kind    ItemKind
		Quality int // note this is as well as the quality for the item kind
	}

	// ItemRarity indicates how rare an item is
	ItemRarity struct {
		Name    string
		Fg      termbox.Attribute
		Bg      termbox.Attribute
		Quality int
	}

	// ItemModifier modifies the item in some way, e.g. "broken"
	ItemModifier struct {
		Name    string
		Quality int

		// Modify is called when the modifier is attached to an object.
		Modify func(i *Item, g *Game)

		// OnUse is used differently for different kinds of item. For example,
		// a weapon will call this when it's attacked with.
		OnUse func(i *Item, g *Game, data map[string]interface{})

		OnTick func(i *Item, g *Game)
	}
)
