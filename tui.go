package main

import (
	"os"
	"regexp"
	"slices"

	"github.com/gdamore/tcell/v2"
)

type Window uint

const (
	StateList Window = iota
	StateSearch
)

type Item struct {
	Name     string
	Selected bool
}

type State struct {
	Window        Window
	Items         []Item
	ItemsFound    []*Item
	ItemCursor    *Item
	SearchContent string
}

func (state *State) SearchRegexp() (*regexp.Regexp, error) {
	if state.SearchContent == "" {
		return nil, nil
	}
	escapedContent := regexp.QuoteMeta(state.SearchContent)
	return regexp.Compile(escapedContent)
}

func (state *State) SearchItems() {
	state.ItemsFound = []*Item{}
	rg, err := state.SearchRegexp()
	if state.SearchContent == "" || err != nil {
		for i := range state.Items {
			state.ItemsFound = append(state.ItemsFound, &state.Items[i])
		}
		return
	}
	for i := range state.Items {
		if rg.MatchString(state.Items[i].Name) {
			state.ItemsFound = append(state.ItemsFound, &state.Items[i])
		}
	}
}

func drawText(screen tcell.Screen, x, y int, text string, style tcell.Style) {
	for i, r := range text {
		screen.SetContent(x+i, y, r, nil, style)
	}
}

func redraw(screen tcell.Screen, state State) {
	screen.Clear()
	w, h := screen.Size()

	switch state.Window {
	case StateSearch:
		y := h - 1
		searchLabel := "Search: "
		for i, r := range searchLabel {
			screen.SetContent(i, y, r, nil, tcell.StyleDefault.Foreground(tcell.ColorWhite))
		}
		for i, r := range state.SearchContent {
			screen.SetContent(len(searchLabel)+i, y, r, nil, tcell.StyleDefault.Foreground(tcell.ColorWhite))
		}
		screen.ShowCursor(len(searchLabel)+len(state.SearchContent), y)
	case StateList:
		screen.HideCursor()
		for i, item := range state.ItemsFound {
			style := tcell.StyleDefault.Foreground(tcell.ColorWhite)
			if item.Selected {
				style = style.Background(tcell.ColorBlue)
			}
			drawText(screen, 2, i+2, item.Name, style)
		}
	}

	top := 0
	bottom := h - 1
	if state.Window == StateSearch {
		bottom = h - 2
	}
	for x := range w {
		screen.SetContent(x, top, '─', nil, tcell.StyleDefault.Foreground(tcell.ColorWhite))
		screen.SetContent(x, bottom, '━', nil, tcell.StyleDefault.Foreground(tcell.ColorWhite))
	}
	for y := top; y < h-1; y++ {
		screen.SetContent(0, y, '│', nil, tcell.StyleDefault.Foreground(tcell.ColorWhite))
		screen.SetContent(w-1, y, '│', nil, tcell.StyleDefault.Foreground(tcell.ColorWhite))
	}
	screen.SetContent(0, top, '┌', nil, tcell.StyleDefault.Foreground(tcell.ColorWhite))
	screen.SetContent(w-1, top, '┐', nil, tcell.StyleDefault.Foreground(tcell.ColorWhite))
	screen.SetContent(0, bottom, '┕', nil, tcell.StyleDefault.Foreground(tcell.ColorWhite))
	screen.SetContent(w-1, bottom, '┙', nil, tcell.StyleDefault.Foreground(tcell.ColorWhite))

	screen.Show()
}

func main() {
	screen, err := tcell.NewScreen()
	if err != nil {
		println(err)
	}
	if err := screen.Init(); err != nil {
		println(err)
	}
	defer screen.Fini()

	state := State{Items: []Item{{Name: "a"}, {Name: "b"}, {Name: "c"}, {Name: "d"}}}
	state.SearchItems()
	redraw(screen, state)

	for {
		ev := screen.PollEvent()
		switch tev := ev.(type) {
		case *tcell.EventKey:
			if tev.Key() == tcell.KeyCtrlC {
				screen.Fini()
				os.Exit(0)
			}
			switch state.Window {
			case StateSearch:
				switch tev.Key() {
				case tcell.KeyEsc:
					state.SearchContent = ""
					state.Window = StateList
					redraw(screen, state)
				case tcell.KeyTab, tcell.KeyEnter:
					state.Window = StateList
					redraw(screen, state)
				case tcell.KeyBackspace, tcell.KeyBackspace2:
					if len(state.SearchContent) > 0 {
						state.SearchContent = state.SearchContent[:len(state.SearchContent)-1]
						redraw(screen, state)
					}
				default:
					if tev.Rune() != 0 && tev.Modifiers() == 0 {
						state.SearchContent += string(tev.Rune())
						state.SearchItems()
						redraw(screen, state)
					}
				}
			case StateList:
				if tev.Rune() == 'q' {
					screen.Fini()
					os.Exit(0)
				}
				switch tev.Key() {
				case tcell.KeyEsc, tcell.KeyCtrlC:
					screen.Fini()
					os.Exit(0)
				case tcell.KeyUp:
					newIndex := slices.Index(state.ItemsFound, state.ItemCursor) - 1
					if newIndex >= 0 {
						state.ItemCursor = state.ItemsFound[newIndex]
					}
					redraw(screen, state)
				case tcell.KeyDown:
					newIndex := slices.Index(state.ItemsFound, state.ItemCursor) + 1
					if newIndex <= len(state.ItemsFound) {
						state.ItemCursor = state.ItemsFound[newIndex]
					}
					redraw(screen, state)
				case tcell.KeyTab:
					state.Window = StateSearch
					redraw(screen, state)
				case tcell.KeyEnter:
					// Item selected, handle as needed
				}
			}
		case *tcell.EventResize:
			redraw(screen, state)
		}
	}
}
