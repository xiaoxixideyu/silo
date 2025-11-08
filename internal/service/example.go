package service

import (
	"context"
	"silo/internal/infra/filter"
	"silo/internal/infra/repo"
	"silo/internal/service/types"
	"silo/pkg/database/search"
	"silo/pkg/errs"
	"silo/pkg/platform"

	"github.com/samber/lo/mutable"
)

type ExampleService interface {
	Reverse(ctx context.Context, info string) string
	Page(ctx context.Context, req *types.ExamplePage) (*types.ExamplePageResult, error)
}

type ExampleServiceImpl struct {
	exampleRepo    repo.ExampleRepo
	platformClient *platform.PlatformClient
}

func NewExampleService(exampleRepo repo.ExampleRepo, platformClient *platform.PlatformClient) ExampleService {
	return &ExampleServiceImpl{
		exampleRepo:    exampleRepo,
		platformClient: platformClient,
	}
}

func (s *ExampleServiceImpl) Reverse(ctx context.Context, info string) string {
	data := []rune(info)
	mutable.Reverse(data)
	return string(data)
}

func (s *ExampleServiceImpl) Page(ctx context.Context, req *types.ExamplePage) (*types.ExamplePageResult, error) {
	pageFilter := &filter.ExamplePageFilter{
		Page: req.Page,
		Val:  req.Val,
		Sort: search.SortDESC,
	}
	pageResult, res, err := s.exampleRepo.FindPage(ctx, pageFilter)
	if err != nil {
		return nil, errs.NewCustomError(errs.StatusInternalServerError, err)
	}
	return &types.ExamplePageResult{
		PageResult: pageResult,
		Infos:      res,
	}, nil
}
