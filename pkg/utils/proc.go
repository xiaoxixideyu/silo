package utils

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"reflect"
	"runtime/debug"
	"strconv"

	"go.elastic.co/apm/v2"
)

// PCall recover panic
func PCall(fn func() error) (err error) {
	defer func() {
		if rec := recover(); rec != nil {
			// Try to use logger from context here to help trace error cause
			stackTrace := debug.Stack()
			stackTraceAsRawStringLiteral := strconv.Quote(string(stackTrace))
			slog.Error("panic - PCall", "panicData", rec, "stackTrace", stackTraceAsRawStringLiteral)

			err = fmt.Errorf("internal error - %v", rec)
		}
	}()
	return fn()
}

// PContextCall recover panic
func PContextCall(ctx context.Context, fn func(ctx context.Context) error) (err error) {
	defer func() {
		if rec := recover(); rec != nil {
			// Try to use logger from context here to help trace error cause
			stackTrace := debug.Stack()
			stackTraceAsRawStringLiteral := strconv.Quote(string(stackTrace))
			slog.Error("panic - PCall", "panicData", rec, "stackTrace", stackTraceAsRawStringLiteral)
			apm.CaptureError(ctx, fmt.Errorf("panic %v", rec)).Send()

			err = fmt.Errorf("internal error - %v", rec)
		}
	}()
	return fn(ctx)
}

// PArgsCall calls a method that returns an interface and an error and recovers in case of panic
func PArgsCall(method reflect.Method, args []reflect.Value) (rets any, err error) {
	defer func() {
		if rec := recover(); rec != nil {
			// Try to use logger from context here to help trace error cause
			stackTrace := debug.Stack()
			stackTraceAsRawStringLiteral := strconv.Quote(string(stackTrace))
			slog.Error(fmt.Sprintf("panic - pitaya/dispatch: methodName=%s panicData=%v stackTrace=%s", method.Name, rec, stackTraceAsRawStringLiteral))

			if s, ok := rec.(string); ok {
				err = errors.New(s)
			} else {
				err = fmt.Errorf("rpc call internal error - %s: %v", method.Name, rec)
			}
		}
	}()

	r := method.Func.Call(args)
	// r can have 0 length in case of notify handlers
	// otherwise it will have 2 outputs: an interface and an error
	if len(r) == 2 {
		if v := r[1].Interface(); v != nil {
			err = v.(error)
		} else if !r[0].IsNil() {
			rets = r[0].Interface()
		} else {
			err = fmt.Errorf("rpc call internal error - ErrReplyShouldBeNotNull")
		}
	}
	return
}
