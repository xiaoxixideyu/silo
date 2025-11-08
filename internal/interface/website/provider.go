package website

import (
	"silo/internal/interface/website/handlers"
	"silo/internal/service"

	"github.com/google/wire"
)

// HandlerSet 包含所有handler的wire set
var HandlerSet = wire.NewSet(
	NewWebsiteHandler,
	handlers.NewExampleHandler,
	service.ServiceSet,
)
