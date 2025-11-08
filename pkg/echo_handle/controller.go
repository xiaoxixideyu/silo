package echo_handle

import (
	"net/http"
	"silo/pkg/errs"
	"strconv"

	"github.com/labstack/echo/v4"
)

// Controller .
type Controller struct{}

// JSONResponse JSONResponse
func (a *Controller) JSONResponse(c echo.Context, data any) error {
	resp := JSONResponse{
		Code:   errs.StatusOK,
		Status: errs.StatusOK.String(),
		Data:   data,
	}
	return c.JSON(http.StatusOK, &resp)
}

// Redirect Redirect
func (a *Controller) Redirect(c echo.Context, url string) error {
	return c.Redirect(http.StatusMovedPermanently, url)
}

// RedirectTemporary .
func (a *Controller) RedirectTemporary(c echo.Context, url string) error {
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

// SimpleJSONResponse .
func (a *Controller) SimpleJSONResponse(c echo.Context, code int, status string, message string) error {
	resp := JSONResponse{
		Status:  status,
		Message: message,
	}

	return c.JSON(code, resp)
}

func (a *Controller) GetID(ctx echo.Context) (int64, error) {
	var param = "id"
	idStr := ctx.Param(param)
	if idStr == "" {
		return 0, errs.NewBadRequestError().WithMessageF("The %s is empty", param)
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, errs.NewBadRequestError().WithMessageF("The %s is empty", param)
	}
	return id, nil
}
