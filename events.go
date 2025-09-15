package main

import (
	"os"
	"slices"
	"time"

	"github.com/gdamore/tcell/v2"
)

func handleEvent(screen tcell.Screen, state *State) {
	ev := screen.PollEvent()
	_, h := screen.Size()
	switch tev := ev.(type) {
	case *tcell.EventKey:
		if tev.Key() == tcell.KeyCtrlC {
			screen.Fini()
			os.Exit(0)
			return
		}
		switch state.Window {
		case StateSearch:
			state.InputSearch.HandleInput(tev)
			goto Redraw
		case StateNew:
			state.InputNew.HandleInput(tev)
			goto Redraw
		case StateRename:
			state.InputRename.HandleInput(tev)
			goto Redraw
		case StateHelp:
			switch tev.Rune() {
			case 'h', '?', 'q':
				state.Window = StateList
				goto Redraw
			}
			switch tev.Key() {
			case tcell.KeyEsc:
				state.Window = StateList
				goto Redraw
			case tcell.KeyPgUp:
				state.HelpScroll = max(0, state.HelpScroll-h)
				goto Redraw
			case tcell.KeyPgDn:
				lines := len(splitTextIntoLines(help, listBox.PixelWidth))
				state.HelpScroll = min(state.HelpScroll+h, max(0, lines-listBox.PixelHeight))
				goto Redraw
			case tcell.KeyUp:
				state.HelpScroll = max(0, state.HelpScroll-1)
				goto Redraw
			case tcell.KeyDown:
				lines := len(splitTextIntoLines(help, listBox.PixelWidth))
				state.HelpScroll = min(state.HelpScroll+1, max(0, lines-listBox.PixelHeight))
				goto Redraw
			}
		case StateList:
			switch tev.Rune() {
			case 'h', '?':
				state.Window = StateHelp
				goto Redraw
			case 'q':
				screen.Fini()
				os.Exit(0)
				return
			case '/':
				state.Window = StateSearch
				goto Redraw
			case 'a':
				state.Window = StateNew
				goto Redraw
			case 'r':
				state.Items[state.Item].Since = time.Now()
				save(*state)
				goto Redraw
			case 'd':
				if len(state.ItemsFound) > 0 {
					state.Items = slices.Delete(state.Items, state.Item, state.Item+1)
				}
				state.SearchItems(state.ItemFound - 1)
				save(*state)
				goto Redraw
			case 'D':
				for _, item := range state.ItemsFound {
					index := slices.Index(state.Items, item)
					state.Items = slices.Delete(state.Items, index, index+1)
				}
				state.SearchContent = ""
				state.SearchItems(0)
				save(*state)
				goto Redraw
			case 'A':
				state.Window = StateRename
				goto Redraw
			}
			switch tev.Key() {
			case tcell.KeyEsc:
				screen.Fini()
				os.Exit(0)
				return
			case tcell.KeyPgUp:
				state.UpdateSelected(state.ItemFound - h)
				goto Redraw
			case tcell.KeyPgDn:
				state.UpdateSelected(state.ItemFound + h)
				goto Redraw
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
		}
	case *tcell.EventResize:
		goto Redraw
	default:
		return
	}
Redraw:
	redraw(screen, state)
}
