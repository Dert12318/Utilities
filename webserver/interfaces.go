package webserver

import (
	"time"

	"github.com/labstack/echo/v4"

	"github.com/Dert12318/Utilities/logs"
	"github.com/Dert12318/Utilities/validator"
)

type (
	Infrastructure interface {
		Register(ec *echo.Echo) error
		Listen() error
	}

	Resource interface {
		Echo() *echo.Echo
		ServiceName() string
		ServerPort() int
		EnableProfiler() bool
		Validator() validator.Validator
		ServerGracefullyDuration() time.Duration
		Logger() logs.Logger
		Close() error
	}

	IWebServer interface {
		Serve() error
	}
)
