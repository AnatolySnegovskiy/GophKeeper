package ui

import (
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestShowMainMenu(t *testing.T) {
	menu := &Menu{
		app:   tview.NewApplication(),
		title: "Test Title",
	}

	menu.ShowMainMenu()
	assert.NotNil(t, menu.app)
}
