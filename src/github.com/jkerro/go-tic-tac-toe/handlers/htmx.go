package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func HtmxOnly(c echo.Context, handler func(echo.Context) error) error {
	headerKey := http.CanonicalHeaderKey("HX-Request")
	hxRequestHeader := c.Request().Header[headerKey]
	errorFunc := func() error {
		return c.String(http.StatusBadRequest, "This endpoint is HTMX only")
	}
	if hxRequestHeader == nil {
		return errorFunc()
	}
	if hxRequest, err := strconv.ParseBool(hxRequestHeader[0]); err != nil || !hxRequest {
		return errorFunc()
	} else {
		return handler(c)
	}
}
