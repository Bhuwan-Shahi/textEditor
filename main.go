package main

import (
	"fmt"
	"os"

	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

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
	defer termbox.Close()
	printMessage(25, 11, termbox.ColorDefault, termbox.ColorDefault, "GRO - A bare bone text editor")
	termbox.Flush()
	//always listening to the keyboard
	termbox.PollEvent()
}

func main() {
	runEditor()
}
