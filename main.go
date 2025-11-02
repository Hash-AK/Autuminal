package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/inancgumus/screen"
	"golang.org/x/term"
)

type Leaf struct {
	X          int
	Y          int
	Charactere rune
	Speed      int
	Color      string
}
type Weather struct {
	CurrentCondition []struct {
		FeelsLikeC  string `json:"FeelsLikeC"`
		FeelsLikeF  string `json:"FeelsLikeF"`
		WeatherDesc []struct {
			Value string `json:"value"`
		} `json:"weatherDesc"`
	} `json:"current_condition"`
}

func PrintAt(buffer *strings.Builder, x, y int, char rune, printColor string) {
	output := fmt.Sprintf("\033[?25l\033[%d;%dH%s%c%s", y+1, x+1, printColor, char, ColorReset)
	buffer.WriteString(output)
	//fmt.Print(x, ", ", y)i col
}
func SaveTodo(todo []string, file string) {
	os.Truncate(file, 0)
	f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Print(err)
	}
	defer f.Close()
	for i := range todo {
		todo[i] = strings.TrimSpace(todo[i])
		if todo[i] != "" {
			_, err := fmt.Fprintf(f, "%s\n", todo[i])
			if err != nil {
				log.Print(err)
			}
		}
	}

}

// const FgBrown = "\033[38;5;130m"
const FgBrown = "\033[38;2;150;67;33m"
const FgRed = "\033[38;2;150;0;0m"
const FgGreen = "\033[32m"
const Italic = "\033[3m"
const FgYellow = "\033[33m"
const FgGray = "\033[38;5;255m"
const Underline = "\033[4m"
const ColorReset = "\033[0m"

var activePanel = "journal"
var selectedTodoItem = 0
var tempUnit = "c"
var enableHacked = true
var isHacked = false
var tempChan = make(chan string, 1)
var tempDescChan = make(chan string, 1)

