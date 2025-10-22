package main

import (
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/inancgumus/screen"
)

type Leaf struct {
	X          int
	Y          int
	Charactere rune
}

func PrintAt(x, y int, char rune) {
	fmt.Printf("\033[?25l\033[%d;%dH%c", y+1, x+1, char)
	//fmt.Print(x, ", ", y)
}
func main() {
	leaf1 := Leaf{
		X:          40,
		Y:          0,
		Charactere: '0',
	}
	fmt.Println("bob")
	for {
		screen.Clear()
		leaf1.Y++
		PrintAt(leaf1.X, leaf1.Y, leaf1.Charactere)
		time.Sleep(100000000)
		if leaf1.Y == 50 {
			leaf1.Y = 0
			leaf1.X = rand.IntN(100)
			randomChar := rand.IntN(4)
			switch randomChar {
			case 0:
				leaf1.Charactere = '0'
			case 1:
				leaf1.Charactere = '*'
			case 2:
				leaf1.Charactere = 'o'
			case 3:
				leaf1.Charactere = '¤'
			default:
				leaf1.Charactere = '~'
			}
		}
	}

}
