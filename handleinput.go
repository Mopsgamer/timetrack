package main

import (
	"unicode"

	"github.com/gdamore/tcell/v2"
)

type InputState struct {
	Value     string
	Cursor    int
	Max       int
	OnChange  func(text string)
	OnProceed func(text string)
	OnEscape  func()
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
	onChange := func(text string) {
		if text == "" {
			input.Cursor = 0
		}
		input.Value = text
		if input.OnChange != nil {
			input.OnChange(text)
		}
	}
	onEscape := func() {
		if input.OnEscape != nil {
			input.OnEscape()
		}
	}
	onProceed := func(text string) {
		if input.OnProceed != nil {
			input.OnProceed(text)
		}
		onChange("")
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
		onChange(input.Value)
		return
	}
	switch tev.Key() {
	case tcell.KeyEsc:
		input.SetCursor(0)
		onChange("")
		onEscape()
	case tcell.KeyCtrlZ, tcell.KeyCtrlL:
		input.SetCursor(0)
		onChange("")
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
		onChange(input.Value)
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
	case tcell.KeyTab:
		onChange(input.Value)
		onEscape()
	case tcell.KeyEnter:
		onProceed(input.Value)
		onEscape()
	case tcell.KeyBackspace, tcell.KeyBackspace2:
		if input.Cursor == 0 {
			break
		}
		input.Value = input.Value[:input.Cursor-1] +
			input.Value[input.Cursor:]
		input.SetCursor(input.Cursor - 1)
		onChange(input.Value)
	case tcell.KeyDelete, tcell.KeyCtrlD:
		if input.Cursor == len(input.Value) {
			break
		}
		input.Value = input.Value[:input.Cursor] +
			input.Value[input.Cursor+1:]
		onChange(input.Value)
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
		onChange(input.Value)
	}
}
