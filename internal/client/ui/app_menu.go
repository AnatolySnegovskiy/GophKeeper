package ui

import (
	"context"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"log/slog"
)

func (m *Menu) showAppMenu(app *tview.Application, logger *slog.Logger) {
	title := tview.NewTextView().
		SetText(m.title).
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true)

	list := tview.NewList().
		AddItem("1. Отправить файл", "", '1', func() {
			m.showSendFileForm(app, logger)
		}).
		AddItem("2. Получить файл", "", '2', func() {
		}).
		AddItem("3. Синхронизировать файлы", "", '2', func() {
		}).
		SetSelectedFocusOnly(true)

	mainLayout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(title, 3, 1, false).
		AddItem(list, 0, 1, true)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc:
			app.Stop()
			return nil
		}
		return event
	})

	app.SetRoot(mainLayout, true).SetFocus(list)
}

func (m *Menu) showSendFileForm(app *tview.Application, logger *slog.Logger) {
	form := tview.NewForm()
	info := tview.NewTextView().
		SetDynamicColors(true).
		SetText("")
	progressBar := NewProgressBar(100)

	sendButtonHandler := func() {
		filePath := form.GetFormItemByLabel("File Path").(*tview.InputField).GetText()
		progressChan := make(chan int)

		if filePath == "" {
			info.SetText("File Path cannot be empty")
			return
		}

		go func() {
			res, err := m.grpcClient.StoreData(context.Background(), filePath, progressChan)
			if err != nil {
				app.QueueUpdateDraw(func() {
					info.SetText(fmt.Sprintf("[red]Error: %s", err))
				})
			}
			close(progressChan)
			app.QueueUpdateDraw(func() {
				info.SetText(fmt.Sprintf("[green]Success: %s", res.Message))
			})
		}()

		go func() {
			for progress := range progressChan {
				app.QueueUpdateDraw(func() {
					progressBar.SetProgress(progress)
				})
			}
		}()
	}

	form.AddFormItem(info).
		AddInputField("File Path", "", 30, nil, nil).
		AddButton("Send", sendButtonHandler).
		AddButton("Cancel", func() { m.showAppMenu(app, logger) }).
		AddFormItem(progressBar)

	form.SetBorder(true).SetTitle("Отправить файл").SetTitleAlign(tview.AlignLeft)
	app.SetRoot(form, true)
}
