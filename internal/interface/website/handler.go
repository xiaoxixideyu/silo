package website

import (
	"silo/internal/interface/website/handlers"
	"silo/pkg/echo_handle"

	"github.com/labstack/echo/v4"
)

// @title QQFun Website API
// @version 1.0
// @description This is the API documentation for QQFun Website service
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8081
// @BasePath /api/v1

type WebsiteHandler interface {
	echo_handle.HandlerInterface
}

type WebsiteHandlerImpl struct {
	handlers []echo_handle.HandlerInterface
}

func NewWebsiteHandler(
	exampleHandler handlers.ExampleHandler,
) WebsiteHandler {
	return &WebsiteHandlerImpl{
		handlers: []echo_handle.HandlerInterface{
			exampleHandler,
		},
	}
}

func (h *WebsiteHandlerImpl) InitRouter(e *echo.Group) {
	g := e.Group("/api/v1")
	for _, handler := range h.handlers {
		handler.InitRouter(g)
	}
}
