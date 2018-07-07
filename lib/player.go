package lib

import (
	"math/rand"
	"time"

	termbox "github.com/nsf/termbox-go"
)

// A Player is the user's player, and stores things such as position.
type Player struct {
	X, Y int

	Money      int
	Health     int
	Experience int
	Attack     int
	Defense    int
	Magic      int

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
				X:          x,
				Y:          y,
				Money:      0,
				Health:     100,
				Experience: 0,
				Attack:     1,
				Defense:    1,
				Magic:      1,
				Map:        m,
			}
		}
	}
}

// Move translates the player (dx, dy) units, but only if it will still
// be in a valid position.
func (p *Player) Move(dx, dy int, g *Game) {
	nx, ny := p.X+dx, p.Y+dy

	switch g.Level.At(nx, ny) {
	case TileBox, TileWall, TileOutside, TileMerchant:
		return

	case TileTrapdoor:
		prev := g.Level
		g.Level = MakeMap(g.Level.Depth + 1)

		if levelChangeConfirm(g) {
			g.Player.Map = g.Level

			for g.Level.At(g.Player.X, g.Player.Y) != TileFloor {
				g.Player.X = rand.Intn(g.Level.Width())
				g.Player.Y = rand.Intn(g.Level.Height())
			}
		} else {
			g.Level = prev
		}

		return

	case TileChest:
		p.Money += rand.Intn(100) + 50
		p.Experience += rand.Intn(10) + 5
		p.Map.Set(nx, ny, TileFloor)
	}

	p.X = nx
	p.Y = ny
}

// Render renders a Player to the terminal, assuming the top-left of the
// map is at (x, y)
func (p *Player) Render(x, y int) {
	ch := []rune{
		'▲',
		'▶',
		'▼',
		'◀',
	}[p.Direction]

	termbox.SetCell(x+p.X*2, y+p.Y, ch, termbox.ColorCyan, termbox.ColorDefault)
	termbox.SetCell(x+p.X*2+1, y+p.Y, ' ', termbox.ColorCyan, termbox.ColorDefault)
}

func levelChangeConfirm(g *Game) bool {
	stop := make(chan bool, 1)

	delayText(
		1, 0, time.Millisecond*10,
		`
Entering level ^B%d^!...

^r health:  %d
^g money:   %d
^y xp:      %d
^c attack:  %d
^w defense: %d
^m magic:   %d ^!

Press ^BRETURN^! to enter the next level
Press ^BESC^! to stay on the current level`,
		termbox.ColorDefault, termbox.ColorDefault, stop,
		g.Level.Depth, g.Player.Health,
		g.Player.Money, g.Player.Experience,
		g.Player.Attack, g.Player.Defense,
		g.Player.Magic)

	for {
		switch evt := termbox.PollEvent(); evt.Type {
		case termbox.EventKey:
			if evt.Key == termbox.KeyEnter {
				stop <- true
				return true
			}

			if evt.Key == termbox.KeyEsc {
				stop <- true
				return false
			}
		}
	}
}
