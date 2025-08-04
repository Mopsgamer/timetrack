package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

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
