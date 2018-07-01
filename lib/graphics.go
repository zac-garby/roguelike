package lib

import (
	"github.com/nsf/termbox-go"
)

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

// Render renders a Player to the terminal, assuming the top-left of the
// map is at (x, y)
func (p *Player) Render(x, y int) {
	termbox.SetCell(x+p.X*2, y+p.Y, '#', termbox.ColorCyan, termbox.ColorDefault)
	termbox.SetCell(x+p.X*2+1, y+p.Y, ' ', termbox.ColorCyan, termbox.ColorDefault)
}
