package interfaces

import "context"

//go:generate mockgen -source=transaction.go -destination=mocks/transaction.go

// Transaction .
type Transaction interface {
	// transaction id
	ID() string
	// bind entity after entity save
	BindEntity(Entity)
	// pop after commit hooks from all bounded entities
	PopAfterCommitHooks() []func()
	// after commit
	AfterCommit(func())
	// pop rollback hooks from all bounded entities
	PopRollbackHooks() []func()
	// rollback
	OnRollback(func())
}

// TXContext .
type TXContext interface {
	context.Context
	Transaction
}

// TXFunc .
type TXFunc func(ctx context.Context, tx Transaction) error

// TXBeginner .
type TXBeginner interface {
	Transaction(txCtx TXContext, fn TXFunc) error
}

// TXManager .
type TXManager interface {
	Transaction(ctx context.Context, fn TXFunc) error
}
