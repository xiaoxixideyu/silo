package service

import (
	"silo/internal/infra/repo"
	"silo/pkg/platform"

	"github.com/google/wire"
)

// ServiceSet 包含所有 service 的wire set
var ServiceSet = wire.NewSet(
	repo.RepoSet,
	platform.NewPlatformClient,
	NewExampleService,
)
