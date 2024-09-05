package ui

import (
	"context"
	"fmt"
	"github.com/rivo/tview"
	"goph_keeper/internal/services/grpc/goph_keeper/v1"
)

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

		form := tview.NewForm().
			AddFormItem(info).
			AddFormItem(progressBar)
		form.SetBorder(true).SetTitle("Download File").SetTitleAlign(tview.AlignLeft)
		flex := tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(form, 0, 1, true)

		m.app.SetRoot(flex, true)

		progressChan := make(chan int)

		go func() {
			err := m.grpcClient.DownloadFile(context.Background(), entry.Uuid, directoryPath, progressChan)
			if err != nil {
				m.app.QueueUpdateDraw(func() {
					info.SetText(fmt.Sprintf("[red]Error: %s", err))
				})
			}
			close(progressChan)
			m.app.QueueUpdateDraw(func() {
				info.SetText(fmt.Sprintf("[green]Success: %t", true))
			})
		}()

		go func() {
			for progress := range progressChan {
				m.app.QueueUpdateDraw(func() {
					progressBar.SetProgress(progress)
				})
			}
			if progressBar.current >= 100 {
				form.AddButton("OK", func() {
					rollbackFilesMenu()
				})
			}
		}()
	}

	m.explore(Dir, f)
}
