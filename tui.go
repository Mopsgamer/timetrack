package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
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

func (state *State) SearchItems() {
	state.ItemsFound = []*Item{}
	rg, err := state.SearchRegexp()
	state.ItemFound = -1
	state.Item = -1
	if state.SearchContent == "" || err != nil {
		for i := range state.Items {
			state.ItemsFound = append(state.ItemsFound, state.Items[i])
		}
		goto SetCurrent
	}
	for i := range state.Items {
		if rg.MatchString(state.Items[i].Name) {
			state.ItemsFound = append(state.ItemsFound, state.Items[i])
		}
	}
SetCurrent:
	if len(state.ItemsFound) > 0 {
		state.ItemFound = 0
	}
	state.Item = slices.Index(state.Items, state.ItemsFound[state.ItemFound])
}

var homedir, _ = os.UserHomeDir()
var path = filepath.Join(homedir, ".timetrack.json")

func save(state State) {
	bytes, err := json.Marshal(state)
	if err != nil {
		panic(err)
	}
	os.WriteFile(path, bytes, 0666)
}

func load(state *State) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		panic(err)
	}
	if err := json.Unmarshal(bytes, state); err != nil {
		panic(err)
	}
	state.SearchItems()
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

	top := 0
	bottom := h - 1
	if state.Window == StateSearch || state.Window == StateNew {
		bottom = h - 2
	}

	drawItems(screen, state.ItemsFound, state.ItemFound, true, tcell.WindowSize{PixelWidth: w - 2, PixelHeight: bottom - 1, Height: 1, Width: 2})

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

func drawItems(screen tcell.Screen, items []*Item, current int, highlight bool, size tcell.WindowSize) {
	half := size.PixelHeight / 2
	from, to := 0, len(items)
	if current > half {
		if current > len(items)-half {
			from, to = len(items)-size.PixelHeight, len(items)
		} else {
			from, to = current-half-1, current+half
		}
	}

	for i, item := range items[from:to] {
		style := tcell.StyleDefault.Foreground(tcell.ColorWhite)
		since := time.Since(item.Since).Round(time.Second).String()
		if highlight && item == items[current] {
			style = tcell.StyleDefault.Foreground(tcell.ColorBlack).Background(tcell.ColorGhostWhite)
			drawText(screen, size.Width+len(item.Name), i+size.Height, strings.Repeat(" ", size.PixelWidth-size.Width-len(item.Name)-len(since)), style)
		}
		drawText(screen, size.Width, i+size.Height, item.Name, style)
		drawTextRight(screen, size.PixelWidth-size.Width+1, i+size.Height, since, style)
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
	load(&state)
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
				handleInput(tev, state.SearchContent,
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
					newItem := state.NewItem
					state.Items = append(state.Items, &newItem)
					state.Window = StateList
					state.SearchItems()
					redraw(screen, state)
					state.NewItem.Name = ""
					save(state)
				default:
					handleInput(tev, state.NewItem.Name,
						func(text string) {
							state.NewItem.Name = text
						},
						func() {
							state.Window = StateList
						},
					)
				}
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
						state.Items = slices.Delete(state.Items, state.Item, state.Item+1)
						state.SearchItems()
						redraw(screen, state)
					}
					save(state)
				case 'D':
					for _, item := range state.ItemsFound {
						index := slices.Index(state.Items, item)
						state.Items = slices.Delete(state.Items, index, index+1)
					}
					state.SearchContent = ""
					state.SearchItems()
					save(state)
				case 'a':
					state.Window = StateNew
				}
				switch tev.Key() {
				case tcell.KeyEsc, tcell.KeyCtrlC:
					screen.Fini()
					os.Exit(0)
				case tcell.KeyUp:
					newIndex := state.ItemFound - 1
					if newIndex >= 0 {
						state.ItemFound = newIndex
					}
					redraw(screen, state)
				case tcell.KeyDown:
					newIndex := state.ItemFound + 1
					if newIndex < len(state.ItemsFound) {
						state.ItemFound = newIndex
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
