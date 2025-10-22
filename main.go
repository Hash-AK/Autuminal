package main

import (
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/inancgumus/screen"
	"golang.org/x/crypto/ssh/terminal"
)

type Leaf struct {
	X          int
	Y          int
	Charactere rune
	Speed      int
}

func PrintAt(x, y int, char rune) {
	fmt.Printf("\033[?25l\033[%d;%dH%c", y+1, x+1, char)
	//fmt.Print(x, ", ", y)
}
func main() {
	terminalWidth, _, _ := terminal.GetSize(0)

	defer fmt.Printf("\033[?25h")
	var leaves []Leaf
	for count := 0; count <= 20; count++ {
		randomX := rand.IntN(terminalWidth)
		randomY := 0
		randomCharNum := rand.IntN(5)
		var randomChar rune
		switch randomCharNum {
		case 0:
			randomChar = '0'
		case 1:
			randomChar = '*'
		case 2:
			randomChar = 'o'
		case 3:
			randomChar = '¤'
		case 4:
			randomChar = '`'
		default:
			randomChar = '~'
		}
		randomSpeed := rand.IntN(5)
		if randomSpeed == 0 {
			randomSpeed++
		}
		leaves = append(leaves, Leaf{X: randomX, Y: randomY, Charactere: randomChar, Speed: randomSpeed})

	}
	for {
		screen.Clear()
		//terminalWidth, terminalHeight, _ := terminal.GetSize(0)
		for id := range leaves {
			leaves[id].Y = leaves[id].Y + leaves[id].Speed
			PrintAt(leaves[id].X, leaves[id].Y, leaves[id].Charactere)

		}
		time.Sleep(100000000)

	}
}
