package mapper

import (
	"silo/internal/infra/entity"
	"silo/internal/interface/admin/dto"
	"silo/internal/service/types"
)

//go:generate go run github.com/jmattheis/goverter/cmd/goverter gen ./

// goverter:variables
var (
	ToTypesExamplePage func(req *dto.ExamplePageRequest) *types.ExamplePage
	//ToDtoExamplePageResponse func(source *types.ExamplePageResult) *dto.ExamplePageResponse
)

// goverter:ignore
func ToDtoExamplePageResponse(source *types.ExamplePageResult) *dto.ExamplePageResponse {
	if source == nil {
		return nil
	}

	// 手动转换，处理复杂的实体结构
	dtoResponse := &dto.ExamplePageResponse{
		PageResult: source.PageResult,
	}

	// 转换 Infos 数组，直接使用 entity.Example（goverter 无法处理私有字段）
	if source.Infos != nil {
		dtoResponse.Infos = make([]*entity.Example, len(source.Infos))
		for i, entityExample := range source.Infos {
			if entityExample != nil {
				// 直接复制指针，忽略内部私有字段
				dtoResponse.Infos[i] = entityExample
			}
		}
	}

	return dtoResponse
}
