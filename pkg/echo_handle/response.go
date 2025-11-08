package echo_handle

import (
	"silo/pkg/errs"
)

// JSONResponse .
type JSONResponse struct {
	Code    errs.StatusCode `json:"code"`
	Status  string          `json:"status"`
	Data    any             `json:"data,omitempty"`
	Message string          `json:"message,omitempty"`
}
