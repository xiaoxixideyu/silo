package request

// CreateExampleRequest create example request
type CreateExampleRequest struct {
	Info string `json:"info"`
	Val  int64  `json:"val"`
}

// UpdateExampleRequest update example request
type UpdateExampleRequest struct {
	Info string `json:"info,omitempty"`
	Val  int64  `json:"val,omitempty"`
}

// GetExamplesRequest get example list request
type GetExamplesRequest struct {
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
	Search string `json:"search,omitempty"`
}

// ReverseRequest reverse request
type ReverseRequest struct {
	Info string `json:"info"`
}
