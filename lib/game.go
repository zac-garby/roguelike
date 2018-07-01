package lib

import "time"

// The Game stores the game state so it can be easily passed around.
type Game struct {
	Level    *Map
	Player   *Player
	UI       *UI
	LastMove time.Time
}
