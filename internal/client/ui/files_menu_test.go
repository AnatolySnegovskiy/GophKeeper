package ui

import (
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
	"testing"
)

func testShowFilesMenu(t *testing.T) {
	menu := &Menu{
		app:   tview.NewApplication(),
		title: "Test Title",
	}

	menu.showFilesMenu()
	assert.NotNil(t, menu.app)
}
