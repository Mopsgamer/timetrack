package main

import (
	"regexp"
	"slices"
	"time"
)

type State struct {
	Window        Window `json:"-"`
	Items         []*Item
	ItemsFound    []*Item    `json:"-"`
	ItemFound     int        `json:"-"`
	Item          int        `json:"-"`
	SearchContent string     `json:"-"`
	NewItem       Item       `json:"-"`
	InputSearch   InputState `json:"-"`
	InputNew      InputState `json:"-"`
}

func LoadState(state *State) {
	state.Items = []*Item{{Name: "First time opened timetrack", Since: time.Now()}} // will be erased by load
	load(state)

	state.Window = StateList

	state.InputSearch = InputState{
		Value: state.SearchContent,
		OnInput: func(text string) {
			state.SearchContent = text
			state.SearchItems()
		},
		OnEscape: func() {
			state.Window = StateList
		},
	}
	state.InputNew = InputState{
		Value: state.SearchContent,
		OnInput: func(text string) {
			state.NewItem.Name = text
		},
		OnEscape: func() {
			state.Window = StateList
		},
	}

	state.SearchItems()
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
