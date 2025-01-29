package models

import (
	"log/slog"
)

type BaseModel struct {
	*slog.Logger
}

func (b *BaseModel) ifErrorLog(err error) error {
	if err != nil {
		b.Logger.Error(err.Error())
	}
	return err
}