const treeArt = `
 ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢱⣸⠀⠀⠀⠀⠀⠀⠀⠀⡄⡄⢀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⠀⠀⠀⠀⠀⠀⠀
     ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⡀⢶⢠⠀⢀⡸⡄⠒⢺⠀⣸⣀⡀⣦⠽⠑⠁⠀⠀⠀⠀⠀⠀⠀⣆⣀⠗⠂⠀⠀⡆⢠⠃⡠⠜⠒⠀⠀⠀⠀⠀⠀
     ⠀⠀⠀⠀⠀⠀⠀⠀⠀⡄⠀⢤⠞⢳⠊⠓⠣⢸⡸⣲⠇⣘⣦⠚⢗⣻⠉⠻⡴⠂⢀⣀⠀⠀⣠⠂⠀⡇⠀⠀⠀⠀⡚⡲⢃⡉⠀⠀⠀⠀⠀⠀⠀⠀⠀
     ⠀⠀⠀⠀⠀⠀⠀⠀⠐⠺⠤⢼⡀⡞⢶⠦⣤⡖⠯⠭⡽⠟⡲⠀⠀⣆⠴⠊⢀⠀⠈⠅⡜⠒⠁⠀⠀⠉⢱⠀⠀⠀⠈⣑⡼⠁⢀⢠⢠⠄⢠⠆⠀⠀⠀
     ⠀⠀⠀⠀⠀⠀⠢⢄⢳⣁⣀⠆⠃⣇⡇⠜⠍⢳⡄⢰⢃⡈⡩⣲⠾⡃⢀⠀⠘⠤⢁⣠⠃⢠⢼⣇⣰⢃⣼⠀⠀⠀⣩⡾⠦⣆⠷⣅⠜⠉⠁⠀⠀⠀⠀
     ⠀⠀⠀⢦⠀⠈⠒⡥⣽⢁⠌⢹⢶⡤⡧⣾⠀⠀⠙⣾⣤⠖⠿⡿⣄⡗⢴⢣⡌⢲⣩⠚⠸⣌⣍⠹⣸⣚⡙⢷⣤⠞⠡⢄⣀⡳⣎⠀⠀⠀⠀⠀⠀⠀⠀
     ⠀⢄⣣⡈⠦⡜⣸⡹⣰⠃⡀⢱⣛⣰⣑⢽⣧⠀⢰⣿⡇⠰⠋⠑⡜⡗⡞⠋⠂⠘⢦⠳⣠⠿⠦⣼⢩⣤⢊⡾⠋⠀⠀⠀⠋⠀⢨⠏⠀⠀⠀⠀⠀⠀⠀
 ⠀⠀⢁⠇⠀⡏⠀⠈⢾⡄⠙⣤⠃⣟⠀⠋⣿⣅⡾⢻⢀⡀⡆⣰⣥⣟⢱⣞⣀⠀⣨⠧⣯⡀⠀⢸⣈⣷⡟⠀⢀⢦⠀⠀⠀⢠⠏⠀⠀⡀⣷⠀⠀⠀⠀
    ⠤⢲⠚⢒⢻⠙⢶⣴⢺⠉⠒⡧⠔⠛⠲⢤⣸⣿⠁⣼⡶⠿⠿⣽⣓⣸⢿⣓⡶⣚⢧⡷⣿⢫⣦⣸⣿⠏⢹⡴⠋⠸⡄⠀⠀⡞⠀⢰⣰⢣⠊⠀⣰⡠⠀
    ⠀⠈⡄⠀⢭⡇⡀⠉⠻⣇⠀⡇⠀⠀⠀⣀⡝⢿⡆⣿⢁⢀⡴⠋⣏⣏⡼⠋⡷⣇⡝⣇⣿⡜⠋⣿⣿⡆⣼⡝⡄⣠⢹⠀⣸⠁⠀⠀⠀⠛⣄⣸⡖⠊⠀
    ⠐⠴⣅⡆⠘⡎⢢⠀⠀⢹⣎⣷⠀⠀⣀⡕⠻⢚⣿⣿⡉⠉⠳⣄⣰⠟⠑⢶⠁⠹⢴⠁⡇⣠⣴⠿⣏⣾⡇⢹⡃⡗⢸⣷⢃⣠⠔⠋⠀⢠⠃⠀⠑⠹⠀
    ⠀⢤⢎⣈⡲⠵⣈⠉⠓⣾⠙⣾⣇⠀⠀⠛⣆⡇⢻⣿⡇⠀⣠⡾⠛⢶⡆⠈⣇⣰⠏⢰⣿⢏⡏⢠⣏⣼⠞⠉⠉⠱⣿⢿⡭⣄⠀⠀⢠⠏⠀⠀⠀⠀⠀
    ⠐⠚⠒⠂⠼⣄⠀⠉⠢⣼⡀⠈⢻⣆⠠⡄⠳⡇⢸⣿⣧⣾⡟⠀⠀⢸⡇⠀⣸⠋⠀⣼⡏⢾⠛⣿⢹⡏⠀⠀⢀⡼⠃⢘⠂⢨⠀⢀⡞⠀⢀⠄⢀⠆⡀
    ⠀⠀⠀⠀⠀⠈⠳⣄⠀⠈⠳⣄⠀⣿⣆⠸⡠⠜⣆⣿⣿⠏⠀⠀⠀⢸⡇⢰⠇⠀⢀⣿⠁⣿⢰⡇⣼⠁⠀⢠⡞⠁⠀⠸⣚⣮⠵⠟⠓⠦⣸⠀⡤⠼⠓
    ⠀⠀⠀⠀⠀⠀⠀⠙⢦⣀⣀⣈⠳⣜⢿⣯⠀⠀⢈⣿⡿⠦⣤⣀⠀⢸⣷⡏⠀⠀⣸⣿⡾⠋⣿⢁⡟⠀⣰⣯⣤⠶⠞⣋⠽⢓⣒⡡⠤⠒⠛⠳⢧⡀⡄
    ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠉⠉⠉⠉⠙⠳⣿⣷⡀⢸⣿⡇⠀⠀⠉⠛⢾⣿⠀⠀⠀⣿⡟⠁⣸⣿⣾⣿⣿⠟⢉⣠⣴⠞⠋⠉⠉⠉⠂⠀⠀⠀⠀⠈⠃⠀
    ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⢻⣿⣾⣿⡇⠀⠀⠀⠀⢸⣿⠀⠀⢸⡟⢀⣼⡿⠋⣼⣿⣿⡿⠛⠉⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
    ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠙⣿⣿⡇⠀⠀⠀⠀⠀⣿⡀⠀⣿⣷⡿⠋⠀⢠⣿⠟⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
    ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢸⣿⣿⣄⠀⠀⠀⠀⢿⡇⣸⣿⠟⠀⠀⢀⣾⡏⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
    ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢸⣿⣿⣿⣆⠀⠀⠀⣸⣷⣿⡇⠀⠀⠀⣼⡟⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
    ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠸⣿⡏⢿⣿⣦⣀⣾⣿⢯⣿⠀⠀⠀⣼⡟⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
    ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢠⣿⣿⣮⣿⣿⣿⡿⠁⣸⡟⠀⠀⣼⡟⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
    ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠘⢿⣿⣿⣿⣿⡟⠀⢠⣿⠃⠀⣼⡿⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
    ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠹⣿⣿⣿⣷⣠⣾⣿⣤⣾⣿⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
    ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠸⣿⣿⣿⣿⣿⣿⠟⠋⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
    ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣿⣿⣿⣿⣿⡟⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
    ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣿⣿⣿⣿⣿⡟⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
    ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣿⣿⣿⣿⣿⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
    ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣿⣿⣿⣿⣿⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
    ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣿⣿⣿⣿⣿⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
    ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣿⣿⣿⣿⣿⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
    ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣿⣿⣿⣿⣿⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
    ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣿⣿⣿⣿⣿⣷⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
    ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢰⣿⣿⣿⣿⣿⣿⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
    ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢸⣿⣿⣿⣿⣿⣿⡆⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
    ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣿⣿⣿⣿⣿⣿⣿⣇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
    ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣸⣿⣿⣿⣿⣿⣿⣿⣿⡄⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
    ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣰⡿⠿⠛⠻⣿⣿⠿⠿⠿⢿⣄⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
    ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠁⠀⠀⠀⠀⠈⠡⠀⠀⠀⠀⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀`

