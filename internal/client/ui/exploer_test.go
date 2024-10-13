package ui

import (
	"github.com/rivo/tview"
	"testing"
)

func TestExplore(t *testing.T) {
	menu := &Menu{
		app:   tview.NewApplication(),
		title: "Test Title",
	}
	callback := func(filePath string, rollback func()) {
	}
	menu.explore(Dir, callback)
}

func TestShowDirectoryContents(t *testing.T) {
	menu := &Menu{
		app:   tview.NewApplication(),
		title: "Test Title",
	}
	callback := func(filePath string, rollback func()) {
	}
	menu.showDirectoryContents(&explore{list: tview.NewList()}, callback)
}
