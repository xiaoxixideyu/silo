package example

////go:generate go run accessor MyStruct,MyStruct2

import (
	"accessor/example/other"
	another "accessor/example/ref"
	"time"

	"silo/pkg/types/decimal"
)

type MyStruct struct {
	ModelBase    // ignore by default
	*other.SBase `accessor:",inline"`
	other.SBase2 `accessor:",inline"`
	*Current     `accessor:",inline"`
	Current2     `accessor:",inline"`
	Another      *another.Another `accessor:"customNameForAnother"`
	DecimalField decimal.Decimal  `accessor:";eq:Equal"` // use custom eq method: Equal
	StrEnum      StrEnum          `accessor:";eq:=="`    // use == for equal

	StrField         string
	BoolField        bool
	IntField         int
	TimeField        time.Time
	TimePointerField *time.Time
	DecimalMap       map[int]decimal.Decimal
	DecimalPMap      map[int]*decimal.Decimal
	MapAny           map[string]any
}

type StrEnum string

type MyStruct2 struct {
	ModelBase
	other.SBase `accessor:",inline"`
}

type Current struct {
	CurrentChild  `accessor:",inline"`
	CurrentField1 int
	CurrentField2 int
}

type CurrentChild struct {
	CurrentChildField1 int
	CurrentChildField2 int
}

type Current2 struct {
	Current2Field1 int
	Current2Field2 int
}

type ModelBase struct {
	FromBase map[string]any
	changes  map[string]any `accessor:"-"`
}

// GetChanges .
func (e *ModelBase) GetChanges() map[string]any {
	if e.changes == nil {
		return nil
	}
	result := e.changes
	e.changes = nil
	return result
}

// Update .
func (e *ModelBase) Update(name string, value any) {
	if e.changes == nil {
		e.changes = make(map[string]any)
	}
	e.changes[name] = value
}
