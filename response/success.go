package response

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/Dert12318/Utilities/models/v2"
)

func (r *HttpResponse) SuccessResponse(ec echo.Context, message string, data interface{}) error {
	if message == "" {
		message = http.StatusText(http.StatusOK)
	}
	return ec.JSON(http.StatusOK, models.BuildSuccessResponse(message, http.StatusOK, data))
}

func (r *HttpResponse) SuccessResponseWithCode(ec echo.Context, code int, message string, data interface{}) error {
	if message == "" {
		message = http.StatusText(http.StatusOK)
	}
	return ec.JSON(http.StatusOK, models.BuildSuccessResponse(message, code, data))
}
