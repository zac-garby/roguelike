package main

import (
	"log"
	"os"
	"time"

	"github.com/Zac-Garby/roguelike/lib"
	"github.com/nsf/termbox-go"
)

var (
	game      *lib.Game
	moveDelay = 0.08
)

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	lf, err := os.OpenFile("roguelike.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening log file: %v", err)
	}
	defer lf.Close()

	log.SetOutput(lf)
	log.Println("starting game")
	defer func() {
		log.Println("closing game")
	}()

	termbox.SetOutputMode(termbox.Output256)
	termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse)

	lib.SpawnWorkers()

	game = &lib.Game{
		Level:    lib.MakeMap(1),
		LastMove: time.Now(),
	}

	game.Player = lib.NewPlayer(game)

	game.UI = &lib.UI{
		Game: game,
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
	if time.Now().Sub(game.LastMove).Seconds() > moveDelay {
		switch key {
		case termbox.KeyArrowLeft:
			game.Player.Move(-1, 0)
			game.Player.Direction = 3
		case termbox.KeyArrowRight:
			game.Player.Move(1, 0)
			game.Player.Direction = 1
		case termbox.KeyArrowUp:
			game.Player.Move(0, -1)
			game.Player.Direction = 0
		case termbox.KeyArrowDown:
			game.Player.Move(0, 1)
			game.Player.Direction = 2
		}

		game.LastMove = time.Now()
	}
}

func redraw() {
	game.Render()
	termbox.Flush()
}
