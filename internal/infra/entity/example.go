package entity

import (
	"silo/internal/infra/model"
	"silo/pkg/database/seedwork"
)

type Example struct {
	seedwork.Entity
	model.Example
}

func NewExample(
	ID int64,
	Info string,
	Val int64,
) *Example {
	e := &Example{
		Example: model.Example{
			ID:   ID,
			Info: Info,
			Val:  Val,
		},
	}
	e.MarkNew()
	return e
}