func drawTree(buffer *strings.Builder, width, height int) {
	if !isHacked {
		lines := strings.Split(treeArt, "\n")
		artHeight := len(lines)
		startY := height - artHeight
		if startY < 1 {
			startY = 1
		}
		for i, line := range lines {
			visualLen := len([]rune(line)) / 2
			startX := (width - visualLen) / 2
			if startX < 1 {
				startX = 1
			}
			y := startY + i
			output := fmt.Sprintf("\033[%d;%dH%s%s%s", y, startX, FgBrown, line, ColorReset)
			buffer.WriteString(output)
		}
	} else {
		chance := rand.IntN(3)
		switch chance {
		case 0, 1:
			lines := strings.Split(treeArt, "\n")
			artHeight := len(lines)
			startY := height - artHeight
			if startY < 1 {
				startY = 1
			}
			for i, line := range lines {
				visualLen := len([]rune(line)) / 2
				startX := (width - visualLen) / 2
				if startX < 1 {
					startX = 1
				}
				y := startY + i
				output := fmt.Sprintf("\033[%d;%dH%s%s%s", y, startX, FgRed, line, ColorReset)
				buffer.WriteString(output)
			}
		case 2, 3:
			return

		}
	}
}
func PrintAtColor(buffer *strings.Builder, x, y int, char rune, colorCode string) {
	output := fmt.Sprintf("\033[?25l\033[%d;%dH%s%c%s", y+1, x+1, colorCode, char, ColorReset)
	buffer.WriteString(output)
}

