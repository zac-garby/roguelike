package main

import "math/rand"

const (
	wallBoxChance  = 0.0250
	otherBoxChance = 0.0125
	chestChance    = 0.0200
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

func (m *Map) width() int {
	return len(m.Tiles[0])
}

func (m *Map) height() int {
	return len(m.Tiles)
}

// postprocess processes a Map, adding in interesting tiles such as boxes,
// more defined walls, etc...
func (m *Map) postprocess() {
	for y := 0; y < m.height(); y++ {
		for x := 0; x < m.width(); x++ {
			t := m.Tiles[y][x]

			switch t {
			case TileOutside:
				n := m.neighbours(x, y, TileFloor, TileBox)
				if n > 0 {
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
		tx, ty := rand.Intn(m.width()), rand.Intn(m.height())
		if m.Tiles[ty][tx] == TileFloor && m.neighbours(tx, ty, TileFloor) == 8 {
			m.Tiles[ty][tx] = TileTrapdoor
			break
		}
	}
}

// neighbours gets the number of neighbours of a cell which are off a
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

		if cx < 0 || cy < 0 || cx >= m.width() || cy >= m.height() {
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
