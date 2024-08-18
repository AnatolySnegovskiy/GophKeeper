package ui

import (
	"context"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (m *Menu) showAppMenu() {
	title := tview.NewTextView().
		SetText(m.title).
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true)

	list := tview.NewList().
		AddItem("1. Отправить файл", "", '1', func() {
			m.showSendFileForm()
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

	m.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc:
			m.app.Stop()
			return nil
		}
		return event
	})

	m.app.SetRoot(mainLayout, true).SetFocus(list)
}

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

		form := tview.NewForm().
			AddFormItem(info).
			AddFormItem(progressBar)
		form.SetBorder(true).SetTitle("Send File").SetTitleAlign(tview.AlignLeft)
		flex := tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(form, 0, 1, true)

		m.app.SetRoot(flex, true)

		progressChan := make(chan int)

		go func() {
			res, err := m.grpcClient.StoreData(context.Background(), filePath, progressChan)
			if err != nil {
				m.app.QueueUpdateDraw(func() {
					info.SetText(fmt.Sprintf("[red]Error: %s", err))
				})
			}
			close(progressChan)
			m.app.QueueUpdateDraw(func() {
				info.SetText(fmt.Sprintf("[green]Success: %s", res.Message))
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
	}

	m.explore(f)
}