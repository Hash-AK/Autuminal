package main

import (
	"fmt"
	"log"
	"math/rand/v2"
	"os"
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

const FgBrown = "\033[38;5;94m"
const ColorReset = "\033[0m"

func PrintAtColor(x, y int, char rune, colorCode string) {
	fmt.Printf("\033[?25l\033[%d;%dH%s%c%s", y+1, x+1, colorCode, char, ColorReset)
}
func generateLeaves() {

}

func main() {
	// defered later in the code  : terminalWidth, _, _ := term.GetSize(0)
	defer fmt.Printf("\033[?25h")

	var leaves []Leaf
	var terminalWidth int
	var terminalHeight int
	var reservedHeight int
	var currentJournalLine string
	var numberOfLine = 1
	var textBoxWidth int
	var textBoxBorderWidth int
	oldState, err := term.MakeRaw(0)
	if err != nil {
		panic(err)
	}
	defer term.Restore(0, oldState)

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
	inputChan := make(chan string)
	doneChan := make(chan bool)
	go func() {

		buffer := make([]byte, 1)

		for {
			os.Stdin.Read(buffer)
			if buffer[0] == 3 {
				doneChan <- true
				return
			}
			if buffer[0] == 13 {
				inputChan <- currentJournalLine
				currentJournalLine = ""
				numberOfLine = 1

			} else if buffer[0] == 8 || buffer[0] == 127 {
				if len(currentJournalLine) > 0 {
					currentJournalLine = currentJournalLine[:len(currentJournalLine)-1]
					numberOfLine = len(currentJournalLine) / textBoxWidth
					if len(currentJournalLine)%textBoxWidth != 0 {
						numberOfLine++
					}
					if numberOfLine == 0 {
						numberOfLine = 1
					}
				}
			} else {
				currentJournalLine += string(buffer)
				//textBoxWidth = terminalWidth - 4
				textBoxWidth = textBoxBorderWidth - 4
				numberOfLine = len(currentJournalLine) / textBoxWidth
				if len(currentJournalLine)%textBoxWidth != 0 {
					numberOfLine++
				}
				if numberOfLine == 0 {
					numberOfLine = 1
				}
			}
		}
	}()
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
		//tree
		//PrintAtColor(terminalWidth-6, reservedHeight, '🭅', FgBrown)
		//PrintAtColor(terminalWidth-5, reservedHeight, '🮋', FgBrown)
		//PrintAtColor(terminalWidth-4, reservedHeight, '🮋', FgBrown)
		//PrintAtColor(terminalWidth-3, reservedHeight, '🮋', FgBrown)
		//PrintAtColor(terminalWidth-1, reservedHeight, '🮋', FgBrown)
		//PrintAtColor(terminalWidth, reservedHeight, '🭎', FgBrown)
		//PrintAtColor(terminalWidth-5, reservedHeight-1, '▋', FgBrown)

		terminalWidth, terminalHeight, _ = term.GetSize(0)
		reservedHeight = terminalHeight - (3 + numberOfLine)
		textBoxBorderWidth = terminalWidth / 1

		PrintAt(0, reservedHeight+1, '╭', color.FgGreen)
		for x := 1; x < textBoxBorderWidth; x++ {
			PrintAt(x, reservedHeight+1, '─', color.FgGreen)
		}
		PrintAt(textBoxBorderWidth, reservedHeight+1, '╮', color.FgGreen)
		PrintAt(0, reservedHeight+2, '│', color.FgGreen)
		PrintAt(textBoxBorderWidth, reservedHeight+2, '│', color.FgGreen)
		PrintAt(2, reservedHeight+2, '>', color.FgYellow)
		// above this line never change
		if numberOfLine > 1 {
			for i := 0; i < numberOfLine; i++ {
				PrintAt(0, reservedHeight+i+3, '│', color.FgGreen)
				PrintAt(textBoxBorderWidth, reservedHeight+3+i, '│', color.FgGreen)
				PrintAt(0, terminalHeight, '╰', color.FgGreen)
				for x := 1; x < textBoxBorderWidth; x++ {
					PrintAt(x, terminalHeight, '─', color.FgGreen)
				}
				PrintAt(textBoxBorderWidth, reservedHeight+i+4, '╯', color.FgGreen)

			}
		} else {
			PrintAt(0, reservedHeight+3, '│', color.FgGreen)
			PrintAt(textBoxBorderWidth, reservedHeight+3, '│', color.FgGreen)
			PrintAt(0, reservedHeight+4, '╰', color.FgGreen)
			for x := 1; x < textBoxBorderWidth; x++ {
				PrintAt(x, reservedHeight+4, '─', color.FgGreen)
			}
			PrintAt(textBoxBorderWidth, reservedHeight+4, '╯', color.FgGreen)

		}

		fmt.Printf("\033[%d;%dH", reservedHeight+3, 4)
		//fmt.Print(currentJournalLine)
		textToDraw := currentJournalLine
		lines := numberOfLine
		for i := 0; i < lines; i++ {
			start := i * textBoxWidth
			end := start + textBoxWidth
			if end > len(textToDraw) {
				end = len(textToDraw)
			}

			lineSubString := textToDraw[start:end]
			y := reservedHeight + 3 + i
			x := 4
			fmt.Printf("\033[%d;%dH", y, x)
			fmt.Print(lineSubString)
		}
		//PrintAt(terminalWidth-6, reservedHeight, '/', color.FgHiRed)
		select {
		case input := <-inputChan:
			f, err := os.OpenFile("journal.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				log.Println(err)
			}
			defer f.Close()
			now := time.Now()
			formatedTime := now.Format("2006-01-02 15:04:05")
			if _, err := f.WriteString(formatedTime + " : " + input + "\n"); err != nil {
				log.Println(err)
			}
		case <-doneChan:
			return

		default:

		}
		for id := range leaves {
			leaves[id].Y = leaves[id].Y + leaves[id].Speed

			if leaves[id].Y >= reservedHeight {
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
			} else {
				PrintAt(leaves[id].X, leaves[id].Y, leaves[id].Charactere, leaves[id].Color)

			}
		}
		time.Sleep(100000000)

	}
}
