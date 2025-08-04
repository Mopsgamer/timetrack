package main

import (
	"os"
	"regexp"
	"slices"
	"time"

	"github.com/gdamore/tcell/v2"
)

type Window uint

const (
	StateList Window = iota
	StateSearch
	StateNew
)

type Item struct {
	Name     string
	Selected bool
	Since    time.Time
}

type State struct {
	Window        Window
	Items         []Item
	ItemsFound    []*Item
	ItemCursor    *Item
	SearchContent string
	NewItem       Item
}

func (state *State) SearchRegexp() (*regexp.Regexp, error) {
	if state.SearchContent == "" {
		return nil, nil
	}
	return regexp.Compile(state.SearchContent)
}

func (state *State) SearchItems() {
	state.ItemsFound = []*Item{}
	rg, err := state.SearchRegexp()
	state.ItemCursor = nil
	if state.SearchContent == "" || err != nil {
		for i := range state.Items {
			state.ItemsFound = append(state.ItemsFound, &state.Items[i])
		}
		goto SetCursor
	}
	for i := range state.Items {
		if rg.MatchString(state.Items[i].Name) {
			state.ItemsFound = append(state.ItemsFound, &state.Items[i])
		}
	}
SetCursor:
	if len(state.ItemsFound) > 0 {
		state.ItemCursor = state.ItemsFound[0]
	}
}

func drawText(screen tcell.Screen, x, y int, text string, style tcell.Style) {
	for i, r := range text {
		screen.SetContent(x+i, y, r, nil, style)
	}
}

func drawTextRight(screen tcell.Screen, x, y int, text string, style tcell.Style) {
	drawText(screen, x-len(text)+1, y, text, style)
}

func redraw(screen tcell.Screen, state State) {
	screen.Clear()
	w, h := screen.Size()
	style := tcell.StyleDefault.Foreground(tcell.ColorWhite)

	switch state.Window {
	case StateSearch:
		x := 0
		y := h - 1
		searchLabel := "Search: /"
		drawText(screen, x, y, searchLabel, style)
		x = len(searchLabel) - 1
		if _, err := state.SearchRegexp(); err != nil {
			style = style.Foreground(tcell.ColorDarkRed)
		}
		drawText(screen, x, y, "/"+state.SearchContent+"/", style)
		x = len(searchLabel) + len(state.SearchContent)
		screen.ShowCursor(x, y)
	case StateNew:
		x := 0
		y := h - 1
		text := "New: " + state.NewItem.Name
		drawText(screen, x, y, text, style)
		x = len(text)
		screen.ShowCursor(x, y)
	case StateList:
		screen.HideCursor()
	}
	for i, item := range state.ItemsFound {
		style := tcell.StyleDefault.Foreground(tcell.ColorWhite)
		if item.Selected {
			style = tcell.StyleDefault.Foreground(tcell.ColorBlack).Background(tcell.ColorGhostWhite)
		}
		if item == state.ItemCursor {
			style = style.Underline(true)
		}
		drawText(screen, 2, i+1, item.Name, style)
		drawTextRight(screen, w-1, i+1, time.Since(item.Since).Round(time.Second).String(), style)
	}

	top := 0
	bottom := h - 1
	if state.Window == StateSearch || state.Window == StateNew {
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

func handleInput(screen tcell.Screen, str string, tev *tcell.EventKey, onInput func(text string), onEscape func()) {
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

func main() {
	screen, err := tcell.NewScreen()
	if err != nil {
		println(err)
	}
	if err := screen.Init(); err != nil {
		println(err)
	}
	defer screen.Fini()

	state := State{}
	state.SearchItems()
	redraw(screen, state)

	go func() {
		for {
			time.Sleep(50 * time.Millisecond)
			redraw(screen, state)
		}
	}()

	for {
		ev := screen.PollEvent()
		switch tev := ev.(type) {
		case *tcell.EventKey:
			switch tev.Key() {
			case tcell.KeyCtrlC:
				screen.Fini()
				os.Exit(0)
			}
			switch state.Window {
			case StateSearch:
				handleInput(screen, state.SearchContent, tev,
					func(text string) {
						state.SearchContent = text
						state.SearchItems()
						redraw(screen, state)
					},
					func() {
						state.Window = StateList
					},
				)
			case StateNew:
				switch tev.Key() {
				case tcell.KeyEnter:
					state.NewItem.Since = time.Now()
					state.Items = append(state.Items, state.NewItem)
					state.SearchItems()
					redraw(screen, state)
					state.NewItem.Name = ""
				}
				handleInput(screen, state.NewItem.Name, tev,
					func(text string) {
						state.NewItem.Name = text
					},
					func() {
						state.Window = StateList
					},
				)
			case StateList:
				switch tev.Rune() {
				case 'q':
					screen.Fini()
					os.Exit(0)
				case '/':
					state.Window = StateSearch
					redraw(screen, state)
				case 'd':
					if len(state.ItemsFound) > 0 {
						index := slices.Index(state.Items, *state.ItemCursor)
						state.Items = slices.Delete(state.Items, index, index+1)
						state.SearchItems()
						redraw(screen, state)
					}
				case 'a':
					state.Window = StateNew
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
					if newIndex < len(state.ItemsFound) {
						state.ItemCursor = state.ItemsFound[newIndex]
					}
					redraw(screen, state)
				case tcell.KeyTab, tcell.KeyCtrlF:
					state.Window = StateSearch
					redraw(screen, state)
				}
			}
		case *tcell.EventResize:
			redraw(screen, state)
		}
	}
}
