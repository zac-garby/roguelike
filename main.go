package main

import (
	"fmt"
)

func main() {
	spawnWorkers()
	m := makeMap()
	m.postprocess()

	for _, row := range m.Tiles {
		for _, tile := range row {
			s := "  "

			switch tile {
			case TileFloor:
				s = "  "
			case TileWall:
				s = "\x1b[47m  \x1b[0m"
			case TileOutside:
				s = "  "
			case TileBox:
				s = "\x1b[33m[]"
			case TileChest:
				s = "\x1b[92m$$"
			case TileTrapdoor:
				s = "\x1b[94m()"
			}

			fmt.Print(s)
		}

		fmt.Println()
	}
}
