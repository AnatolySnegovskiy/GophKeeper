package ui

import (
	"path/filepath"

	"github.com/rivo/tview"
)

func (m *Menu) showFilesMenu() {
	title := tview.NewTextView().
		SetText("Files Menu").
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true)

	list := tview.NewList().
		AddItem("Download", "Download a file", 'd', func() {
			m.showServerFilesMenu(string(filepath.Separator))
		}).
		AddItem("Upload", "Upload a file", 'u', func() {
			m.showSendFileForm()
		}).AddItem("Back", "", 'e', func() {
		m.showAppMenu()
	})

	mainLayout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(title, 3, 1, false).
		AddItem(list, 0, 1, true)

	m.app.SetRoot(mainLayout, true).SetFocus(list)
}
