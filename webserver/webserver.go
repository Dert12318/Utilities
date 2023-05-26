package webserver

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"gopkg.in/DataDog/dd-trace-go.v1/profiler"

	"github.com/Dert12318/Utilities/logs/logrus"
)

type (
	WebServer struct {
		infrastructure Infrastructure
		resource       Resource

		afterRun   []HookFunction
		beforeRun  []HookFunction
		afterExit  []HookFunction
		beforeExit []HookFunction
	}
)

func New(infrastructure Infrastructure, resource Resource) IWebServer {
	return &WebServer{
		infrastructure: infrastructure,
		resource:       resource,
		afterRun:       make([]HookFunction, 0),
		beforeRun:      make([]HookFunction, 0),
		afterExit:      make([]HookFunction, 0),
		beforeExit:     make([]HookFunction, 0),
	}
}

func (w *WebServer) Serve() error {
	err := w.initialize()
	if err != nil {
		return err
	}

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)

	if err := w.applyHooks(w.beforeRun); err != nil {
		return err
	}

	go w.serve(sig)

	if err := w.applyHooks(w.afterRun); err != nil {
		return err
	}
	<-sig

	if err := w.applyHooks(w.beforeExit); err != nil {
		return err
	}
	w.exit()
	if err := w.applyHooks(w.afterExit); err != nil {
		return err
	}
	return nil
}

func (w *WebServer) initialize() error {
	w.resource.Echo().Use(middleware.Recover())
	w.resource.Echo().Validator = w.resource.Validator()

	w.resource.Echo().Logger = logrus.DefaultLog()
	w.resource.Echo().Logger.SetLevel(log.INFO)

	w.resource.Echo().Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAcceptEncoding, echo.HeaderAuthorization},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete, http.MethodOptions},
	}))

	w.resource.Echo().GET("/", func(context echo.Context) error {
		return context.JSON(http.StatusOK, w.resource.ServiceName())
	})

	if err := w.infrastructure.Register(w.resource.Echo()); err != nil {
		w.resource.Logger().Error("error on register http")
		return err
	}

	if err := w.infrastructure.Listen(); err != nil {
		w.resource.Logger().Error("error on register listener")
		return err
	}

	return nil
}

func (w *WebServer) serve(sig chan os.Signal) {
	if w.resource.EnableProfiler() {
		if err := profiler.Start(profiler.WithService(w.resource.ServiceName())); err != nil {
			w.resource.Logger().Error("failed to start profiler")
		}
	}

	if err := w.resource.Echo().Start(fmt.Sprintf(":%d", w.resource.ServerPort())); err != nil {
		w.resource.Logger().Errorf("http server interrupted %s", err.Error())
		sig <- syscall.SIGINT
	} else {
		w.resource.Logger().Info("starting apps")
	}
}

func (w *WebServer) applyHooks(hooks []HookFunction) error {
	for i := 0; i < len(hooks); i++ {
		if err := hooks[i](w.resource); err != nil {
			return err
		}
	}
	return nil
}

func (w *WebServer) exit() {
	ctx, cancel := context.WithTimeout(context.Background(), w.resource.ServerGracefullyDuration())
	defer cancel()

	if err := w.resource.Echo().Shutdown(ctx); err != nil {
		w.resource.Logger().Error("failed to shutdown echo http server %s", err)
	}

	w.resource.Logger().Info("closing resource")
	if err := w.resource.Close(); err != nil {
		w.resource.Logger().Error("failed to close resource %s", err)
	}
}
