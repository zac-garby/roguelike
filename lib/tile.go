package lib

import (
	"math/rand"

	"github.com/nsf/termbox-go"
)

// The set of tile types
const (
	_ int = iota

	TileFloor
	TileWall
	TileOutside
	TileBox
	TileChest
	TileTrapdoor
	TileMerchant
)

// A Tile is a single block in the world
type Tile interface {
	Render(x, y int)
	Passable() bool
	Description() string
	OnWalk(g *Game)
	OnInteract(g *Game)
	Type() int
}

type tileDefaults struct{}

func (t *tileDefaults) Passable() bool      { return true }
func (t *tileDefaults) Description() string { return "idk" }
func (t *tileDefaults) OnWalk(g *Game)      {}
func (t *tileDefaults) OnInteract(g *Game)  {}

type (
	// A FloorTile is a walkable tile
	FloorTile struct {
		*tileDefaults
	}

	// A WallTile is a tile which the player cannot move through
	WallTile struct {
		*tileDefaults
	}

	// An OutsideTile is a tile outside the walls, which the player can't get to
	OutsideTile struct {
		*tileDefaults
	}

	// A BoxTile is a tile which spawns randomly near walls and can be broken
	// by shooting it
	BoxTile struct {
		*tileDefaults
	}

	// A ChestTile is a tile which, when touched, will open and give the player
	// some money and xp
	ChestTile struct {
		*tileDefaults
		Open bool
	}

	// A TrapdoorTile is a tile which, when touched, transports the player to the
	// next level
	TrapdoorTile struct {
		*tileDefaults
	}

	// A MerchantTile is a tile which allows the player to buy and sell items
	MerchantTile struct {
		*tileDefaults
	}
)

// Render() definitions

// Render renders a tile to the terminal
func (f *FloorTile) Render(x, y int) {
	writeText(x, y, -1, "  ", termbox.ColorDefault, termbox.ColorDefault)
}

// Render renders a tile to the terminal
func (f *WallTile) Render(x, y int) {
	writeText(x, y, -1, "  ", termbox.ColorDefault, 0x10)
}

// Render renders a tile to the terminal
func (f *OutsideTile) Render(x, y int) {
	writeText(x, y, -1, "  ", termbox.ColorDefault, termbox.ColorWhite)
}

// Render renders a tile to the terminal
func (f *BoxTile) Render(x, y int) {
	writeText(x, y, -1, "▨ ", termbox.ColorYellow|termbox.AttrBold, termbox.ColorDefault)
}

// Render renders a tile to the terminal
func (f *ChestTile) Render(x, y int) {
	if f.Open {
		writeText(x, y, -1, "$ ", termbox.ColorRed|termbox.AttrBold, termbox.ColorDefault)
	} else {
		writeText(x, y, -1, "$ ", termbox.ColorGreen|termbox.AttrBold, termbox.ColorDefault)
	}
}

// Render renders a tile to the terminal
func (f *TrapdoorTile) Render(x, y int) {
	writeText(x, y, -1, "[]", 0x0d, termbox.ColorDefault)
}

// Render renders a tile to the terminal
func (f *MerchantTile) Render(x, y int) {
	writeText(x, y, -1, "M ", termbox.ColorMagenta|termbox.AttrBold, termbox.ColorDefault)
}

// Type() definitions

// Type gets the type of a tile
func (f *FloorTile) Type() int {
	return TileFloor
}

// Type gets the type of a tile
func (f *WallTile) Type() int {
	return TileWall
}

// Type gets the type of a tile
func (f *OutsideTile) Type() int {
	return TileOutside
}

// Type gets the type of a tile
func (f *BoxTile) Type() int {
	return TileBox
}

// Type gets the type of a tile
func (f *ChestTile) Type() int {
	return TileChest
}

// Type gets the type of a tile
func (f *TrapdoorTile) Type() int {
	return TileTrapdoor
}

// Type gets the type of a tile
func (f *MerchantTile) Type() int {
	return TileMerchant
}

// Description() definitions

// Description returns a human-readable description of a tile
func (f *FloorTile) Description() string {
	return "Stop looking at the floor, do something interesting"
}

// Description returns a human-readable description of a tile
func (f *WallTile) Description() string {
	return "Keeps you inside the map ;)"
}

// Description returns a human-readable description of a tile
func (f *OutsideTile) Description() string {
	return "Uh.. How did you get out here??"
}

// Description returns a human-readable description of a tile
func (f *BoxTile) Description() string {
	return "Shoot this, maybe it will drop something cool. Or maybe not, who knows?"
}

// Description returns a human-readable description of a tile
func (f *ChestTile) Description() string {
	if f.Open {
		return "Looks like you've already opened this. Too bad :("
	}

	return "Lucky you! You found a chest. Interact with it (press S) to get some money."
}

// Description returns a human-readable description of a tile
func (f *TrapdoorTile) Description() string {
	return "Walk on this to go to the next level. Make sure you've done everything you want to here!"
}

// Description returns a human-readable description of a tile
func (f *MerchantTile) Description() string {
	return "wNaT tO tRadDE sOmE StuFF?!11?1!!"
}

// Passable() definitions

// Passable returns true if the tile can be walked through, false otherwise
func (f *WallTile) Passable() bool { return false }

// Passable returns true if the tile can be walked through, false otherwise
func (f *BoxTile) Passable() bool { return false }

// Passable returns true if the tile can be walked through, false otherwise
func (f *ChestTile) Passable() bool { return false }

// Passable returns true if the tile can be walked through, false otherwise
func (f *MerchantTile) Passable() bool { return false }

// OnWalk() definitions

// OnWalk is a callback which is fired when the tile is stepped on by the player
func (f *TrapdoorTile) OnWalk(g *Game) {
	prev := g.Level
	g.Level = MakeMap(g.Level.Depth + 1)

	if levelChangeConfirm(g) {
		g.Player.Game.Level = g.Level

		for g.Level.At(g.Player.X, g.Player.Y).Type() != TileFloor {
			g.Player.X = rand.Intn(g.Level.Width())
			g.Player.Y = rand.Intn(g.Level.Height())
		}
	} else {
		g.Level = prev
	}
}

// OnInteract() definitions

// OnInteract is a callback which is fired when the tile interacted with by the player
func (f *ChestTile) OnInteract(g *Game) {
	if f.Open {
		return
	}

	g.Player.Money += rand.Intn(100) + 50
	g.Player.Experience += rand.Intn(10) + 5

	f.Open = true
}
