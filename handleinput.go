package main

import (
	"github.com/gdamore/tcell/v2"
)

func handleInput(tev *tcell.EventKey, str string, onInput func(text string), onEscape func()) {
	switch tev.Key() {
	case tcell.KeyEsc:
		onInput("")
		onEscape()
	case tcell.KeyTab, tcell.KeyEnter:
		onInput(str)
		onEscape()
	case tcell.KeyBackspace, tcell.KeyBackspace2:
		if len(str) > 0 {
			str = str[:len(str)-1]
			onInput(str)
		}
	default:
		if tev.Rune() != 0 && tev.Modifiers() == 0 {
			str += string(tev.Rune())
			onInput(str)
		}
	}
}
