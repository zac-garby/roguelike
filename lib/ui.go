package lib

import (
	"fmt"

	"github.com/nsf/termbox-go"
)

// The UI stores information about various entities and can draw
// a nice UI to show info
type UI struct {
	Game *Game
}

// Render renders the UI to the screen, relative to the given coords.
func (u *UI) Render(x, y int) {
	fg, bg := termbox.ColorDefault, termbox.ColorDefault

	writeText(x, y+0, "      x: %d", fg, bg, u.Game.Player.X)
	writeText(x, y+1, "      y: %d", fg, bg, u.Game.Player.Y)
	writeText(x, y+2, "  depth: %d", fg, bg, u.Game.Level.Depth)
	writeText(x, y+4, " health: %d", termbox.ColorRed, bg, u.Game.Player.Health)
	writeText(x, y+5, "  money: %d", termbox.ColorGreen, bg, u.Game.Player.Money)
	writeText(x, y+6, "     xp: %d", termbox.ColorYellow, bg, u.Game.Player.Experience)
	writeText(x, y+7, " attack: %d", termbox.ColorCyan, bg, u.Game.Player.Attack)
	writeText(x, y+8, "defense: %d", termbox.ColorWhite, bg, u.Game.Player.Defense)
	writeText(x, y+9, "  magic: %d", termbox.ColorMagenta, bg, u.Game.Player.Magic)

	fg = 0x09
	writeText(x, y+12, "use ESC to exit game", fg, bg)
	writeText(x, y+13, "use SPACE to exit to menu", fg, bg)
	writeText(x, y+14, "use ARROWS to move", fg, bg)
}

func writeText(x, y int, text string, fg, bg termbox.Attribute, args ...interface{}) {
	str := fmt.Sprintf(text, args...)

	for i := 0; i < len(str); i++ {
		ch := rune(str[i])
		termbox.SetCell(x+i, y, ch, fg, bg)
	}
}
