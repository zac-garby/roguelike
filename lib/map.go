package lib

import (
	"math/rand"

	termbox "github.com/nsf/termbox-go"
)

var (
	wallBoxChance  = 0.0350
	otherBoxChance = 0.0000
	chestChance    = 0.0300
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
)

// A Map represents a level in the game.
type Map struct {
	// Depth determines the level of this map in the game.
	Depth int

	// Tiles stores the tiles in a 2d matrix.
	Tiles [][]int
}

// Width returns the width of the map
func (m *Map) Width() int {
	return len(m.Tiles[0])
}

// Height returns the height of the map
func (m *Map) Height() int {
	return len(m.Tiles)
}

// At returns the tile type at (x, y)
func (m *Map) At(x, y int) int {
	return m.Tiles[y][x]
}

// Set sets the tile type at (x, y)
func (m *Map) Set(x, y, t int) {
	m.Tiles[y][x] = t
}

// Postprocess processes a Map, adding in interesting tiles such as boxes,
// more defined walls, etc...
func (m *Map) Postprocess() {
	for y := 0; y < m.Height(); y++ {
		for x := 0; x < m.Width(); x++ {
			t := m.Tiles[y][x]

			switch t {
			case TileOutside:
				if m.neighbours(x, y, TileFloor, TileBox) > 0 {
					m.Tiles[y][x] = TileWall
				}

			case TileFloor:
				r := rand.Float64()

				if m.neighbours(x, y, TileOutside, TileWall, TileBox) > 1 && r < wallBoxChance {
					m.Tiles[y][x] = TileBox
				} else if m.neighbours(x, y, TileOutside, TileWall) == 0 && r < chestChance {
					m.Tiles[y][x] = TileChest
				} else if r < otherBoxChance {
					m.Tiles[y][x] = TileBox
				}
			}
		}
	}

	for {
		tx, ty := rand.Intn(m.Width()), rand.Intn(m.Height())
		if m.Tiles[ty][tx] == TileFloor && m.neighbours(tx, ty, TileFloor) == 8 {
			m.Tiles[ty][tx] = TileTrapdoor
			break
		}
	}
}

// neighbours gets the number of neighbours of a cell which are of a
// certain type.
func (m *Map) neighbours(x, y int, types ...int) int {
	coords := [][]int{
		{x - 1, y - 1}, {x, y - 1}, {x + 1, y - 1},
		{x - 1, y}, {x + 1, y},
		{x - 1, y + 1}, {x, y + 1}, {x + 1, y + 1},
	}

	count := 0

	for _, coord := range coords {
		cx, cy := coord[0], coord[1]

		if cx < 0 || cy < 0 || cx >= m.Width() || cy >= m.Height() {
			continue
		}

		for _, tile := range types {
			if m.Tiles[cy][cx] == tile {
				count++
			}
		}
	}

	return count
}

// Render renders a Map instance to the terminal at the given
// coordinates
func (m *Map) Render(x, y int) {
	for i, row := range m.Tiles {
		for j, tile := range row {
			var (
				s  = "  "
				fg = termbox.ColorDefault
				bg = termbox.ColorDefault
			)

			switch tile {
			case TileFloor:

			case TileWall:
				bg = 0x10
			case TileOutside:
				bg = termbox.ColorWhite
			case TileBox:
				s = "[]"
				fg = termbox.ColorYellow | termbox.AttrBold
			case TileChest:
				s = "$ "
				fg = termbox.ColorGreen | termbox.AttrBold
			case TileTrapdoor:
				s = "()"
				fg = 0x0d
			}

			termbox.SetCell(x+j*2, y+i, rune(s[0]), fg, bg)
			termbox.SetCell(x+1+j*2, y+i, rune(s[1]), fg, bg)
		}
	}
}
