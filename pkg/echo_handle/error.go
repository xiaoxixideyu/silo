package echo_handle

import (
	"encoding/json"
	"net/http"
	"silo/pkg/errs"

	"github.com/labstack/echo/v4"
)

type errorHandlerInfo struct {
	debug          bool
	withOutMessage bool
}

var info = errorHandlerInfo{}

func SetErrorHandlerInfo(debug bool, withOutMessage bool) {
	info.debug = debug
	info.withOutMessage = withOutMessage
}

func ErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}
	// if debug is false and WithOutMessage is true, make response without message
	withOutMessage := !info.debug && info.withOutMessage

	var he *echo.HTTPError

	switch e := err.(type) {
	case *errs.CustomError:
		he = e.ToHTTPError(withOutMessage)
	case *echo.HTTPError:
		he = e
		if he.Internal != nil {
			if herr, ok := he.Internal.(*echo.HTTPError); ok {
				he = herr
			}
		}
	default:
		he = echo.NewHTTPError(http.StatusInternalServerError, e.Error())
	}

	// Issue #1426
	code := he.Code
	message := he.Message

	switch m := he.Message.(type) {
	case string:
		message = map[string]any{"code": code, "message": m}
	case json.Marshaler:
		// do nothing - this type knows how to format itself to JSON
	case error:
		message = map[string]any{"code": code, "message": m.Error()}
	}

	// Send response
	if c.Request().Method == http.MethodHead { // Issue #608
		err = c.NoContent(he.Code)
	} else {
		err = c.JSON(code, message)
	}
	if err != nil {
		c.Echo().Logger.Error(err)
	}
}
