package main

import (
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/fatih/color"
	"github.com/inancgumus/screen"
	"golang.org/x/term"
)

type Leaf struct {
	X          int
	Y          int
	Charactere rune
	Speed      int
	Color      color.Attribute
}

func PrintAt(x, y int, char rune, printColor color.Attribute) {
	color.Set(printColor)
	fmt.Printf("\033[?25l\033[%d;%dH%c", y+1, x+1, char)
	color.Unset()
	//fmt.Print(x, ", ", y)i col
}
func generateLeaves() {

}
func main() {
	// defered later in the code  : terminalWidth, _, _ := term.GetSize(0)
	defer fmt.Printf("\033[?25h")
	var leaves []Leaf
	var terminalWidth int
	var terminalHeight int

	for {
		if w, h, err := term.GetSize(1); err == nil {
			terminalWidth, terminalHeight = w, h
			break
		}
		if w, h, err := term.GetSize(0); err == nil {
			terminalWidth, terminalHeight = w, h
			break
		}
		if w, h, err := term.GetSize(2); err == nil {
			terminalWidth, terminalHeight = w, h
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
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
		var randomColor color.Attribute
		randomColorNum := rand.IntN(3)
		switch randomColorNum {
		case 0:
			randomColor = color.FgRed
		case 1:
			randomColor = color.FgYellow
		case 2:
			randomColor = color.FgHiRed
		case 3:
			randomColor = color.FgHiYellow

		}
		randomSpeed := rand.IntN(5)
		if randomSpeed == 0 {
			randomSpeed++
		}
		leaves = append(leaves, Leaf{X: randomX, Y: randomY, Charactere: randomChar, Speed: randomSpeed, Color: randomColor})

	}

	for {

		screen.Clear()
		//terminalWidth, terminalHeight, _ := terminal.GetSize(0)
		for id := range leaves {
			leaves[id].Y = leaves[id].Y + leaves[id].Speed
			PrintAt(leaves[id].X, leaves[id].Y, leaves[id].Charactere, leaves[id].Color)
			if leaves[id].Y >= terminalHeight {
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
				var randomColor color.Attribute
				randomColorNum := rand.IntN(3)
				switch randomColorNum {
				case 0:
					randomColor = color.FgYellow
				case 1:
					randomColor = color.FgRed
				case 2:
					randomColor = color.FgHiYellow
				case 3:
					randomColor = color.FgHiRed
				}
				leaves[id].Y = randomY
				leaves[id].X = randomX
				leaves[id].Charactere = randomChar
				leaves[id].Speed = randomSpeed
				leaves[id].Color = randomColor
			}
		}
		time.Sleep(100000000)

	}
}
