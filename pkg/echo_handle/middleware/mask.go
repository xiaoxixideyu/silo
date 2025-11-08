package middleware

import (
	"github.com/labstack/echo/v4"
)

// CtxKeyMasker request/response mask field
const CtxKeyMasker = "masker"

// NewMaskerMiddleware for request or response mask
func NewMaskerMiddleware(m *Masker) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set(CtxKeyMasker, m)
			return next(c)
		}
	}
}

// NewMasker .
func NewMasker(reqMasker MaskerFunc, respMasker MaskerFunc) *Masker {
	return &Masker{
		RequestMasker:  reqMasker,
		ResponseMasker: respMasker,
	}
}

// MaskerFunc .
type MaskerFunc = func(m any) string

// Masker .
type Masker struct {
	RequestMasker  MaskerFunc
	ResponseMasker MaskerFunc
}

// MaskerRequest .
func (m *Masker) MaskerRequest(req any) any {
	if m.RequestMasker != nil {
		return m.RequestMasker(req)
	}

	return req
}

// MaskerResponse .
func (m *Masker) MaskerResponse(resp any) any {
	if m.ResponseMasker != nil {
		return m.ResponseMasker(resp)
	}

	return resp
}
