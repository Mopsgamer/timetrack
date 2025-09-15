package main

import (
	"regexp"
	"slices"
	"time"
)

type State struct {
	Window        Window `json:"-"`
	HelpScroll    int    `json:"-"`
	Items         []*Item
	ItemsFound    []*Item    `json:"-"`
	ItemFound     int        `json:"-"`
	Item          int        `json:"-"`
	SearchContent string     `json:"-"`
	InputSearch   InputState `json:"-"`
	InputNew      InputState `json:"-"`
	InputRename   InputState `json:"-"`
}

func LoadState(state *State) {
	state.Items = []*Item{{Name: "First time opened timetrack", Since: time.Now()}} // will be erased by load
	load(state)

	state.Window = StateList

	state.InputSearch = InputState{
		Value: state.SearchContent,
		OnChange: func(text string) {
			state.SearchContent = text
			state.SearchItems(0)
		},
		OnEscape: func() {
			state.Window = StateList
		},
	}
	state.InputNew = InputState{
		OnProceed: func(text string) {
			newItem := Item{
				Name:  text,
				Since: time.Now(),
			}
			state.Items = append(state.Items, &newItem)
			state.Window = StateList
			state.SearchItems(0)
			save(*state)
		},
		OnEscape: func() {
			state.Window = StateList
		},
	}
	state.InputRename = InputState{
		OnProceed: func(text string) {
			state.Items[state.Item].Name = state.InputRename.Value
			save(*state)
		},
		OnEscape: func() {
			state.InputRename.Value = ""
			state.InputRename.Cursor = 0
			state.Window = StateList
		},
	}

	state.SearchItems(0)
}

func (state *State) SearchRegexp() (*regexp.Regexp, error) {
	if state.SearchContent == "" {
		return nil, nil
	}
	return regexp.Compile(state.SearchContent)
}

func (state *State) UpdateSelected(foundIndex int) {
	foundIndex = max(0, min(foundIndex, len(state.ItemsFound)-1))
	if len(state.ItemsFound) > 0 {
		state.ItemFound = foundIndex
		state.Item = slices.Index(state.Items, state.ItemsFound[foundIndex])
	}
}

func (state *State) SearchItems(foundItem int) {
	state.ItemsFound = []*Item{}
	rg, err := state.SearchRegexp()
	state.ItemFound = -1
	state.Item = -1
	if state.SearchContent == "" || err != nil {
		for i := range state.Items {
			state.ItemsFound = append(state.ItemsFound, state.Items[i])
		}
		state.UpdateSelected(foundItem)
		return
	}
	for i := range state.Items {
		if rg.MatchString(state.Items[i].Name) {
			state.ItemsFound = append(state.ItemsFound, state.Items[i])
		}
	}
	state.UpdateSelected(foundItem)
}
