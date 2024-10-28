package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

var ROWS, COLS int //width and height of the terminal
var offsetX, offsetY int
var sourceFile string

var textBuffer = [][]rune{}

func readFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		textBuffer = append(textBuffer, []rune(line))
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}
}

// Displaying the text buffers
func displayTextBuffer() {
	var row, col int
	for row = 0; row < ROWS; row++ {
		textBufferRow := row + offsetY
		for col = 0; col < COLS; col++ {
			textBufferCol := col + offsetX
			if textBufferRow >= 0 && textBufferRow < len(textBuffer) && textBufferCol < len(textBuffer[textBufferRow]) {
				if textBuffer[textBufferRow][textBufferCol] != '\t' {
					termbox.SetChar(col, row, textBuffer[textBufferRow][textBufferCol])
				} else {
					termbox.SetCell(col, row, rune(' '), termbox.ColorDefault, termbox.ColorGreen)
				}
			} else if row+offsetY > len(textBuffer) {
				termbox.SetCell(0, row, rune('*'), termbox.ColorBlue, termbox.ColorDefault)
				termbox.SetChar(col, row, rune('\n'))
			}
		}
	}

}

func printMessage(col, row int, fg, bg termbox.Attribute, message string) {
	for _, ch := range message {
		termbox.SetCell(col, row, ch, fg, bg)
		col += runewidth.RuneWidth(ch)
	}
}
func runEditor() {
	err := termbox.Init()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if len(os.Args) > 1 {
		sourceFile = os.Args[1]
		readFile(sourceFile)
	} else {
		textBuffer = append(textBuffer, []rune{})
		sourceFile = "out.txt"

	}
	ROWS, COLS = termbox.Size()
	ROWS--
	if COLS < 78 {
		COLS = 78
	}
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	displayTextBuffer()
	defer termbox.Close()
	termbox.Flush()
	//always listening to the keyboard
	termbox.PollEvent()
}

func main() {
	runEditor()
}
