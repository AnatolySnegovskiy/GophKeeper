package ui

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"os"
	"path/filepath"
	"strings"
)

type TypeSelection int

const (
	File TypeSelection = iota
	Dir  TypeSelection = iota
)

type explore struct {
	list          *tview.List
	currentPath   string
	typeSelection TypeSelection
}

func (m *Menu) explore(typeSelection TypeSelection, callback func(filePath string, rollback func())) {
	exp := &explore{
		list:          tview.NewList().ShowSecondaryText(false),
		currentPath:   "",
		typeSelection: typeSelection,
	}
	getDrives := func() []string {
		if GetGOOS() != "windows" {
			return []string{"/"}
		}
		var drives []string
		for _, drive := range "ABCDEFGHIJKLMNOPQRSTUVWXYZ" {
			_, err := os.Stat(string(drive) + ":\\")
			if err == nil {
				drives = append(drives, string(drive)+":\\")
			}
		}
		return drives
	}

	drives := getDrives()
	for _, drive := range drives {
		drive := drive
		exp.list.AddItem(drive, "", 0, func() {
			exp.currentPath = drive
			m.showDirectoryContents(exp, callback)
		})
	}
	exp.list.AddItem("Back", "", 0, func() {
		m.showFilesMenu()
	})
	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(exp.list, 0, 1, true)

	var selectButton *tview.Button
	if typeSelection == Dir {
		selectButton = tview.NewButton("Select Directory").SetSelectedFunc(func() {
			currentItemIndex := exp.list.GetCurrentItem()
			if currentItemIndex >= 0 && currentItemIndex < exp.list.GetItemCount() {
				mainText, _ := exp.list.GetItemText(currentItemIndex)
				pathDIr := filepath.Join(exp.currentPath, mainText)
				callback(pathDIr, func() {
					m.app.SetRoot(exp.list, true)
				})
			}
		})
		layout.AddItem(selectButton, 1, 1, false)
	}

	bindingKeyMap(layout, selectButton, m, exp, callback)

	m.app.SetRoot(layout, true).SetFocus(exp.list)
}

func (m *Menu) showDirectoryContents(exp *explore, callback func(filePath string, rollback func())) {
	exp.list.Clear()
	files, err := os.ReadDir(exp.currentPath)
	if err != nil {
		exp.list.AddItem(fmt.Sprintf("Error: %v", err), "", 0, nil)
		return
	}

	if isDriveRoot(exp.currentPath) {
		exp.list.AddItem("..", "", 0, func() {
			m.explore(exp.typeSelection, callback)
		})
	} else if !isRootPath(exp.currentPath) {
		exp.list.AddItem("..", "", 0, func() {
			exp.currentPath = filepath.Dir(exp.currentPath)
			m.showDirectoryContents(exp, callback)
		})
	}

	for _, file := range files {
		fileName := file.Name()
		if file.IsDir() {
			exp.list.AddItem(fileName+"/", "", 0, func() {
				exp.currentPath = filepath.Join(exp.currentPath, fileName)
				m.showDirectoryContents(exp, callback)
			})
		} else if exp.typeSelection == File {
			exp.list.AddItem(fileName, "", 0, func() {
				exp.currentPath = filepath.Join(exp.currentPath, fileName)
				callback(exp.currentPath, func() {
					m.app.SetRoot(exp.list, true)
				})
			})
		}
	}

	exp.list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRight:
			currentItemIndex := exp.list.GetCurrentItem()
			if currentItemIndex >= 0 && currentItemIndex < exp.list.GetItemCount() {
				mainText, _ := exp.list.GetItemText(currentItemIndex)
				if strings.HasSuffix(mainText, "/") {
					exp.currentPath = filepath.Join(exp.currentPath, strings.TrimSuffix(mainText, "/"))
					m.showDirectoryContents(exp, callback)
				}
			}
		case tcell.KeyLeft:
			exp.currentPath = filepath.Dir(exp.currentPath)
			m.showDirectoryContents(exp, callback)
		default:
			return event
		}
		return event
	})
}

func bindingKeyMap(layout *tview.Flex, selectButton *tview.Button, m *Menu, exp *explore, callback func(filePath string, rollback func())) {
	layout.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:
			if selectButton != nil {
				if m.app.GetFocus() == exp.list {
					m.logger.Info("Switching focus to selectButton")
					m.app.SetFocus(selectButton)
				} else {
					m.logger.Info("Switching focus to exp.list")
					m.app.SetFocus(exp.list)
				}
			}
			return nil
		default:
			return event
		}
	})

	exp.list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRight:
			currentItemIndex := exp.list.GetCurrentItem()
			if currentItemIndex >= 0 && currentItemIndex < exp.list.GetItemCount() {
				mainText, _ := exp.list.GetItemText(currentItemIndex)
				if strings.HasSuffix(mainText, "/") {
					exp.currentPath = filepath.Join(exp.currentPath, strings.TrimSuffix(mainText, "/"))
					m.showDirectoryContents(exp, callback)
				}
			}
		case tcell.KeyLeft:
			exp.currentPath = filepath.Dir(exp.currentPath)
			m.showDirectoryContents(exp, callback)
		case tcell.KeyTab:
			if selectButton != nil {
				m.logger.Info("Switching focus to selectButton from exp.list")
				m.app.SetFocus(selectButton)
			}
			return nil
		default:
			return event
		}
		return event
	})
}
