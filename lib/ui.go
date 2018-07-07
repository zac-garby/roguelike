package lib

import (
	"fmt"
	"time"

	"github.com/nsf/termbox-go"
)

// The UI stores information about various entities and can draw
// a nice UI to show info
type UI struct {
	Game           *Game
	Minibuf        string
	MinibufTimeout time.Time
}

// Render renders the UI to the screen, relative to the given coords.
func (u *UI) Render(x, y int) {
	fg, bg := termbox.ColorDefault, termbox.ColorDefault

	writeText(x, y+0, -1, "      x: %d", fg, bg, u.Game.Player.X)
	writeText(x, y+1, -1, "      y: %d", fg, bg, u.Game.Player.Y)
	writeText(x, y+2, -1, "  depth: %d", fg, bg, u.Game.Level.Depth)
	writeText(x, y+4, -1, " health: ^r%d%%^!", fg, bg, u.Game.Player.Health)
	writeText(x, y+5, -1, "  money: ^g£%d^!", fg, bg, u.Game.Player.Money)
	writeText(x, y+6, -1, "     xp: ^y%d^!", fg, bg, u.Game.Player.Experience)
	writeText(x, y+7, -1, " attack: ^c%d^!", fg, bg, u.Game.Player.Attack)
	writeText(x, y+8, -1, "defense: ^w%d^!", fg, bg, u.Game.Player.Defense)
	writeText(x, y+9, -1, "  magic: ^m%d^!", fg, bg, u.Game.Player.Magic)

	fg = 0x09
	writeText(x, y+12, -1, "^wESC^! to exit the game", fg, bg)
	writeText(x, y+13, -1, "^wQ^! to exit to the menu", fg, bg)
	writeText(x, y+14, -1, "^w▲▼◀▶^! to move", fg, bg)
	writeText(x, y+15, -1, "^wIJKL^! to turn on the spot", fg, bg)
	writeText(x, y+16, -1, "^wSPACE^! to shoot", fg, bg)
	writeText(x, y+17, -1, "^wS^! to interact with a tile", fg, bg)
	writeText(x, y+18, -1, "^wD^! to inspect a tile", fg, bg)

	if time.Now().After(u.MinibufTimeout) {
		u.Minibuf = ""
	}

	if len(u.Minibuf) > 0 {
		fg = termbox.ColorDefault
		w, _ := termbox.Size()
		writeText(x, y+21, w-3, u.Minibuf, fg, bg)
	}
}

// SetMinibuf sets the UI's mini-buffer for a certain amount of time.
func (u *UI) SetMinibuf(text string, duration time.Duration) {
	u.Minibuf = text
	u.MinibufTimeout = time.Now().Add(duration)
}

func writeText(sx, sy, wx int, text string, fg, bg termbox.Attribute, args ...interface{}) {
	dfg := fg
	x := sx
	y := sy

	str := []rune(fmt.Sprintf(text, args...))

	for i := 0; i < len(str); i++ {
		ch := rune(str[i])
		if ch == '\n' {
			x = sx
			y++
		} else if ch == '\t' {
			x += 4
			continue
		} else if ch == '^' {
			if i+1 >= len(str) || i+1 < 0 {
				continue
			}

			switch str[i+1] {
			case 'b':
				fg = termbox.ColorBlue
			case 'c':
				fg = termbox.ColorCyan
			case '!':
				fg = dfg
			case 'g':
				fg = termbox.ColorGreen
			case 'm':
				fg = termbox.ColorMagenta
			case 'r':
				fg = termbox.ColorRed
			case 'w':
				fg = termbox.ColorWhite
			case 'y':
				fg = termbox.ColorYellow
			case 'B':
				fg |= termbox.AttrBold
			}

			i++

			continue
		}

		for x > wx && wx >= 0 {
			x = sx
			y++
		}

		termbox.SetCell(x, y, ch, fg, bg)
		x++
	}
}

func delayText(x, y int, delay time.Duration, text string, fg, bg termbox.Attribute, stop chan bool, args ...interface{}) {
	go func() {
		str := fmt.Sprintf(text, args...)

		for i := 0; i < len(str); i++ {
			if stop != nil && len(stop) > 0 {
				break
			}

			termbox.Clear(fg, bg)
			writeText(x, y, -1, str[:i+1], fg, bg)
			termbox.Flush()
			time.Sleep(delay)
		}
	}()
}
