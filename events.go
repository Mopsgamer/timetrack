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
			state.InputSearch.HandleInput(tev)
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
				state.InputNew.HandleInput(tev)
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
				goto Redraw
			case 'q':
				screen.Fini()
				os.Exit(0)
				return
			case '/':
				state.Window = StateSearch
				redraw(screen, *state)
				goto Redraw
			case 'd':
				if len(state.ItemsFound) > 0 {
					state.Items = slices.Delete(state.Items, state.Item, state.Item+1)
					state.SearchItems()
					redraw(screen, *state)
				}
				save(*state)
				goto Redraw
			case 'D':
				for _, item := range state.ItemsFound {
					index := slices.Index(state.Items, item)
					state.Items = slices.Delete(state.Items, index, index+1)
				}
				state.SearchContent = ""
				state.SearchItems()
				save(*state)
				goto Redraw
			case 'a':
				state.Window = StateNew
				redraw(screen, *state)
				goto Redraw
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
				goto Redraw
			case tcell.KeyDown:
				state.UpdateSelected(state.ItemFound + 1)
				goto Redraw
			case tcell.KeyTab, tcell.KeyCtrlF:
				state.Window = StateSearch
				goto Redraw
			}
		Redraw:
			redraw(screen, *state)
			return
		}
	case *tcell.EventResize:
		redraw(screen, *state)
	}
}
