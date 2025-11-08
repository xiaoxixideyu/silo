package repo

import (
	"context"
	"silo/internal/infra/entity"
	"silo/internal/infra/filter"
	"silo/internal/infra/model"
	"silo/pkg/database/interfaces"
	"silo/pkg/database/search"
)

type ExampleRepo interface {
	FindPage(ctx context.Context, filter *filter.ExamplePageFilter) (*search.PageResult, []*entity.Example, error)
}

type exampleRepoImpl struct {
	dao interfaces.DaoCache
}

func NewExampleRepo(dao interfaces.DaoCache) ExampleRepo {
	return &exampleRepoImpl{
		dao: dao,
	}
}

func (r *exampleRepoImpl) FindPage(ctx context.Context, filter *filter.ExamplePageFilter) (*search.PageResult, []*entity.Example, error) {
	res := make([]*entity.Example, 0)
	p, err := r.dao.Page(ctx, model.TableNameExample, filter, &res)
	return p, res, err
}
