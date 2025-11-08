package admin

import (
	"silo/internal/interface/admin/handlers"
	"silo/pkg/echo_handle"

	"github.com/labstack/echo/v4"
)

// @title QQFun Admin API
// @version 1.0
// @description This is the API documentation for QQFun Admin service
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1

type AdminHandler interface {
	echo_handle.HandlerInterface
}

type adminHandlerImpl struct {
	handlers []echo_handle.HandlerInterface
}

func NewAdminHandler(
	exampleHandler handlers.ExampleHandler,
) AdminHandler {
	return &adminHandlerImpl{
		handlers: []echo_handle.HandlerInterface{
			exampleHandler,
		},
	}
}

func (h *adminHandlerImpl) InitRouter(e *echo.Group) {
	group := e.Group("/api/v1")
	for _, handler := range h.handlers {
		handler.InitRouter(group)
	}
}
