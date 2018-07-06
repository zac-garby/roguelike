package lib

import "time"

// The Game stores the game state so it can be easily passed around.
type Game struct {
	Level    *Map
	Player   *Player
	UI       *UI
	LastMove time.Time
}

// Render renders the game to termbox
func (g *Game) Render() {
	g.Level.Render(2, 1)
	g.Player.Render(2, 1)
	g.UI.Render(100, 1)
}
