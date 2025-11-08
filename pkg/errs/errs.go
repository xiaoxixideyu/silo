package errs

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	pkgErrors "github.com/pkg/errors"
)

// Max stack trace depth
const MAX_STACK_TRACE_DEPTH = 3
const MAX_CALLER_SKIP = 3

// ErrorCode should implement this interface
type ErrorCode interface {
	Code() int
	String() string
	Recoverable() bool
}

// CustomError ce
type CustomError struct {
	Code       ErrorCode         `json:"code"`
	Status     string            `json:"status,omitempty"`
	Message    string            `json:"message,omitempty"`
	Data       map[string]string `json:"data,omitempty"`
	Additional map[string]any    `json:"additional,omitempty"` // additional data to be included in the response
	HTTPCode   int               `json:"-"`
	Internal   any               `json:"-"` // Stores the error returned by an external dependency
	stack      []uintptr         `json:"-"` // Stack trace for debugging
}

// Error makes it compatible with `error` interface.
func (e *CustomError) Error() string {
	msg := strings.Builder{}
	msg.WriteString("code=")
	msg.WriteString(strconv.Itoa(e.Code.Code()))
	if e.Message != "" {
		msg.WriteString(", message=")
		msg.WriteString(e.Message)
	}
	if e.Internal != nil {
		msg.WriteString(", internal=")
		msg.WriteString(fmt.Sprintf("%v", e.Internal))
	}

	return msg.String()
}

// ToHTTPError ToHTTPError
func (e *CustomError) ToHTTPError(withOutMessage bool) *echo.HTTPError {
	var err *echo.HTTPError
	if e.HTTPCode == 0 {
		err = echo.NewHTTPError(http.StatusOK)
	} else {
		err = echo.NewHTTPError(e.HTTPCode)
	}

	msg := map[string]any{}
	msg["code"] = e.Code.Code()

	if !withOutMessage {
		msg["status"] = e.Status
		if e.Message != "" {
			msg["message"] = e.Message
		}
	}

	if e.Data != nil {
		msg["data"] = e.Data
	}

	if e.Additional != nil {
		for k, v := range e.Additional {
			msg[k] = v
		}
	}

	err.Message = msg

	return err
}

// WithHTTPCode .
func (e *CustomError) WithHTTPCode(code int) *CustomError {
	e.HTTPCode = code
	return e
}

// WithMessage WithMessage
func (e *CustomError) WithMessage(message ...string) *CustomError {
	e.Message = strings.Join(message, "")
	return e
}

// WithMessageF WithMessageF
func (e *CustomError) WithMessageF(format string, args ...any) *CustomError {
	e.Message = fmt.Sprintf(format, args...)
	return e
}

// WithError WithError
func (e *CustomError) WithError(err error) *CustomError {
	e.Internal = err
	return e
}

// WithoutInternalError .
func (e *CustomError) WithoutInternalError() *CustomError {
	e.Internal = nil
	return e
}

// WithData WithData
func (e *CustomError) WithData(data map[string]string) *CustomError {
	if e.Data == nil {
		e.Data = data
	} else {
		for key, value := range data {
			e.Data[key] = value
		}
	}
	return e
}

// WithDataKV WithData
func (e *CustomError) WithDataKV(key string, value string) *CustomError {
	if e.Data == nil {
		e.Data = make(map[string]string)
	}

	e.Data[key] = value
	return e
}

// ClearData clears all data
func (e *CustomError) ClearData() *CustomError {
	e.Data = nil

	return e
}

// IsServerError .
func (e *CustomError) IsServerError() bool {
	return e.Code == StatusInternalServerError
}

// WithAdditional Additional
func (e *CustomError) WithAdditional(key string, value any) *CustomError {
	if e.Additional == nil {
		e.Additional = make(map[string]any)
	}

	e.Additional[key] = value

	return e
}

// WithStack WithData
func (e *CustomError) WithStack() *CustomError {
	return e.WithError(pkgErrors.WithStack(e))
}

func (e *CustomError) GetData(key string) string {
	if e.Data == nil {
		return ""
	}

	return e.Data[key]
}

// StackTrace returns the stack trace of the error.
func (e *CustomError) StackTrace() string {
	if e.stack == nil {
		return e.Error()
	}

	frames := runtime.CallersFrames(e.stack)
	trace := strings.Builder{}
	trace.WriteString(e.Error())
	trace.WriteString(", ")
	return framesToString(&trace, frames)
}

