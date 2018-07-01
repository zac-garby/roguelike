package main

import (
	"time"

	"github.com/Zac-Garby/roguelike/lib"
	"github.com/nsf/termbox-go"
)

var (
	level     *lib.Map
	player    *lib.Player
	ui        *lib.UI
	lastMove  = time.Now()
	moveDelay = 0.08
)

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}

	defer termbox.Close()
	termbox.SetOutputMode(termbox.Output256)
	termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse)

	lib.SpawnWorkers()
	level = lib.MakeMap()
	player = lib.NewPlayer(level)
	ui = &lib.UI{
		Player: player,
		Level:  level,
	}

	redraw()

mainloop:
	for {
		switch evt := termbox.PollEvent(); evt.Type {
		case termbox.EventKey:
			if evt.Key == termbox.KeyEsc {
				break mainloop
			}

			handleKey(evt.Key)
		}

		redraw()
	}
}

func handleKey(key termbox.Key) {
	if time.Now().Sub(lastMove).Seconds() > moveDelay {
		switch key {
		case termbox.KeyArrowLeft:
			player.Move(-1, 0, level)
			player.Direction = 3
		case termbox.KeyArrowRight:
			player.Move(1, 0, level)
			player.Direction = 1
		case termbox.KeyArrowUp:
			player.Move(0, -1, level)
			player.Direction = 0
		case termbox.KeyArrowDown:
			player.Move(0, 1, level)
			player.Direction = 2
		}

		lastMove = time.Now()
	}
}

func redraw() {
	level.Render(2, 1)
	player.Render(2, 1)
	ui.Render(100, 1)
	termbox.Flush()
}