func drawBox(buffer *strings.Builder, x, y, width, height int, colorToDraw string) {
	if !isHacked {
		PrintAt(buffer, x, y, '╭', colorToDraw)
		for i := 1; i < width-1; i++ {
			PrintAt(buffer, x+i, y, '─', colorToDraw)
		}
		PrintAt(buffer, x+width-1, y, '╮', colorToDraw)
		for i := 1; i < height-1; i++ {
			PrintAt(buffer, x, y+i, '│', colorToDraw)
			PrintAt(buffer, x+width-1, y+i, '│', colorToDraw)
		}
		PrintAt(buffer, x, y+height-1, '╰', colorToDraw)
		for i := 1; i < width-1; i++ {
			PrintAt(buffer, x+i, y+height-1, '─', colorToDraw)

		}
		PrintAt(buffer, x+width-1, y+height-1, '╯', colorToDraw)
	} else {
		chance := rand.IntN(3)
		switch chance {
		case 0, 1:
			PrintAt(buffer, x, y, '╭', FgRed)
			for i := 1; i < width-1; i++ {
				PrintAt(buffer, x+i, y, '─', FgRed)
			}
			PrintAt(buffer, x+width-1, y, '╮', FgRed)
			for i := 1; i < height-1; i++ {
				PrintAt(buffer, x, y+i, '│', FgRed)
				PrintAt(buffer, x+width-1, y+i, '│', FgRed)
			}
			PrintAt(buffer, x, y+height-1, '╰', FgRed)
			for i := 1; i < width-1; i++ {
				PrintAt(buffer, x+i, y+height-1, '─', FgRed)

			}
		case 2, 3:
			PrintAt(buffer, x, y, '╭', FgYellow)
			for i := 1; i < width-1; i++ {
				PrintAt(buffer, x+i, y, '─', FgYellow)
			}
			PrintAt(buffer, x+width-1, y, '╮', FgYellow)
			for i := 1; i < height-1; i++ {
				PrintAt(buffer, x, y+i, '│', FgYellow)
				PrintAt(buffer, x+width-1, y+i, '│', FgYellow)
			}
			PrintAt(buffer, x, y+height-1, '╰', FgYellow)
			for i := 1; i < width-1; i++ {
				PrintAt(buffer, x+i, y+height-1, '─', FgYellow)

			}
		}
	}
}
func fetchWeather(unit string) {
	response, err := http.Get("https://wttr.in/?format=j1")
	var temperature string
	var WeatherDesc string
	if err != nil {
		tempChan <- "Weather N/A"
	}
	defer response.Body.Close()
	var weatherData Weather
	err = json.NewDecoder(response.Body).Decode(&weatherData)
	if err != nil {
		tempChan <- "Weather N/A"
	}
	if len(weatherData.CurrentCondition) > 0 {
		current := weatherData.CurrentCondition[0]
		if unit == "c" {
			temperature = current.FeelsLikeC

		} else {
			temperature = current.FeelsLikeF
		}
		desc := current.WeatherDesc[0]
		WeatherDesc = desc.Value

	}
	finalTemp := strings.Split(strings.TrimSpace(temperature), "\n")
	tempChan <- finalTemp[0]
	tempDescChan <- WeatherDesc

}
func main() {
	var buffer strings.Builder
	defer fmt.Printf("\033[?25h")

	var leaves []Leaf
	var terminalWidth int
	var terminalHeight int
	var reservedHeight int
	var currentJournalLine string
	var newTodoInput = ""
	var isAddingTodo = false
	var textBoxWidth int
	var textBoxBorderWidth int
	var boxHeight int
	var frameCount = 0
	var todoLines []string
	var showSettings = false
	var showTree = true
	var leafStyle = "autumn"
	var temperature string
	var weatherDesc string
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
	inputChan := make(chan byte, 10)
	doneChan := make(chan bool, 1)
	saveJournalChan := make(chan string, 1)
	go func() {

		buffer := make([]byte, 3)

		for {
			n, _ := os.Stdin.Read(buffer)
			if n == 1 {
				inputChan <- buffer[0]
			} else if n == 3 && buffer[0] == 27 && buffer[1] == 91 {
				if activePanel == "todo" {
					if buffer[2] == 'A' {
						inputChan <- 'k'
					}
					if buffer[2] == 'B' {
						inputChan <- 'j'
					}
				}
			}
		}
	}()
	go func() {
		for {
			fetchWeather(tempUnit)
			time.Sleep(15 * time.Minute)
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
		var randomColor string
		randomColorNum := rand.IntN(3)
		switch randomColorNum {
		case 0:
			randomColor = FgRed
		case 1:
			randomColor = FgYellow
		case 2:
			randomColor = FgRed
		case 3:
			randomColor = FgYellow

		}
		randomSpeed := rand.IntN(5)
		if randomSpeed == 0 {
			randomSpeed++
		}
		leaves = append(leaves, Leaf{X: randomX, Y: randomY, Charactere: randomChar, Speed: randomSpeed, Color: randomColor})

	}

	for {
		buffer.Reset()
		frameCount++
		todoContent, err := os.ReadFile("todo.txt")
		if err == nil {
			todoLines = strings.Split(string(todoContent), "\n")

		}
		for len(inputChan) > 0 {
			key := <-inputChan
			if key == 9 {
				if activePanel == "journal" {
					activePanel = "todo"
				} else if activePanel == "todo" {
					activePanel = "journal"
				}
			}

			if showSettings {
				switch key {

				case 27:

					showSettings = false

				case 's':
					enableHacked = !enableHacked
				case 'l':
					if leafStyle == "autumn" {
						leafStyle = "matrix"

					} else {
						leafStyle = "autumn"
					}
				case 't':
					showTree = !showTree
				case 'u':
					if tempUnit == "c" {
						tempUnit = "f"
					} else {
						tempUnit = "c"
					}
					go fetchWeather(tempUnit)
				}
				continue
			} else {
				if isAddingTodo {
					switch key {
					case 13:
						if strings.TrimSpace(newTodoInput) != "" {
							todoLines = append(todoLines, newTodoInput)
							SaveTodo(todoLines, "todo.txt")
						}
						isAddingTodo = false
						newTodoInput = ""
					case 8, 127:
						if len(newTodoInput) > 0 {
							newTodoInput = newTodoInput[:len(newTodoInput)-1]
						}

					case 27:
						isAddingTodo = false
						newTodoInput = ""

					default:
						newTodoInput += string(key)

					}
					continue
				}
				if activePanel == "journal" {
					switch key {

					case 3:
						doneChan <- true
						return
					case 13:
						saveJournalChan <- currentJournalLine
						currentJournalLine = ""

					case 8, 127:
						if len(currentJournalLine) > 0 {
							currentJournalLine = currentJournalLine[:len(currentJournalLine)-1]
						}

					default:
						currentJournalLine += string(key)
					}
					continue
				} else if activePanel == "todo" {
					switch key {
					case 'k':
						if selectedTodoItem > 0 {
							selectedTodoItem--
						}
					case 'j':
						if selectedTodoItem < len(todoLines)-1 {
							selectedTodoItem++
						}
					case 'd':
						if len(todoLines) > 1 {
							todoLines = append(todoLines[:selectedTodoItem], todoLines[selectedTodoItem+1:]...)

							if selectedTodoItem > len(todoLines)-1 {
								selectedTodoItem--
							}
							SaveTodo(todoLines, "todo.txt")

						}
					case 'a':
						isAddingTodo = true
						newTodoInput = ""
						continue

					case 3:
						doneChan <- true
						return
					}
					continue
				}

			}
		}
		for len(tempChan) > 0 {
			temperature = <-tempChan
		}
		for len(tempDescChan) > 0 {
			weatherDesc = <-tempDescChan
		}
		if !enableHacked {
			isHacked = false
		}
		terminalWidth, terminalHeight, _ = term.GetSize(0)
		textBoxBorderWidth = (terminalWidth / 3) * 2
		textBoxWidth = textBoxBorderWidth - 4
		textToDraw := currentJournalLine
		currentTextBoxWidth := textBoxWidth
		var lines int
		if textBoxWidth > 0 {
			lines = (len(textToDraw) / textBoxWidth) + 1
		} else {
			lines = 1
		}
		if lines < 4 {
			boxHeight = 4

		} else {
			boxHeight = 1 + lines
		}
		reservedHeight = terminalHeight - boxHeight - 2
		if !showSettings {
			if showTree {
				drawTree(&buffer, terminalWidth, reservedHeight+2)
			}
			drawBox(&buffer, 0, reservedHeight, textBoxBorderWidth, boxHeight+1, FgGreen)
			buffer.WriteString(fmt.Sprintf("\033[%d;%dH%s%s%s", reservedHeight+1, (textBoxBorderWidth/2)-4, FgGreen, "─Journal─", ColorReset))
			PrintAt(&buffer, 2, reservedHeight+1, '>', FgYellow)

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
					buffer.WriteString(fmt.Sprintf("\033[%d;%dH", y, x))
					buffer.WriteString(fmt.Sprint(lineSubString))

				}
			}

			drawBox(&buffer, textBoxBorderWidth, reservedHeight, terminalWidth/3, boxHeight+1, FgYellow)

			buffer.WriteString(fmt.Sprintf("\033[%d;%dH%s%s%s", reservedHeight+1, textBoxBorderWidth+(terminalWidth/3)/2-3, FgYellow, "─TODO─", ColorReset))

			buffer.WriteString(fmt.Sprintf("\033[%d;%dH", reservedHeight+2, textBoxBorderWidth+3))
			buffer.WriteString(fmt.Sprint(time.Now().Format("Mon, 02 Jan 2006 15:04 MST")))
			buffer.WriteString(fmt.Sprintf("\033[%d;%dH", reservedHeight+3, textBoxBorderWidth+4))
			lineNum := 0
			bufferHeight := 4
			if activePanel == "todo" && isAddingTodo {
				bufferHeight = 5

			} else {
				bufferHeight = 4
			}
			for i, line := range todoLines {
				buffer.WriteString(fmt.Sprintf("\033[%d;%dH", reservedHeight+3+lineNum, textBoxBorderWidth+3))
				todoWidth := terminalWidth - (textBoxBorderWidth + 4)
				if len(line) > todoWidth {
					line = line[:todoWidth]
				}
				if activePanel == "todo" {
					if i == selectedTodoItem {
						line = fmt.Sprintf("%s%s%s", FgRed, ">", ColorReset) + line
					}
				}

				buffer.WriteString(line)
				lineNum++
				if lineNum >= terminalHeight-reservedHeight-bufferHeight {
					break
				}

			}
			if activePanel == "todo" && isAddingTodo {
				inputY := reservedHeight + 3 + lineNum
				inputX := textBoxBorderWidth + 3
				prompt := "New: "
				buffer.WriteString(fmt.Sprintf("\033[%d;%dH%s%s%s", inputY, inputX, FgYellow, prompt, ColorReset))
				buffer.WriteString(newTodoInput)
				buffer.WriteString(fmt.Sprintf("\033[?25h\033[%d;%dH", inputY, inputX+len(prompt)+len(newTodoInput)))
			}
			statusBarY := terminalHeight
			statusBarText := "Ctrl+C : Quit | /settings: Open Settings" + " | " + weatherDesc + ", " + temperature + "°" + tempUnit
			paddedText := fmt.Sprintf("%-*s", terminalWidth, statusBarText)
			buffer.WriteString(fmt.Sprintf("\033[%d;%dH\033[48;5;235m%s\033[49m", statusBarY, 1, paddedText))
			select {
			case input := <-saveJournalChan:
				f, err := os.OpenFile("journal.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				if err != nil {
					log.Println(err)
				}
				now := time.Now()
				formatedTime := now.Format("2006-01-02 15:04:05")
				if _, err := f.WriteString(formatedTime + " : " + input + "\n"); err != nil {
					log.Println(err)
				}
				f.Close()
				if enableHacked {
					if strings.Contains(input, "scary") || strings.Contains(input, "spooky") {
						isHacked = true
					}
				} else {
					isHacked = false
				}
				if strings.Contains(input, "stop") {
					isHacked = false
				}

				if strings.Contains(input, "/settings") {
					showSettings = true
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
					if !isHacked {
						if leafStyle == "autumn" {
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
						} else {
							katakana := []rune("アァカサタナハマヤャラワガザダバパイィキシチニヒミリヰギジヂビピウゥクスツヌフムユュルグズブヅプエェケセテネヘメレヱゲゼデベペオォコソトノホモヨョロヲゴゾドボポヴ")
							randomChar = katakana[rand.IntN(len(katakana))]
						}
					} else {
						if leafStyle == "autumn" {
							switch randomCharNum {
							case 0:
								randomChar = '?'
							case 1:
								randomChar = '!'
							case 2:
								randomChar = '%'
							case 3:
								randomChar = '$'
							case 4:
								randomChar = '\\'
							default:
								randomChar = '='
							}
						} else {

							katakana := []rune("アァカサタナハマヤャラワガザダバパイィキシチニヒミリヰギジヂビピウゥクスツヌフムユュルグズブヅプエェケセテネヘメレヱゲゼデベペオォコソトノホモヨョロヲゴゾドボポヴ")
							randomChar = katakana[rand.IntN(len(katakana))]
						}
					}
					randomSpeed := rand.IntN(5)
					if randomSpeed == 0 {
						randomSpeed++
					}
					if isHacked {
						randomSpeed++
					}
					var randomColor string
					randomColorNum := rand.IntN(3)
					if !isHacked {
						if leafStyle == "autumn" {
							switch randomColorNum {
							case 0:
								randomColor = FgYellow
							case 1:
								randomColor = FgRed
							case 2:
								randomColor = FgYellow
							case 3:
								randomColor = FgRed
							}
						} else {
							randomColor = FgGreen
						}
					} else {
						randomColor = FgRed
					}
					leaves[id].Y = randomY
					leaves[id].X = randomX
					leaves[id].Charactere = randomChar
					leaves[id].Speed = randomSpeed
					leaves[id].Color = randomColor
				} else {
					PrintAt(&buffer, leaves[id].X, leaves[id].Y, leaves[id].Charactere, leaves[id].Color)

				}
			}
		} else {
			drawBox(&buffer, 5, 5, terminalWidth-10, terminalHeight-10, FgBrown)
			title := "─Settings (Press ESC to quit)─"
			titleX := 5 + (terminalWidth-len(title))/2
			buffer.WriteString(fmt.Sprintf("\033[%d;%dH%s%s%s", 6, titleX, FgBrown, title, ColorReset))
			enableHackedText := "[S]cary mode : "
			enableHackedStatus := "[ON]"
			if !enableHacked {
				enableHackedStatus = "[OFF]"
			}
			buffer.WriteString(fmt.Sprintf("\033[%d;%dH%s%s", 8, 10, enableHackedText, enableHackedStatus))
			leafStyleText := "[L]eaf style :"
			buffer.WriteString(fmt.Sprintf("\033[%d;%dH%s<%s>", 10, 10, leafStyleText, leafStyle))

			showTreeText := "[T]ree : "
			showTreeStatus := "[ON]"
			if !showTree {
				showTreeStatus = "[OFF]"
			}
			buffer.WriteString(fmt.Sprintf("\033[%d;%dH%s%s", 12, 10, showTreeText, showTreeStatus))
			tempUnitText := "Weather temperature [u]nit : "
			tempUnitValue := "celsius"
			if tempUnit == "c" {
				tempUnitValue = "celsius"
			} else {
				tempUnitValue = "fahrenheit"
			}
			buffer.WriteString(fmt.Sprintf("\033[%d;%dH%s<%s>", 14, 10, tempUnitText, tempUnitValue))
		}

		screen.Clear()
		fmt.Print(buffer.String())
		time.Sleep(time.Millisecond * 150)
	}
}
