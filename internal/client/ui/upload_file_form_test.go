package ui

import (
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestShowSendFileForm(t *testing.T) {
	menu := &Menu{
		app:   tview.NewApplication(),
		title: "Test Title",
	}
	menu.showSendFileForm()
	assert.NotNil(t, menu.app)
}
