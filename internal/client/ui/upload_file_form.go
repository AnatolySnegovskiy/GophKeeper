package ui

import (
	"context"
	"fmt"
	v1 "goph_keeper/internal/services/grpc/goph_keeper/v1"
	"strings"

	"github.com/rivo/tview"
)

func (m *Menu) showSendFileForm() {
	progressBar := NewProgressBar(100)
	info := tview.NewTextView().
		SetDynamicColors(true).
		SetText("")

	f := func(filePath string, rollback func()) {
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

			if !strings.HasPrefix(userPath, "/") {
				userPath = "/" + userPath
			}

			userPath = strings.TrimSuffix(userPath, "/")
			progressChan := make(chan int)

			go func() {
				res, err := m.grpcClient.UploadFile(context.Background(), filePath, userPath, v1.DataType_DATA_TYPE_BINARY, progressChan)
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
