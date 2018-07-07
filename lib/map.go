package lib

import (
	"math/rand"
)

var (
	wallBoxChance = 0.0350
	chestChance   = 0.0300
	numMerchants  = 3
)

// A Map represents a level in the game.
type Map struct {
	// Depth determines the level of this map in the game.
	Depth int

	// Tiles stores the tiles in a 2d matrix.
	Tiles [][]Tile
}

// Width returns the width of the map
func (m *Map) Width() int {
	return len(m.Tiles[0])
}

// Height returns the height of the map
func (m *Map) Height() int {
	return len(m.Tiles)
}

// At returns the tile at (x, y)
func (m *Map) At(x, y int) Tile {
	return m.Tiles[y][x]
}

// Set sets the tile at (x, y)
func (m *Map) Set(x, y int, t Tile) {
	m.Tiles[y][x] = t
}

// Postprocess processes a Map, adding in interesting tiles such as boxes,
// more defined walls, etc...
func (m *Map) Postprocess() {
	for y := 0; y < m.Height(); y++ {
		for x := 0; x < m.Width(); x++ {
			t := m.At(x, y)

			switch t.Type() {
			case TileOutside:
				if m.neighbours(x, y, TileFloor, TileBox) > 0 {
					m.Set(x, y, &WallTile{})
				}

			case TileFloor:
				r := rand.Float64()

				if m.neighbours(x, y, TileOutside, TileWall, TileBox) > 1 && r < wallBoxChance {
					m.Set(x, y, &BoxTile{})
				} else if m.neighbours(x, y, TileOutside, TileWall) == 0 && r < chestChance {
					m.Set(x, y, &ChestTile{
						Open: false,
					})
				}
			}
		}
	}

	for {
		tx, ty := rand.Intn(m.Width()), rand.Intn(m.Height())
		if m.At(tx, ty).Type() == TileFloor && m.neighbours(tx, ty, TileFloor) == 8 {
			m.Set(tx, ty, &TrapdoorTile{})
			break
		}
	}

	for i := 0; i < numMerchants; i++ {
		for {
			tx, ty := rand.Intn(m.Width()), rand.Intn(m.Height())
			if m.At(tx, ty).Type() == TileFloor && m.neighbours(tx, ty, TileFloor) == 8 {
				m.Set(tx, ty, &MerchantTile{})
				break
			}
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
			if m.At(cx, cy).Type() == tile {
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
			tile.Render(x+j*2, y+i)
		}
	}
}
