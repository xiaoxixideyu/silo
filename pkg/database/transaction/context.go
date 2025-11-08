package transaction

import (
	"context"
	"silo/pkg/database/interfaces"
)

type transactionKey struct {
}

type transactionContext struct {
	context.Context
	interfaces.Transaction
}

// NewContext .
func NewContext(ctx context.Context, tx interfaces.Transaction) interfaces.TXContext {
	return &transactionContext{
		Context:     context.WithValue(ctx, transactionKey{}, tx),
		Transaction: tx,
	}
}

// FromContext .
func FromContext(ctx context.Context) interfaces.Transaction {
	val := ctx.Value(transactionKey{})
	if val == nil {
		return nil
	}

	tx, ok := val.(interfaces.Transaction)
	if !ok {
		return nil
	}

	return tx
}
