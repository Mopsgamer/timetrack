package main

type Window uint

const (
	StateHelp Window = iota
	StateList
	StateSearch
	StateNew
	StateRename
)
