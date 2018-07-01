package lib

import (
	"math/rand"
)

// A Player is the user's player, and stores things such as position.
type Player struct {
	X, Y      int
	Money     int
	Map       *Map
	Direction int // 0: top, 1: right, 2: bottom, 3: left
}

// NewPlayer creates a new player, choosing random valid coordinates.
func NewPlayer(m *Map) *Player {
	for {
		var (
			x = rand.Intn(m.Width())
			y = rand.Intn(m.Height())
		)

		if m.At(x, y) == TileFloor {
			return &Player{
				X:     x,
				Y:     y,
				Money: 0,
				Map:   m,
			}
		}
	}
}

// Move translates the player (dx, dy) units, but only if it will still
// be in a valid position.
func (p *Player) Move(dx, dy int, m *Map) {
	nx, ny := p.X+dx, p.Y+dy

	switch m.At(nx, ny) {
	case TileBox, TileWall, TileOutside:
		return

	case TileChest:
		p.Money += rand.Intn(100) + 50
		p.Map.Set(nx, ny, TileFloor)
	}

	p.X = nx
	p.Y = ny
}
