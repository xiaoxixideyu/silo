package tracer

import (
	"context"

	"go.elastic.co/apm/v2"
)

// Tracer .
type Tracer = apm.Tracer

// Transaction .
type Transaction = apm.Transaction

// WithContext .
func WithContext(ctx context.Context, trans *Transaction) context.Context {
	return apm.ContextWithTransaction(ctx, trans)
}

// FromContext .
func FromContext(ctx context.Context) *Transaction {
	return apm.TransactionFromContext(ctx)
}

// SpanFromContext .
func SpanFromContext(ctx context.Context, name string, spanType string) *apm.Span {
	tx := FromContext(ctx)
	if tx == nil {
		return nil
	}
	parentSpan := apm.SpanFromContext(ctx)

	return tx.StartSpan(name, spanType, parentSpan)
}

// DefaultTracer .
func DefaultTracer() *Tracer {
	tracer := apm.DefaultTracer()

	return tracer
}
