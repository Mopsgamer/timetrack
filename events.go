package main

import (
	"os"
	"slices"
	"time"

	"github.com/gdamore/tcell/v2"
)

func handleEvent(screen tcell.Screen, state *State) {
	ev := screen.PollEvent()
	switch tev := ev.(type) {
	case *tcell.EventKey:
		switch tev.Key() {
		case tcell.KeyCtrlC:
			screen.Fini()
			os.Exit(0)
			return
		}
		switch state.Window {
		case StateSearch:
			handleInput(tev, state.SearchContent,
				func(text string) {
					state.SearchContent = text
					state.SearchItems()
					redraw(screen, *state)
				},
				func() {
					state.Window = StateList
				},
			)
			return
		case StateNew:
			switch tev.Key() {
			case tcell.KeyEnter:
				state.NewItem.Since = time.Now()
				newItem := state.NewItem
				state.Items = append(state.Items, &newItem)
				state.Window = StateList
				state.SearchItems()
				redraw(screen, *state)
				state.NewItem.Name = ""
				save(*state)
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
			return
		case StateList, StateHelp:
			switch tev.Rune() {
			case 'h', '?':
				switch state.Window {
				case StateList:
					state.Window = StateHelp
				case StateHelp:
					state.Window = StateList
				}
				redraw(screen, *state)
				return
			case 'q':
				screen.Fini()
				os.Exit(0)
				return
			case '/':
				state.Window = StateSearch
				redraw(screen, *state)
				return
			case 'd':
				if len(state.ItemsFound) > 0 {
					state.Items = slices.Delete(state.Items, state.Item, state.Item+1)
					state.SearchItems()
					redraw(screen, *state)
				}
				save(*state)
				return
			case 'D':
				for _, item := range state.ItemsFound {
					index := slices.Index(state.Items, item)
					state.Items = slices.Delete(state.Items, index, index+1)
				}
				state.SearchContent = ""
				state.SearchItems()
				save(*state)
				return
			case 'a':
				state.Window = StateNew
				redraw(screen, *state)
				return
			}
			switch tev.Key() {
			case tcell.KeyEsc, tcell.KeyCtrlC:
				switch state.Window {
				case StateList:
					screen.Fini()
					os.Exit(0)
				case StateHelp:
					state.Window = StateList
				}
				return
			case tcell.KeyUp:
				state.UpdateSelected(state.ItemFound - 1)
			case tcell.KeyDown:
				state.UpdateSelected(state.ItemFound + 1)
			case tcell.KeyTab, tcell.KeyCtrlF:
				state.Window = StateSearch
			}
			redraw(screen, *state)
			return
		}
	case *tcell.EventResize:
		redraw(screen, *state)
	}
}
