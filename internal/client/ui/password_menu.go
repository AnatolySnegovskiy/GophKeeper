package ui

import (
	"context"
	"goph_keeper/internal/services/entities"
	v1 "goph_keeper/internal/services/grpc/goph_keeper/v1"
	"os"
	"strings"

	"github.com/rivo/tview"
)

func (m *Menu) showPasswordMenu() {
	res, err := m.grpcClient.GetStoreDataList(context.Background(), v1.DataType_DATA_TYPE_LOGIN_PASSWORD)
	if err != nil {
		m.errorHandler(err, func() {
			m.showAppMenu()
		})
		return
	}

	list := tview.NewList()
	list.AddItem("Add", "", 0, func() {
		m.showPasswordForm(&entities.FilePassword{})
	})

	for _, entry := range res.Entries {
		parts := strings.Split(entry.UserPath, "/")
		firstPart := parts[0]
		uuid := entry.Uuid
		list.AddItem(firstPart, "", 0, func() {
			file, err := m.grpcClient.DownloadFile(context.Background(), uuid, os.TempDir(), nil)
			if err != nil {
				m.errorHandler(err, func() {
					m.showPasswordMenu()
				})
				return
			}
			defer file.Close()
			defer os.Remove(file.Name())
			filePassword := &entities.FilePassword{}
			err = filePassword.FromFile(file)

			if err != nil {
				m.errorHandler(err, func() {
					m.showPasswordMenu()
				})
				return
			}

			m.showPasswordForm(filePassword)
		})
	}
	list.AddItem("Back", "", 0, func() {
		m.showAppMenu()
	})

	m.app.SetRoot(list, true).SetFocus(list)
}

func (m *Menu) showPasswordForm(passwordFile *entities.FilePassword) {
	form := tview.NewForm().
		AddInputField("Title", passwordFile.Title, 20, nil, func(text string) {
			passwordFile.Title = text
		}).
		AddTextArea("Description", passwordFile.Description, 20, 4, 255, func(text string) {
			passwordFile.Description = text
		}).
		AddInputField("Login", passwordFile.Login, 20, nil, func(text string) {
			passwordFile.Login = text
		}).
		AddInputField("Password", passwordFile.Password, 20, nil, func(text string) {
			passwordFile.Password = text
		}).
		AddButton("Submit", func() {
			tmpFile, err := passwordFile.ToFile()
			if err != nil {
				m.errorHandler(err, func() {
					m.showPasswordForm(passwordFile)
				})
				return
			}
			defer os.Remove(tmpFile.Name())
			if passwordFile.Uuid != "" {
				err := m.grpcClient.DeleteFile(context.Background(), passwordFile.Uuid)
				if err != nil {
					m.errorHandler(err, func() {
						m.showPasswordForm(passwordFile)
					})
					return
				}
			}

			_, err = m.grpcClient.UploadFile(context.Background(), tmpFile.Name(), passwordFile.GetName(), v1.DataType_DATA_TYPE_LOGIN_PASSWORD, nil)

			if err != nil {
				m.errorHandler(err, func() {
					m.showPasswordForm(passwordFile)
				})
				return
			}
			m.showPasswordMenu()
		}).
		AddButton("Cancel", func() {
			m.showPasswordMenu()
		})

	form.SetBorder(true).SetTitle("Пароли").SetTitleAlign(tview.AlignLeft)

	m.app.SetRoot(form, true)
}
