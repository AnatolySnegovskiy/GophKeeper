package ui

import (
	"context"
	"fmt"
	"github.com/rivo/tview"
	"goph_keeper/internal/client"
	"goph_keeper/internal/services/grpc/goph_keeper/v1"
	"os"
	"path/filepath"
	"regexp"
)

type ApplicationInterface interface {
	QueueUpdateDraw(func()) *tview.Application
	Stop()
	Draw() *tview.Application
}

func (m *Menu) showDownloadFileForm(entry *v1.ListDataEntry, rollbackFilesMenu func()) {
	progressBar := NewProgressBar(100)
	info := tview.NewTextView().
		SetDynamicColors(true).
		SetText("")

	f := func(directoryPath string, rollback func()) {
		if directoryPath == "" {
			info.SetText("Directory Path cannot be empty")
			return
		}

		form := createDownloadForm(info, progressBar)
		m.app.SetRoot(createFlexLayout(form), true)

		progressChan := make(chan int)

		go handleFileDownload(directoryPath, entry, progressChan, info, m.grpcClient, m.app)
		go handleProgressUpdates(progressChan, progressBar, rollbackFilesMenu, form, m.app)
	}

	m.explore(Dir, f)
}

func createDownloadForm(info *tview.TextView, progressBar *ProgressBar) *tview.Form {
	form := tview.NewForm().
		AddFormItem(info).
		AddFormItem(progressBar)
	form.SetBorder(true).SetTitle("Download File").SetTitleAlign(tview.AlignLeft)
	return form
}

func createFlexLayout(form *tview.Form) *tview.Flex {
	return tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(form, 0, 1, true)
}

func handleFileDownload(directoryPath string, entry *v1.ListDataEntry, progressChan chan int, info *tview.TextView, grpcClient *client.GrpcClient, app ApplicationInterface) {
	directoryPath = cleanDirectoryPath(directoryPath)

	fileInfo, err := os.Stat(directoryPath)
	if err != nil {
		app.QueueUpdateDraw(func() {
			info.SetText(fmt.Sprintf("[red]Error: %s", err))
		})
		return
	}

	if !fileInfo.IsDir() {
		directoryPath = filepath.Dir(directoryPath)
	}

	_, err = grpcClient.DownloadFile(context.Background(), entry.Uuid, directoryPath, progressChan)
	if err != nil {
		app.QueueUpdateDraw(func() {
			info.SetText(fmt.Sprintf("[red]Error: %s", err))
		})
		return
	}
	close(progressChan)
	app.QueueUpdateDraw(func() {
		info.SetText(fmt.Sprintf("[green]Success: %t", true))
	})
}

func cleanDirectoryPath(directoryPath string) string {
	re := regexp.MustCompile(`[./]+$`)
	cleanedPath := re.ReplaceAllString(directoryPath, "")

	return cleanedPath
}

func handleProgressUpdates(progressChan chan int, progressBar *ProgressBar, rollbackFilesMenu func(), form *tview.Form, app ApplicationInterface) {
	for progress := range progressChan {
		app.QueueUpdateDraw(func() {
			progressBar.SetProgress(progress)
		})
	}
	if progressBar.current >= 100 {
		form.AddButton("OK", func() {
			rollbackFilesMenu()
		})
	}
}
