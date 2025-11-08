package dto

import (
	"silo/internal/infra/entity"
	"silo/pkg/database/search"
)

type ExampleRequest struct {
	Example string `json:"example" query:"example" example:"hello world" validate:"required"`
}

type ExampleResponse struct {
	Example string `json:"example" query:"example" example:"hello world"`
}

type ExampleReverseRequest struct {
	Info string `json:"info" query:"info" example:"hello world" validate:"required"`
}

type ExampleReverseResponse struct {
	Info string `json:"info" query:"info" example:"dlrow olleh"`
}

type ExamplePageRequest struct {
	search.Page
	Val int64 `json:"val" example:"100"`
}

type ExamplePageResponse struct {
	*search.PageResult
	Infos []*entity.Example `json:"infos"`
}
