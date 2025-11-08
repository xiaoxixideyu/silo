package echo_handle

import "github.com/labstack/echo/v4"

type HandlerInterface interface {
	InitRouter(group *echo.Group)
}
