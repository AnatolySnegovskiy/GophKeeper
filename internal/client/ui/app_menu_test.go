package ui

import (
	"testing"

	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
)

func TestShowAppMenu(t *testing.T) {
	menu := &Menu{
		app:   tview.NewApplication(),
		title: "Test Title",
	}

	menu.showAppMenu()
	assert.NotNil(t, menu.app)
}
