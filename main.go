package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand/v2"
	"os"
	"sync"
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

// const FgBrown = "\033[38;5;130m"
const FgBrown = "\033[38;2;150;67;33m"
const ColorReset = "\033[0m"

func PrintAtColor(x, y int, char rune, colorCode string) {
	fmt.Printf("\033[?25l\033[%d;%dH%s%c%s", y+1, x+1, colorCode, char, ColorReset)
}
func generateLeaves() {

}
func drawTree(width, floorY int) {
	treeX := width
	treeHeight := 20
	for i := 0; i < treeHeight; i++ {
		y := floorY - i
		PrintAtColor(treeX, y, '🮋', FgBrown)

	}
}
func drawBox(x, y, width, height int, colorToDraw color.Attribute) {
	PrintAt(x, y, '╭', colorToDraw)
	for i := 1; i < width-1; i++ {
		PrintAt(x+i, y, '─', colorToDraw)
	}
	PrintAt(x+width-1, y, '╮', colorToDraw)
	for i := 1; i < height-1; i++ {
		PrintAt(x, y+i, '│', colorToDraw)
		PrintAt(x+width-1, y+i, '│', colorToDraw)
	}
	PrintAt(x, y+height-1, '╰', colorToDraw)
	for i := 1; i < width-1; i++ {
		PrintAt(x+i, y+height-1, '─', colorToDraw)

	}
	PrintAt(x+width-1, y+height-1, '╯', colorToDraw)
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
	var dataMutex sync.Mutex
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
	inputChan := make(chan string, 1)
	doneChan := make(chan bool, 1)
	go func() {

		buffer := make([]byte, 1)

		for {
			os.Stdin.Read(buffer)
			if buffer[0] == 3 {
				doneChan <- true
				return
			}
			dataMutex.Lock()
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
			dataMutex.Unlock()
		}
	}()
	for count := 0; count <= 20; count++ {
		randomX := rand.IntN(terminalWidth - 1)
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

		//tree
		//PrintAtColor(terminalWidth-6, reservedHeight, '🭅', FgBrown)
		//PrintAtColor(terminalWidth-5, reservedHeight, '🮋', FgBrown)
		//PrintAtColor(terminalWidth-4, reservedHeight, '🮋', FgBrown)
		//PrintAtColor(terminalWidth-3, reservedHeight, '🮋', FgBrown)
		//PrintAtColor(terminalWidth-1, reservedHeight, '🮋', FgBrown)
		//PrintAtColor(terminalWidth, reservedHeight, '🭎', FgBrown)
		//PrintAtColor(terminalWidth-5, reservedHeight-1, '▋', FgBrown)
		terminalWidth, terminalHeight, _ = term.GetSize(0)
		textBoxBorderWidth = (terminalWidth / 3) * 2
		dataMutex.Lock()
		textToDraw := currentJournalLine
		lines := numberOfLine
		currentTextBoxWidth := textBoxWidth
		dataMutex.Unlock()
		boxHeight := 4 + lines
		reservedHeight = terminalHeight - boxHeight - 1
		screen.Clear()

		drawTree(terminalWidth, reservedHeight)

		drawBox(0, reservedHeight, textBoxBorderWidth, boxHeight+1, color.FgGreen)
		fmt.Printf("\033[%d;%dH", reservedHeight+1, (textBoxBorderWidth/2)-4)
		color.Set(color.FgGreen)
		color.Set(color.Italic)
		fmt.Print("─Journal─")
		color.Unset()
		PrintAt(2, reservedHeight+1, '>', color.FgYellow)

		for i := 0; i < lines; i++ {
			start := i * currentTextBoxWidth
			end := start + currentTextBoxWidth
			if end > len(textToDraw) {
				end = len(textToDraw)
			}
			if start < end {
				lineSubString := textToDraw[start:end]
				y := reservedHeight + 2 + i
				x := 4
				fmt.Printf("\033[%d;%dH", y, x)
				fmt.Print(lineSubString)

			}
		}
		drawBox(textBoxBorderWidth, reservedHeight, terminalWidth/3, boxHeight+1, color.FgHiYellow)
		color.Set(color.Italic)
		color.Set(color.FgHiYellow)

		fmt.Printf("\033[%d;%dH", reservedHeight+1, textBoxBorderWidth+(terminalWidth/3)/2-3)
		fmt.Print("─TODO─")
		color.Unset()
		color.Set(color.FgHiYellow)

		fmt.Printf("\033[%d;%dH", reservedHeight+2, textBoxBorderWidth+3)
		color.Set(color.Underline)
		fmt.Print(time.Now().Format("Mon, 02 Jan 2006 15:04 MST"))
		fmt.Printf("\033[%d;%dH", reservedHeight+3, textBoxBorderWidth+4)
		color.Unset()
		color.Set(color.FgYellow)
		f, err := os.OpenFile("todo.txt", os.O_CREATE|os.O_RDONLY, 0644)
		if err != nil {
			log.Println(err)
		}
		defer f.Close()
		scanner := bufio.NewScanner(f)
		scanner.Split(bufio.ScanLines)
		lineNum := 0
		for scanner.Scan() {
			fmt.Printf("\033[%d;%dH", reservedHeight+3+lineNum, textBoxBorderWidth+3)
			line := scanner.Text()
			todoWidth := terminalWidth - (textBoxBorderWidth + 4)
			if len(line) > todoWidth {
				line = line[:todoWidth]
			}
			fmt.Print(line)

			lineNum++
			if lineNum >= terminalHeight-reservedHeight-3 {
				break
			}
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
				randomX := rand.IntN(terminalWidth - 1)
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
