package main

import (
	"fmt"

	"github.com/Zac-Garby/roguelike/lib"
)

func main() {
	lib.SpawnWorkers()
	m := lib.MakeMap()
	m.Postprocess()

	for _, row := range m.Tiles {
		for _, tile := range row {
			s := "  "

			switch tile {
			case lib.TileFloor:
				s = "  "
			case lib.TileWall:
				s = "\x1b[47m  \x1b[0m"
			case lib.TileOutside:
				s = "\x1b[107m  \x1b[0m"
			case lib.TileBox:
				s = "\x1b[33m[]"
			case lib.TileChest:
				s = "\x1b[92m$ "
			case lib.TileTrapdoor:
				s = "\x1b[94m()"
			}

			fmt.Print(s)
		}

		fmt.Println()
	}
}
