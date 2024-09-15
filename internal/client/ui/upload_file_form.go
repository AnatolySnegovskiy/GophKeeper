package ui

import (
	"context"
	"fmt"
	"github.com/rivo/tview"
)

func (m *Menu) showSendFileForm() {
	progressBar := NewProgressBar(100)
	info := tview.NewTextView().
		SetDynamicColors(true).
		SetText("")

	f := func(filePath string, rollback func()) {
		if filePath == "" {
			info.SetText("File Path cannot be empty")
			return
		}

		userPathField := tview.NewInputField().
			SetLabel("User Path:").
			SetFieldWidth(40)

		form := tview.NewForm().
			AddFormItem(info).
			AddFormItem(userPathField).
			AddFormItem(progressBar)

		form.AddButton("Upload", func() {
			userPath := userPathField.GetText()

			if userPath == "" {
				info.SetText("User Path cannot be empty")
				return
			}

			progressChan := make(chan int)

			go func() {
				res, err := m.grpcClient.UploadFile(context.Background(), filePath, userPath, progressChan)
				if err != nil {
					m.app.QueueUpdateDraw(func() {
						info.SetText(fmt.Sprintf("[red]Error: %s", err))
					})
				}
				close(progressChan)
				m.app.QueueUpdateDraw(func() {
					info.SetText(fmt.Sprintf("[green]Success: %t", res.Success))
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
						rollback()
					})
				}
			}()
		})

		form.SetBorder(true).SetTitle("Send File").SetTitleAlign(tview.AlignLeft)
		flex := tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(form, 0, 1, true)

		m.app.SetRoot(flex, true)
	}

	m.explore(File, f)
}
