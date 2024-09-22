package ui

import (
	"context"
	"github.com/rivo/tview"
	"goph_keeper/internal/services/entities"
	v1 "goph_keeper/internal/services/grpc/goph_keeper/v1"
	"os"
	"strings"
)

func (m *Menu) showCardsMenu() {
	res, err := m.grpcClient.GetStoreDataList(context.Background(), v1.DataType_DATA_TYPE_CARD)
	if err != nil {
		m.errorHandler(err, func() {
			m.showCardsMenu()
		})
	}

	list := tview.NewList()
	list.AddItem("Add", "", 0, func() {
		m.showCardForm(&entities.FileCard{})
	})

	for _, entry := range res.Entries {
		parts := strings.Split(entry.UserPath, "/")
		firstPart := parts[0]
		uuid := entry.Uuid
		list.AddItem(firstPart, "", 0, func() {
			file, err := m.grpcClient.DownloadFile(context.Background(), uuid, os.TempDir(), nil)
			if err != nil {
				m.errorHandler(err, func() {
					m.showCardsMenu()
				})
				return
			}
			defer file.Close()
			defer os.Remove(file.Name())

			fileCard := &entities.FileCard{}
			err = fileCard.FromFile(file)
			if err != nil {
				m.errorHandler(err, func() {
					m.showCardsMenu()
				})
				return
			}
			fileCard.Uuid = uuid
			m.showCardForm(fileCard)
		})
	}
	list.AddItem("Back", "", 0, func() {
		m.showAppMenu()
	})

	m.app.SetRoot(list, true).SetFocus(list)
}

func (m *Menu) showCardForm(fileCard *entities.FileCard) {
	form := tview.NewForm().
		AddInputField("Card Name", fileCard.CardName, 20, nil, func(text string) {
			fileCard.CardName = text
		}).
		AddInputField("Description", fileCard.Description, 20, nil, func(text string) {
			fileCard.Description = text
		}).
		AddInputField("Card Number", fileCard.CardNumber, 20, nil, func(text string) {
			fileCard.CardNumber = text
		}).
		AddInputField("Expiry Date", fileCard.ExpiryDate, 10, nil, func(text string) {
			fileCard.ExpiryDate = text
		}).
		AddInputField("CVV", fileCard.CVV, 3, nil, func(text string) {
			fileCard.CVV = text
		}).
		AddInputField("Card Holder", fileCard.CardHolder, 20, nil, func(text string) {
			fileCard.CardHolder = text
		}).
		AddButton("Submit", func() {
			tmpFile, err := fileCard.ToFile()
			if err != nil {
				m.errorHandler(err, func() {
					m.showCardForm(fileCard)
				})
				return
			}
			defer os.Remove(tmpFile.Name())

			if fileCard.Uuid != "" {
				err := m.grpcClient.DeleteFile(context.Background(), fileCard.Uuid)
				if err != nil {
					m.errorHandler(err, func() {
						m.showCardForm(fileCard)
					})
					return
				}
			}

			_, err = m.grpcClient.UploadFile(context.Background(), tmpFile.Name(), fileCard.GetName(), v1.DataType_DATA_TYPE_CARD, nil)

			if err != nil {
				m.errorHandler(err, func() {
					m.showCardForm(fileCard)
				})
				return
			}
			m.showCardsMenu()
		}).
		AddButton("Cancel", func() {
			m.showCardsMenu()
		})

	form.SetBorder(true).SetTitle("Карты").SetTitleAlign(tview.AlignLeft)

	m.app.SetRoot(form, true)
}
