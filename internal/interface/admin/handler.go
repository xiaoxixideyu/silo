package admin

import (
	"silo/internal/interface/admin/handlers"

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
	InitRouter(e *echo.Echo)
}

type adminHandlerImpl struct {
	exampleHandler handlers.ExampleHandler
}

func NewAdminHandler(
	exampleHandler handlers.ExampleHandler,
) AdminHandler {
	return &adminHandlerImpl{
		exampleHandler: exampleHandler,
	}
}

func (h *adminHandlerImpl) InitRouter(e *echo.Echo) {
	group := e.Group("/api/v1")
	h.exampleHandler.InitRouter(group)
}
