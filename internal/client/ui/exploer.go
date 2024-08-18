package ui

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func (m *Menu) explore(callback func(filePath string, rollback func())) {
	list := tview.NewList().ShowSecondaryText(false)

	getDrives := func() []string {
		if runtime.GOOS != "windows" {
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
		drive := drive // захват переменной
		list.AddItem(drive, "", 0, func() {
			m.showDirectoryContents(drive, list, callback)
		})
	}

	m.app.SetRoot(list, true)
}

func (m *Menu) showDirectoryContents(path string, list *tview.List, callback func(filePath string, rollback func())) {
	list.Clear()
	files, err := os.ReadDir(path)
	if err != nil {
		list.AddItem(fmt.Sprintf("Error: %v", err), "", 0, nil)
		return
	}

	if isDriveRoot(path) {
		list.AddItem("..", "", 0, func() {
			m.explore(callback)
		})
	} else if !isRootPath(path) {
		parentPath := filepath.Dir(path)
		list.AddItem("..", "", 0, func() {
			m.showDirectoryContents(parentPath, list, callback)
		})
	}

	for _, file := range files {
		fileName := file.Name()
		filePath := filepath.Join(path, fileName)
		if file.IsDir() {
			list.AddItem(fileName+"/", "", 0, func() {
				m.showDirectoryContents(filePath, list, callback)
			})
		} else {
			list.AddItem(fileName, "", 0, func() {
				callback(filePath, func() {
					m.app.SetRoot(list, true)
				})
			})
		}
	}

	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRight:
			currentItemIndex := list.GetCurrentItem()
			if currentItemIndex >= 0 && currentItemIndex < list.GetItemCount() {
				mainText, _ := list.GetItemText(currentItemIndex)
				if strings.HasSuffix(mainText, "/") {
					m.showDirectoryContents(filepath.Join(path, strings.TrimSuffix(mainText, "/")), list, callback)
				}
			}
		case tcell.KeyLeft:
			parentPath := filepath.Dir(path)
			m.showDirectoryContents(parentPath, list, callback)
		default:
			return event
		}
		return event
	})
}
