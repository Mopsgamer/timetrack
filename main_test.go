package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTimet(t *testing.T) {
	check := assert.New(t)
	check.NotPanics(func() {
		ProcessAction("message", nil)
	})
	check.Panics(func() {
		ProcessAction("", errors.New("testing"))
	})
	check.NotPanics(func() {
		ProcessAction("", nil)
	})
	check.NotPanics(func() {
		main()
	})
}
