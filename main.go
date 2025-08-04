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
	StateHelp Window = iota
	StateList
	StateSearch
	StateNew
)

type Item struct {
	Name  string
	Since time.Time
}

type State struct {
	Window        Window `json:"-"`
	Items         []*Item
	ItemsFound    []*Item `json:"-"`
	ItemFound     int     `json:"-"`
	Item          int     `json:"-"`
	SearchContent string  `json:"-"`
	NewItem       Item    `json:"-"`
}

func (state *State) SearchRegexp() (*regexp.Regexp, error) {
	if state.SearchContent == "" {
		return nil, nil
	}
	return regexp.Compile(state.SearchContent)
}

func (state *State) UpdateSelected(foundIndex int) {
	if foundIndex < 0 || foundIndex >= len(state.ItemsFound) {
		return
	}
	if len(state.ItemsFound) > 0 {
		state.ItemFound = foundIndex
		state.Item = slices.Index(state.Items, state.ItemsFound[foundIndex])
	}
}

func (state *State) SearchItems() {
	state.ItemsFound = []*Item{}
	rg, err := state.SearchRegexp()
	state.ItemFound = -1
	state.Item = -1
	if state.SearchContent == "" || err != nil {
		for i := range state.Items {
			state.ItemsFound = append(state.ItemsFound, state.Items[i])
		}
		state.UpdateSelected(0)
		return
	}
	for i := range state.Items {
		if rg.MatchString(state.Items[i].Name) {
			state.ItemsFound = append(state.ItemsFound, state.Items[i])
		}
	}
	state.UpdateSelected(0)
}

const help = "Welcome to timetrack!\n" +
	"This screen available through <?> and <h> keys.\n\n" +
	"Use ▲ and ▼ to move.\n" +
	"<C-c>, <q>, <esc> - quit\n" +
	"<a> - add new record\n" +
	"<d> - delete current record\n" +
	"<C-f>, </> - search records\n" +
	"<D> - delete all found records"

func redraw(screen tcell.Screen, state State) {
	screen.Clear()
	w, h := screen.Size()
	style := tcell.StyleDefault.Foreground(tcell.ColorWhite)

	listBox := tcell.WindowSize{PixelWidth: w - 2, PixelHeight: h - 2, Height: 1, Width: 2}
	if state.Window == StateSearch || state.Window == StateNew {
		listBox.PixelHeight = h - 3
	}

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
	case StateHelp:
		drawParagraph(screen, listBox, help, style)
		screen.HideCursor()
	case StateList:
		screen.HideCursor()
	}

	if len(state.ItemsFound) == 0 {
		drawParagraph(screen, listBox, help, style)
	} else if state.Window != StateHelp {
		drawItems(screen, state.ItemsFound, state.ItemFound, true, listBox)
	}

	for x := range w {
		screen.SetContent(x, listBox.Height-1, '─', nil, tcell.StyleDefault.Foreground(tcell.ColorWhite))
		screen.SetContent(x, listBox.PixelHeight+1, '━', nil, tcell.StyleDefault.Foreground(tcell.ColorWhite))
	}
	for y := listBox.Height; y < h-1; y++ {
		screen.SetContent(0, y, '│', nil, tcell.StyleDefault.Foreground(tcell.ColorWhite))
		screen.SetContent(w-1, y, '│', nil, tcell.StyleDefault.Foreground(tcell.ColorWhite))
	}
	screen.SetContent(0, listBox.Height-1, '┌', nil, tcell.StyleDefault.Foreground(tcell.ColorWhite))
	screen.SetContent(w-1, listBox.Height-1, '┐', nil, tcell.StyleDefault.Foreground(tcell.ColorWhite))
	screen.SetContent(0, listBox.PixelHeight+1, '┕', nil, tcell.StyleDefault.Foreground(tcell.ColorWhite))
	screen.SetContent(w-1, listBox.PixelHeight+1, '┙', nil, tcell.StyleDefault.Foreground(tcell.ColorWhite))

	screen.Show()
}

func main() {
	if slices.Contains(os.Args, "--help") || slices.Contains(os.Args, "-h") {
		println(help)
		os.Exit(1)
		return
	}
	screen, err := tcell.NewScreen()
	if err != nil {
		println(err)
		return
	}
	if err := screen.Init(); err != nil {
		println(err)
		return
	}
	defer screen.Fini()

	state := State{
		Items: []*Item{{Name: "First time opened timetrack", Since: time.Now()}}, // will be erased by load
	}
	load(&state)
	state.Window = StateList
	state.SearchItems()
	redraw(screen, state)

	go func() {
		for {
			time.Sleep(50 * time.Millisecond)
			redraw(screen, state)
		}
	}()

	for {
		handleEvent(screen, &state)
	}
}
