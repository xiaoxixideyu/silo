package handlers

import (
	"fmt"
	"silo/internal/interface/admin/dto"
	"silo/internal/interface/admin/dto/mapper"
	"silo/internal/service"
	"silo/pkg/echo_handle"
	"silo/pkg/errs"

	"github.com/labstack/echo/v4"
)

type ExampleHandler interface {
	echo_handle.HandlerInterface
}

type ExampleHandlerImpl struct {
	echo_handle.Controller
	exampleService service.ExampleService
}

func NewExampleHandler(exampleService service.ExampleService) ExampleHandler {
	return &ExampleHandlerImpl{
		exampleService: exampleService,
	}
}

func (h *ExampleHandlerImpl) InitRouter(g *echo.Group) {
	exampleGroup := g.Group("/example")
	{
		exampleGroup.POST("/hello", h.Hello)
		exampleGroup.GET("/reverse", h.Reverse)
		exampleGroup.POST("/page", h.ExamplePage)
		exampleGroup.POST("/test", h.Test)
	}
}

// Hello godoc
// @Summary Example hello endpoint
// @Description Returns the input example string
// @Tags example
// @Accept json
// @Produce json
// @Param request body dto.ExampleRequest true "Example request"
// @Success 200 {object} dto.ExampleResponse
// @Router /example/hello [post]
func (h *ExampleHandlerImpl) Hello(ctx echo.Context) error {
	req := &dto.ExampleRequest{}
	if err := ctx.Bind(req); err != nil {
		return errs.NewBadRequestError().WithError(err)
	}
	if err := ctx.Validate(req); err != nil {
		return err
	}
	return h.JSONResponse(ctx, dto.ExampleResponse{Example: req.Example})
}

// Reverse godoc
// @Summary Reverse string endpoint
// @Description Returns the reversed input string
// @Tags example
// @Accept json
// @Produce json
// @Param request body dto.ExampleReverseRequest true "Reverse request"
// @Success 200 {object} dto.ExampleReverseResponse
// @Router /example/reverse [get]
func (h *ExampleHandlerImpl) Reverse(ctx echo.Context) error {
	req := &dto.ExampleReverseRequest{}
	if err := ctx.Bind(req); err != nil {
		return errs.NewBadRequestError().WithError(err)
	}
	if err := ctx.Validate(req); err != nil {
		return err
	}
	res := h.exampleService.Reverse(ctx.Request().Context(), req.Info)
	return h.JSONResponse(ctx, dto.ExampleReverseResponse{Info: res})
}

// ExamplePage godoc
// @Summary Get paginated examples
// @Description Returns a paginated list of examples
// @Tags example
// @Accept json
// @Produce json
// @Param request body dto.ExamplePageRequest true "Page request"
// @Success 200 {object} dto.ExamplePageResponse
// @Router /example/page [post]
func (h *ExampleHandlerImpl) ExamplePage(ctx echo.Context) error {
	req := &dto.ExamplePageRequest{}
	if err := ctx.Bind(req); err != nil {
		return errs.NewBadRequestError().WithError(err)
	}
	if err := ctx.Validate(req); err != nil {
		return err
	}
	res, err := h.exampleService.Page(ctx.Request().Context(), mapper.ToTypesExamplePage(req))
	if err != nil {
		return err
	}
	return h.JSONResponse(ctx, mapper.ToDtoExamplePageResponse(res))
}

// Test godoc
// @Summary Test endpoint
// @Description A simple test endpoint that returns a greeting message
// @Tags example
// @Accept json
// @Produce json
// @Param request body dto.TestRequest true "Test request"
// @Success 200 {object} dto.TestResponse
// @Router /example/test [post]
func (h *ExampleHandlerImpl) Test(ctx echo.Context) error {
	req := &dto.TestRequest{}
	if err := ctx.Bind(req); err != nil {
		return errs.NewBadRequestError().WithError(err)
	}
	if err := ctx.Validate(req); err != nil {
		return err
	}

	response := dto.TestResponse{
		Message: fmt.Sprintf("Hello, %s!", req.Name),
		Success: true,
	}

	return h.JSONResponse(ctx, response)
}
