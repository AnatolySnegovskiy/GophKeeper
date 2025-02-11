package ui

import (
	"context"
	"fmt"
	v1 "goph_keeper/internal/services/grpc/goph_keeper/v1"
	"path/filepath"
	"strings"

	"github.com/rivo/tview"
)

func (m *Menu) showServerFilesMenu(currentPath string) {
	m.logger.Info("Show files menu")
	title := tview.NewTextView().
		SetText("Files").
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true)

	list := tview.NewList()
	listFiles, err := m.grpcClient.GetStoreDataList(context.Background(), v1.DataType_DATA_TYPE_BINARY)
	if err != nil {
		m.errorHandler(err, func() {
			m.showServerFilesMenu(currentPath)
		})
		return
	}

	vDirectories := buildVirtualDirectories(listFiles.Entries)
	m.showVirtualDirectoryContents(vDirectories, currentPath, list)

	mainLayout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(title, 3, 1, false).
		AddItem(list, 0, 1, true)

	list.AddItem("Back", "", 0, func() {
		m.showFilesMenu()
	})

	m.app.SetRoot(mainLayout, true).SetFocus(list)
}

func buildVirtualDirectories(entries []*v1.ListDataEntry) map[string]*v1.ListDataEntry {
	vDirs := map[string]*v1.ListDataEntry{}

	for _, entry := range entries {
		vDirs[entry.UserPath] = entry
	}

	return vDirs
}

func (m *Menu) showVirtualDirectoryContents(vDirs map[string]*v1.ListDataEntry, currentPath string, viewList *tview.List) {
	viewList.Clear()
	subDirs := []string{}
	files := map[string]*v1.ListDataEntry{}

	parentPath := filepath.Dir(currentPath)
	viewList.AddItem("..", "", 0, func() {
		m.showVirtualDirectoryContents(vDirs, parentPath, viewList)
		if parentPath == string(filepath.Separator) {
			viewList.AddItem("Back", "", 0, func() {
				m.showFilesMenu()
			})
		}
	})

	// Проходим по всем путям в vDirs
	for path, entry := range vDirs {
		dir := filepath.Dir(path)
		file := filepath.Base(path)
		m.logger.Info(fmt.Sprintf("dir: %s, file: %s", dir, file))
		m.logger.Info(fmt.Sprintf("currentPath: %s", currentPath))
		if !strings.Contains(dir, currentPath) {
			continue
		}

		if dir == currentPath {
			files[file] = entry
		} else {
			if currentPath != string(filepath.Separator) {
				dir = strings.ReplaceAll(dir, currentPath, "")
			}
			parts := strings.Split(dir, string(filepath.Separator))
			if len(parts) > 1 {
				subDirs = append(subDirs, parts[1])
			}
		}
	}

	m.logger.Info(fmt.Sprintf("currentPath: %s, subDirs: %s, files: %s", currentPath, subDirs, files))

	for _, dir := range subDirs {
		m.logger.Info(fmt.Sprintf("currentPath: %s, dir: %s", currentPath, dir))
		pathDir := dir
		viewList.AddItem(pathDir+string(filepath.Separator), "", 0, func() {
			m.showVirtualDirectoryContents(vDirs, filepath.Join(currentPath, pathDir), viewList)
		})
	}

	for name, e := range files {
		fileName := name
		viewList.AddItem(fileName, "", 0, func() {
			modal := tview.NewModal().SetText(fmt.Sprintf("Download File %s?", fileName)).AddButtons([]string{"Yes", "No"}).SetDoneFunc(func(buttonIndex int, buttonLabel string) {
				if buttonLabel == "Yes" {
					m.showDownloadFileForm(e, func() {
						m.showServerFilesMenu(currentPath)
					})
				} else if buttonLabel == "No" {
					m.showServerFilesMenu(currentPath)
				}
			})
			m.app.SetRoot(modal, true)
		})
	}
}
