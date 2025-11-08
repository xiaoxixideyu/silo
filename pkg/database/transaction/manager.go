package transaction

import (
	"context"
	"silo/pkg/database/interfaces"
)

// NewManager .
func NewManager(beginner interfaces.TXBeginner) interfaces.TXManager {
	txFactory := newFactory()
	return &Manager{
		txFactory: txFactory,
		beginner:  beginner,
		// attemptCount: 1,
		// retryIfFunc:  defaultRetryIfFunc,
	}
}

func defaultRetryIfFunc(_ error) bool {
	return true
}

// Manager .
type Manager struct {
	txFactory Factory
	beginner  interfaces.TXBeginner
	// attemptCount int
	// retryIfFunc  func(err error) bool
}

// Transaction .
func (m *Manager) Transaction(ctx context.Context, fn interfaces.TXFunc) error {
	tx := FromContext(ctx)
	if tx != nil {
		return fn(ctx, tx)
	}

	// tx with retry
	// return retry.Do(func() error {
	tx = m.txFactory.New(ctx)
	txCtx := NewContext(ctx, tx)
	afterCommitHooks := []func(){}
	rollbackHooks := []func(){}

	// do transaction
	err := m.beginner.Transaction(txCtx, func(ctx context.Context, tx interfaces.Transaction) error {
		// 返回error前，需要获取rollback hooks
		defer func() {
			rollbackHooks = tx.PopRollbackHooks()
		}()

		if err := fn(ctx, tx); err != nil {
			return err
		}

		afterCommitHooks = tx.PopAfterCommitHooks()

		return nil
	})

	if err != nil {
		m.executeRollbackHooks(rollbackHooks)
		return err
	}

	// execute after commit hooks
	for _, hook := range afterCommitHooks {
		hook()
	}

	return nil
	// }, retry.Attempts(uint(m.attemptCount)), retry.RetryIf(m.retryIfFunc), retry.LastErrorOnly(true))
}

func (m *Manager) executeRollbackHooks(hooks []func()) {
	for _, hook := range hooks {
		hook()
	}
}