func GetStackTrace(skip int, depth int) string {
	stack := make([]uintptr, depth)
	length := runtime.Callers(skip, stack[:])
	stack = stack[:length]
	frames := runtime.CallersFrames(stack)
	trace := strings.Builder{}
	return framesToString(&trace, frames)
}

func framesToString(trace *strings.Builder, frames *runtime.Frames) string {
	trace.WriteString("stack=")
	for {
		frame, more := frames.Next()
		if more {
			trace.WriteString(fmt.Sprintf("%s:%d; ", frame.Function, frame.Line))
		} else {
			trace.WriteString(fmt.Sprintf("%s:%d", frame.Function, frame.Line))
			break
		}
	}
	return trace.String()
}

// WithStack .
func WithStack(err error, msg ...any) *CustomError {
	if err == nil {
		return nil
	}

	return NewInternalServerError(msg...).WithError(pkgErrors.WithStack(err))
}

// WithStackf .
func WithStackf(err error, format string, args ...any) *CustomError {
	if err == nil {
		return nil
	}

	return WithStack(err, fmt.Sprintf(format, args...))
}

// Is error
func Is(err error, target error) bool {
	a, aok := err.(*CustomError)
	b, bok := target.(*CustomError)
	if aok && bok {
		return a.Code == b.Code
	}

	return errors.Is(err, target)
}

func GetCode(err error) ErrorCode {
	if err == nil {
		//无错误信息
		return StatusOK
	}
	var e *CustomError
	ok := errors.As(err, &e)
	if !ok {
		//非自定义错误格式消息
		return StatusUnknown
	}
	return e.Code
}

// IsCode .
func IsCode(err error, code ErrorCode) bool {
	a, ok := err.(*CustomError)
	if ok {
		return a.Code.Code() == code.Code()
	}

	return false
}

// IsInstanceNotFound error
func IsInstanceNotFound(err error) bool {
	return IsCode(err, StatusInstanceNotFound)
}

// IsDuplicated error
func IsDuplicated(err error) bool {
	return IsCode(err, StatusDuplicated)
}

// IsInternalError .
func IsInternalError(err error) bool {
	_, ok := err.(*CustomError)
	if !ok {
		return true
	}

	return IsCode(err, StatusInternalServerError)
}

func IsRecoverable(err error) bool {
	e, ok := err.(*CustomError)
	if ok {
		// recoverable custom error code
		return e.Code.Recoverable()
	}

	return true
}

func IsRequestCanceled(err error) bool {
	if IsCode(err, StatusRequestCanceled) {
		return true
	}

	return errors.Is(err, context.Canceled)
}

func IsRequestTimeout(err error) bool {
	if IsCode(err, StatusRequestTimeout) {
		return true
	}

	return errors.Is(err, context.DeadlineExceeded)
}

func IsCustomError(err error) bool {
	_, ok := err.(*CustomError)
	return ok
}

func ToCustomError(err error) *CustomError {
	if err == nil {
		return nil
	}

	ce, ok := err.(*CustomError)
	if !ok {
		return nil
	}

	return ce
}

// NewCustomError NewCustomError
func NewCustomError(code ErrorCode, msg ...any) *CustomError {
	return newCustomErrorWithStack(getStack(), code, msg...)
}

// newCustomErrorWithStack NewCustomError
func newCustomErrorWithStack(stack []uintptr, code ErrorCode, msg ...any) *CustomError {
	e := &CustomError{
		Code:   code,
		Status: code.String(),
		stack:  stack,
	}

	e.Message = fmt.Sprint(msg...)

	return e
}

func getStack() []uintptr {
	stack := make([]uintptr, MAX_STACK_TRACE_DEPTH)
	length := runtime.Callers(MAX_CALLER_SKIP, stack[:])
	return stack[:length]
}

// NewCustomErrorWithoutStack NewCustomError
func newCustomErrorWithoutStack(code ErrorCode, msg ...any) *CustomError {
	e := &CustomError{
		Code:   code,
		Status: code.String(),
	}

	e.Message = fmt.Sprint(msg...)

	return e
}

// NewCustomErrorf NewCustomError
func NewCustomErrorf(code ErrorCode, format string, msg ...any) *CustomError {
	return newCustomErrorWithStack(getStack(), code, fmt.Sprintf(format, msg...))
}

