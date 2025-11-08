package response

import "time"

// Example example information
type Example struct {
	ID        int64     `json:"id"`
	Info      string    `json:"info"`
	Val       int64     `json:"val"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// GetExamplesResponse get example list response
type GetExamplesResponse struct {
	Examples []Example `json:"examples"`
	Total    int       `json:"total"`
	Page     int       `json:"page"`
	Limit    int       `json:"limit"`
}

// CreateExampleResponse create example response
type CreateExampleResponse struct {
	ID int64 `json:"id"`
}

// UpdateExampleResponse update example response
type UpdateExampleResponse struct {
	Success bool `json:"success"`
}

// DeleteExampleResponse delete example response
type DeleteExampleResponse struct {
	Success bool `json:"success"`
}

// ReverseResponse reverse response
type ReverseResponse struct {
	Info string `json:"info"`
}
