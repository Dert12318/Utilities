package response

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/Dert12318/Utilities/models/v2"
)

var (
	ErrBadRequest     = NewError(http.StatusText(http.StatusBadRequest), http.StatusBadRequest, nil)
	ErrUnauthorized   = NewError(http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized, nil)
	ErrNotFound       = NewError(http.StatusText(http.StatusNotFound), http.StatusNotFound, nil)
	ErrInternalServer = NewError(http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError, nil)
	ErrForbidden      = NewError(http.StatusText(http.StatusForbidden), http.StatusForbidden, nil)
)

type HttpError struct {
	models.ErrorResponse
	err error
}

func (e *HttpError) Error() string {
	return e.Message
}

func NewError(message string, status int, err error) *HttpError {
	if message == "" {
		message = http.StatusText(status)
	}
	if err == nil {
		err = errors.New("")
	}
	return &HttpError{
		ErrorResponse: models.BuildErrorResponse(message, status, err),
		err:           err,
	}
}

func ErrorWrap(base *HttpError, err error) *HttpError {
	if base == nil {
		base = ErrInternalServer
	}
	return NewError(base.Message, base.Code, err)
}

func ErrorWithMessage(base *HttpError, message string, err error) *HttpError {
	if base == nil {
		base = ErrInternalServer
	}
	return NewError(message, base.Code, err)
}

func ErrorWithErrMessage(base *HttpError, err error) *HttpError {
	if base == nil {
		base = ErrInternalServer
	}
	message := ""
	if err != nil {
		message = err.Error()
	}
	return NewError(message, base.Code, err)
}

func (r *HttpResponse) ErrorResponse(ec echo.Context, err error, request ...interface{}) error {
	httpErr, ok := err.(*HttpError)
	if !ok {
		httpErr = ErrorWrap(ErrInternalServer, err)
	}

	var req interface{}
	if len(request) > 0 {
		req = request[0]
	}

	r.logger.Error("", zap.Error(httpErr.err),
		zap.String("method", ec.Request().Method),
		zap.String("uri", ec.Request().RequestURI),
		zap.Any("request", req))

	return ec.JSON(httpErr.Code, models.BuildErrorResponse(httpErr.Message, httpErr.Code, httpErr.err))
}

func (r *HttpResponse) CustomErrorResponse(ec echo.Context, err error, request ...interface{}) error {
	httpErr, ok := err.(*HttpError)
	if !ok {
		httpErr = ErrorWrap(ErrInternalServer, err)
	}

	var req interface{}
	if len(request) > 0 {
		req = request[0]
	}

	r.logger.Error("", zap.Error(httpErr.err),
		zap.String("method", ec.Request().Method),
		zap.String("uri", ec.Request().RequestURI),
		zap.Any("request", req))

	customCode := GetErrorCode(httpErr.err)

	return ec.JSON(httpErr.Code, models.BuildCustomError(httpErr.Message, customCode, httpErr.Code, httpErr.err))
}
