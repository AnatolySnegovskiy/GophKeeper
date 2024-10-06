package ui

import (
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestShowRegistrationForm(t *testing.T) {
	menu := &Menu{
		app:   tview.NewApplication(),
		title: "Test Title",
	}

	menu.showRegistrationForm()
	assert.NotNil(t, menu.app)
}

func TestShowAuthorizationForm(t *testing.T) {
	menu := &Menu{
		app:   tview.NewApplication(),
		title: "Test Title",
	}

	menu.showAuthorizationForm()
	assert.NotNil(t, menu.app)
}