// NewCustomErrorf NewCustomError
func newCustomErrorWithStackf(stack []uintptr, code ErrorCode, format string, msg ...any) *CustomError {
	return newCustomErrorWithStack(stack, code, fmt.Sprintf(format, msg...))
}

// NewBadRequestError .
func NewBadRequestError(msg ...any) *CustomError {
	return newCustomErrorWithStack(getStack(), StatusBadRequest, msg...)
}

func NewBadRequestErrorf(format string, msg ...any) *CustomError {
	return newCustomErrorWithStackf(getStack(), StatusBadRequest, format, msg...)
}

// NewUnauthorizedError .
func NewUnauthorizedError(msg ...any) *CustomError {
	return newCustomErrorWithStack(getStack(), StatusUnauthorized, msg...)
}

// NewNotFoundError .
func NewNotFoundError(msg ...any) *CustomError {
	return newCustomErrorWithStack(getStack(), StatusNotFound, msg...)
}

// NewInternalServerError .
func NewInternalServerError(msg ...any) *CustomError {
	return newCustomErrorWithStack(getStack(), StatusInternalServerError, msg...).WithHTTPCode(http.StatusInternalServerError)
}

// NewInternalServerErrorf .
func NewInternalServerErrorf(format string, msg ...any) *CustomError {
	return newCustomErrorWithStackf(getStack(), StatusInternalServerError, format, msg...).WithHTTPCode(http.StatusInternalServerError)
}

// NewTransactionAbortError .
func NewTransactionAbortError(msg ...any) *CustomError {
	return newCustomErrorWithStack(getStack(), StatusTransactionAbort, msg...)
}

// NewTransactionAbortErrorf .
func NewTransactionAbortErrorf(format string, msg ...any) *CustomError {
	return newCustomErrorWithStackf(getStack(), StatusTransactionAbort, format, msg...).WithHTTPCode(http.StatusInternalServerError)
}

// NewRequestCanceledError .
func NewRequestCanceledError(msg ...any) *CustomError {
	return newCustomErrorWithStack(getStack(), StatusRequestCanceled, msg...)
}

// NewRequestCanceledErrorf .
func NewRequestCanceledErrorf(format string, msg ...any) *CustomError {
	return newCustomErrorWithStackf(getStack(), StatusRequestCanceled, format, msg...)
}

// NewForbiddenError .
func NewForbiddenError(msg ...any) *CustomError {
	return newCustomErrorWithStack(getStack(), StatusForbidden, msg...)
}

// NewServerMaintainError .
func NewServerMaintainError(msg ...any) *CustomError {
	return newCustomErrorWithStack(getStack(), StatusServerMaintain, msg...)
}

// NewInstanceNotFoundError .
func NewInstanceNotFoundError(msg ...any) *CustomError {
	return newCustomErrorWithStack(getStack(), StatusInstanceNotFound, msg...)
}

// NewNotImplementError .
func NewNotImplementError(msg ...any) *CustomError {
	return newCustomErrorWithStack(getStack(), StatusBadRequest, msg...)
}

// NewNotImplementErrorf .
func NewNotImplementErrorf(format string, msg ...any) *CustomError {
	return newCustomErrorWithStackf(getStack(), StatusBadRequest, format, msg...)
}

// NewNotEnoughResourceError .
func NewNotEnoughResourceError(msg ...any) *CustomError {
	return newCustomErrorWithStack(getStack(), StatusNotEnoughResource, msg...)
}

// NewDuplicatedError .
func NewDuplicatedError(msg ...any) *CustomError {
	return newCustomErrorWithStack(getStack(), StatusDuplicated, msg...)
}

func NewRequestTimeoutError(msg ...any) *CustomError {
	return newCustomErrorWithoutStack(StatusRequestTimeout, msg...)
}

func NewRequestTimeoutErrorf(format string, msg ...any) *CustomError {
	return newCustomErrorWithStackf(getStack(), StatusRequestTimeout, format, msg...)
}

func NewJSONUnmarshalError(msg ...any) *CustomError {
	return newCustomErrorWithStack(getStack(), StatusJSONUnmarshalError, msg...)
}

func NewJSONUnmarshalErrorf(format string, msg ...any) *CustomError {
	return newCustomErrorWithStackf(getStack(), StatusJSONUnmarshalError, format, msg...)
}
