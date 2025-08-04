package main

import (
	"unicode"

	"github.com/gdamore/tcell/v2"
)

type InputState struct {
	Value    string
	Cursor   int
	Max      int
	OnInput  func(text string)
	OnEscape func()
}

func (input *InputState) SetCursor(pos int) bool {
	if pos < 1 {
		pos = 0
		goto SetCursor
	}
	if pos >= len(input.Value) {
		pos = len(input.Value)
		goto SetCursor
	}
	if input.Max > 0 && pos > input.Max {
		pos = input.Max
		goto SetCursor
	}
SetCursor:
	if input.Cursor == pos {
		return false
	}
	input.Cursor = pos
	return true
}

func jumpWord(input *InputState, inc int) {
	if len(input.Value) == 0 {
		return
	}
	foundLetter := false
	for input.SetCursor(input.Cursor+inc) && input.Cursor < len(input.Value) {
		isLetter := unicode.IsLetter(rune(input.Value[(input.Cursor)]))
		if !foundLetter && isLetter {
			foundLetter = true
			continue
		}
		if foundLetter && !isLetter {
			if inc < 0 {
				input.SetCursor(input.Cursor - inc)
			}
			break
		}
	}
}

func (input *InputState) HandleInput(tev *tcell.EventKey) {
	onInput := func(text string) {
		if input.OnInput != nil {
			input.OnInput(text)
		}
	}
	onEscape := func() {
		if input.OnEscape != nil {
			input.OnEscape()
		}
	}
	switch tev.Rune() {
	case 'd':
		if (tev.Modifiers() & tcell.ModAlt) == 0 {
			break
		}
		from := input.Cursor
		jumpWord(input, +1)
		to := input.Cursor
		input.Value = input.Value[:from] +
			input.Value[to:]
		input.SetCursor(from)
		onInput(input.Value)
		return
	}
	switch tev.Key() {
	case tcell.KeyEsc:
		input.SetCursor(0)
		onInput("")
		onEscape()
	case tcell.KeyCtrlZ, tcell.KeyCtrlL:
		input.SetCursor(0)
		onInput("")
	case tcell.KeyCtrlB:
		input.SetCursor(input.Cursor - 1)
	case tcell.KeyCtrlF:
		input.SetCursor(input.Cursor + 1)
	case tcell.KeyCtrlW:
		to := input.Cursor
		jumpWord(input, -1)
		from := input.Cursor
		input.Value = input.Value[:from] +
			input.Value[to:]
		onInput(input.Value)
	case tcell.KeyLeft:
		if (tev.Modifiers() & tcell.ModCtrl) == 0 {
			input.SetCursor(input.Cursor - 1)
			break
		}
		jumpWord(input, -1)
	case tcell.KeyRight:
		if (tev.Modifiers() & tcell.ModCtrl) == 0 {
			input.SetCursor(input.Cursor + 1)
			break
		}
		jumpWord(input, +1)
	case tcell.KeyTab, tcell.KeyEnter:
		onInput(input.Value)
		onEscape()
	case tcell.KeyBackspace, tcell.KeyBackspace2:
		if input.Cursor == 0 {
			break
		}
		input.Value = input.Value[:input.Cursor-1] +
			input.Value[input.Cursor:]
		input.SetCursor(input.Cursor - 1)
		onInput(input.Value)
	case tcell.KeyDelete, tcell.KeyCtrlD:
		if input.Cursor == len(input.Value) {
			break
		}
		input.Value = input.Value[:input.Cursor] +
			input.Value[input.Cursor+1:]
		onInput(input.Value)
	default:
		if (input.Max > 0 && len(input.Value) >= input.Max) ||
			tev.Rune() == 0 ||
			tev.Modifiers() > 0 {
			break
		}
		char := string(tev.Rune())
		input.Value = input.Value[:input.Cursor] +
			char +
			input.Value[input.Cursor:]
		input.SetCursor(input.Cursor + 1)
		onInput(input.Value)
	}
}
