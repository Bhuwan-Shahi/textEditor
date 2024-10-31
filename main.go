package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/nsf/termbox-go"
)

var mode int       //view and edit mode
var ROWS, COLS int //width and height of the terminal
var offsetCol, offsetRow int
var currentRow, currentCol = 1, 1
var sourceFile string

var textBuffer = [][]rune{}
var undoBuffer = [][]rune{}
var copyBuffer = [][]rune{}

var modified bool

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

func insertRune(event termbox.Event) {
	insertrune := make([]rune, len(textBuffer[currentRow])+1)
	copy(insertrune[:currentCol], textBuffer[currentRow][:currentCol])
	if event.Key == termbox.KeySpace {
		insertrune[currentCol] = rune(' ')
	} else if event.Key == termbox.KeyTab {
		insertrune[currentCol] = rune(' ')
	} else {
		insertrune[currentCol] = rune(event.Ch)
	}
	copy(insertrune[currentCol+1:], textBuffer[currentRow][currentCol:])
	textBuffer[currentRow] = insertrune
	currentCol++
}

// Scroling
func scrollTextBuffer() {
	if currentRow < offsetRow {
		offsetRow = currentRow
	}
	if currentCol < offsetCol {
		offsetCol = currentCol
	}
	if currentRow >= offsetRow+ROWS {
		offsetRow = currentRow - ROWS + 1
	}
	if currentCol >= offsetCol+COLS {
		offsetCol = currentCol - COLS + 1
	}
}

// Displaying the text buffers
func displayTextBuffer() {
	var row, col int
	for row = 0; row < ROWS; row++ {
		textBufferRow := row + offsetRow
		for col = 0; col < COLS; col++ {
			textBufferCol := col + offsetCol
			if textBufferRow >= 0 && textBufferRow < len(textBuffer) && textBufferCol < len(textBuffer[textBufferRow]) {
				if textBuffer[textBufferRow][textBufferCol] != '\t' {
					termbox.SetChar(col, row, textBuffer[textBufferRow][textBufferCol])
				} else {
					termbox.SetCell(col, row, rune(' '), termbox.ColorDefault, termbox.ColorGreen)
				}
			} else if row+offsetRow > len(textBuffer) {
				termbox.SetCell(0, row, rune('*'), termbox.ColorBlue, termbox.ColorDefault)
				termbox.SetChar(col, row, rune('\n'))
			}
		}
	}

}

// Display status bar
func displayStatusBar() {
	var modeStatus string
	var fileStatus string
	var copyStatus string
	var undoStatus string
	var cursorStatus string
	if mode > 0 {
		modeStatus = " EDIT: "
	} else {
		modeStatus = " VIEW: "
	}
	filename_length := len(sourceFile)
	if filename_length > 8 {
		filename_length = 8
	}
	fileStatus = sourceFile[:filename_length] + " - " + strconv.Itoa(len(textBuffer)) + " lines"
	if modified {
		fileStatus += " modified"
	} else {
		fileStatus += " saved"
	}
	cursorStatus = " Row " + strconv.Itoa(currentRow+1) + ", Col " + strconv.Itoa(currentCol+1) + " "
	if len(copyBuffer) > 0 {
		copyStatus = " [Copy]"
	}
	if len(undoBuffer) > 0 {
		undoStatus = " [Undo]"
	}
	usedSpace := len(modeStatus) + len(fileStatus) + len(cursorStatus) + len(copyStatus) + len(undoStatus)
	spaces := strings.Repeat(" ", COLS-usedSpace)
	message := modeStatus + fileStatus + copyStatus + undoStatus + spaces + cursorStatus
	printMessage(0, ROWS, termbox.ColorBlack, termbox.ColorWhite, message)
}

func printMessage(col, row int, fg, bg termbox.Attribute, message string) {
	for i, ch := range message {
		if col+i >= COLS {
			break
		}
		termbox.SetCell(col+i, row, ch, fg, bg)
	}
}

func getKey() termbox.Event {
	var keyEvent termbox.Event
	switch event := termbox.PollEvent(); event.Type {
	case termbox.EventKey:
		keyEvent = event
	case termbox.EventError:
		panic(event)

	}
	return keyEvent
}

func processKeyPress() {
	keyEvent := getKey()
	if keyEvent.Key == termbox.KeyEsc {
		mode = 0
	} else if keyEvent.Ch != 0 {
		if mode == 1 {
			insertRune(keyEvent)
			modified = true
		} else {
			switch keyEvent.Ch {
			case 'q':
				termbox.Close()
				os.Exit(0)
			case 'e':
				mode = 1
			}
		}
		//Handeling the chars in the form of rune

	} else {
		switch keyEvent.Key {
		case termbox.KeyTab:
			if mode == 1 {
				for i := 0; i < 4; i++ {
					insertRune(keyEvent)
					modified = true
				}
			}
		case termbox.KeySpace:
			if mode == 1 {
				insertRune(keyEvent)
				modified = true
			}
		case termbox.KeyHome:
			currentCol = 0
		case termbox.KeyEnd:
			currentCol = len(textBuffer[currentRow])
		case termbox.KeyPgup:
			if currentRow-int(ROWS/4) > 0 {
				currentRow -= int(ROWS / 4)
			}
		case termbox.KeyPgdn:
			if currentRow+int(ROWS/4) < len(textBuffer)-1 {
				currentRow += int(ROWS / 4)
			}
		case termbox.KeyArrowUp:
			if currentRow != 0 {
				currentRow--
			}
		case termbox.KeyArrowDown:
			if currentRow < len(textBuffer)-1 {
				currentRow++
			}
		case termbox.KeyArrowLeft:
			if currentCol != 0 {
				currentCol--
			} else if currentRow > 0 {
				currentRow--
				currentCol = len(textBuffer[currentRow])
			}
		case termbox.KeyArrowRight:
			if currentCol < len(textBuffer[currentRow]) {
				currentCol++
			} else if currentRow < len(textBuffer)-1 {
				currentRow++
				currentCol = 0
			}
		}
		if currentCol > len(textBuffer[currentRow]) {
			currentCol = len(textBuffer[currentRow])
		}
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
		sourceFile = "out.txt"
		textBuffer = append(textBuffer, []rune{})
	}
	for {
		COLS, ROWS = termbox.Size()
		ROWS--
		if COLS < 78 {
			COLS = 78
		}
		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
		scrollTextBuffer()
		displayTextBuffer()
		displayStatusBar()
		termbox.SetCursor(currentCol-offsetCol, currentRow-offsetRow)
		termbox.Flush()
		processKeyPress()

	}
}

func main() {
	runEditor()
}
