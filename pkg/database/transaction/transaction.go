package transaction

import (
	"context"
	"silo/pkg/database/interfaces"

	"github.com/google/uuid"
)

//go:generate mockgen -source=transaction.go -destination=mocks/transaction.go

// Factory .
type Factory interface {
	New(context.Context) interfaces.Transaction
}

// NewFactory .
func newFactory() Factory {
	return &factoryImpl{}
}

// FactoryImpl .
type factoryImpl struct{}

// New .
func (f *factoryImpl) New(ctx context.Context) interfaces.Transaction {
	return &transaction{
		id:       uuid.New().String(),
		ctx:      ctx,
		entities: []interfaces.Entity{},
	}
}

func newTransaction(ctx context.Context) interfaces.Transaction {
	return &transaction{
		id:       uuid.New().String(),
		ctx:      ctx,
		entities: []interfaces.Entity{},
	}
}

type transaction struct {
	id               string
	ctx              context.Context
	entities         []interfaces.Entity
	afterCommitHooks []func()
	rollbackHooks    []func()
}

// transaction id
func (tx *transaction) ID() string {
	return tx.id
}

// bind entity after entity save
func (tx *transaction) BindEntity(entity interfaces.Entity) {
	// ignore duplicated binding
	for _, v := range tx.entities {
		if v == entity {
			return
		}
	}
	tx.entities = append(tx.entities, entity)
}

// pop after commit hooks from all bounded entities
func (tx *transaction) PopAfterCommitHooks() []func() {
	hooks := []func(){}
	hooks = append(hooks, tx.afterCommitHooks...)
	tx.afterCommitHooks = []func(){}

	for _, v := range tx.entities {
		pending := v.PopAfterCommitHooks()
		hooks = append(hooks, pending...)
	}
	return hooks
}

// AfterCommit .
func (tx *transaction) AfterCommit(hook func()) {
	tx.afterCommitHooks = append(tx.afterCommitHooks, hook)
}

// pop rollback hooks from all bounded entities
func (tx *transaction) PopRollbackHooks() []func() {
	hooks := []func(){}
	hooks = append(hooks, tx.rollbackHooks...)
	tx.rollbackHooks = []func(){}

	return hooks
}

// rollback
func (tx *transaction) OnRollback(hook func()) {
	tx.rollbackHooks = append(tx.rollbackHooks, hook)
}
