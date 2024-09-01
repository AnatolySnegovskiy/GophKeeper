package ui

import (
	"context"
	"fmt"
	"github.com/rivo/tview"
	v1 "goph_keeper/internal/services/grpc/goph_keeper/v1"
	"path/filepath"
	"strings"
)

func (m *Menu) showFilesMenu() {
	m.logger.Info("Show files menu")
	title := tview.NewTextView().
		SetText("Files").
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true)

	list := tview.NewList()
	listFiles, err := m.grpcClient.GetStoreDataList(context.Background())

	m.errorHandler(err)

	vDirectories := buildVirtualDirectories(listFiles.Entries)
	m.showVirtualDirectoryContents(vDirectories, string(filepath.Separator), list)

	mainLayout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(title, 3, 1, false).
		AddItem(list, 0, 1, true)

	m.app.SetRoot(mainLayout, true).SetFocus(list)
}

func buildVirtualDirectories(entries []*v1.ListDataEntry) []string {
	vDirs := make([]string, len(entries))
	for _, entry := range entries {
		vDirs = append(vDirs, entry.UserPath)
	}

	return vDirs
}

func (m *Menu) showVirtualDirectoryContents(vDirs []string, currentPath string, list *tview.List) {
	list.Clear()
	subDirs := []string{}
	files := []string{}

	parentPath := filepath.Dir(currentPath)
	list.AddItem("..", "", 0, func() {
		m.showVirtualDirectoryContents(vDirs, parentPath, list)
	})

	// Проходим по всем путям в vDirs
	for _, path := range vDirs {
		dir := filepath.Dir(path)
		file := filepath.Base(path)
		m.logger.Info(fmt.Sprintf("dir: %s, file: %s", dir, file))
		m.logger.Info(fmt.Sprintf("currentPath: %s", currentPath))
		if !strings.Contains(dir, currentPath) {
			continue
		}

		if dir == currentPath {
			files = append(files, file)
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
		list.AddItem(pathDir+string(filepath.Separator), "", 0, func() {
			m.showVirtualDirectoryContents(vDirs, filepath.Join(currentPath, pathDir), list)
		})
	}

	// Добавляем файлы в список
	for _, file := range files {
		fileItem := file
		list.AddItem(fileItem, "", 0, func() {
			// Обработка выбора файла (например, открытие или отображение его содержимого)
		})
	}
}
