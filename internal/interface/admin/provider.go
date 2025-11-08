package admin

import (
	"silo/internal/interface/admin/handlers"
	"silo/internal/service"

	"github.com/google/wire"
)

// HandlerSet 包含所有handler的wire set
var HandlerSet = wire.NewSet(
	NewAdminHandler,
	handlers.HandlerSet,
	service.ServiceSet,
)
