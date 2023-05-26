package messaging

import (
	"context"
	"sync"

	tntContext "github.com/Dert12318/Utilities/context"
	constantError "github.com/Dert12318/Utilities/errors"
)

type singleEventDispatcher struct {
	handler      HandlerFunc
	errorHandler ErrorHandlerFunc
	middlewares  []MiddlewareFunc
	sync.Mutex
}

func (d *singleEventDispatcher) AddHandler(handler HandlerFunc, errorHandler ErrorHandlerFunc, msgType ...string) {
	d.handler = handler
	d.errorHandler = errorHandler
}

func (d *singleEventDispatcher) Dispatch(dto DispatchDTO) error {
	dto.Log.Debugf("RECEIVE[%v][%v] %v", dto.Msg.MsgID, dto.RequestID, string(dto.Msg.MsgData))
	dispatch := applyMiddleware(d.dispatch, d.middlewares...)
	return dispatch(tntContext.NewWithContext(context.Background()), dto)
}

func (d *singleEventDispatcher) Use(middlewareFunc ...MiddlewareFunc) {
	d.Lock()
	defer d.Unlock()
	d.middlewares = append(d.middlewares, middlewareFunc...)
}

func NewSingleEventDispatcher() Dispatcher {
	return &singleEventDispatcher{
		middlewares: make([]MiddlewareFunc, 0),
	}
}

func (d *singleEventDispatcher) dispatch(ctx *tntContext.Context, dto DispatchDTO) error {
	if ctx == nil {
		c := context.Background()
		ctx = tntContext.NewWithContext(c)
	}

	if dto.Type == Handle {
		return d.handle(ctx, dto.Msg)
	}
	return d.onError(ctx, dto.Msg, dto.Err)
}

func (d *singleEventDispatcher) handle(ctx *tntContext.Context, msg Message) error {
	if d.handler == nil {
		return constantError.MissingHandler
	}
	return d.handler(ctx, msg)
}

func (d *singleEventDispatcher) onError(ctx *tntContext.Context, msg Message, err error) error {
	if d.errorHandler == nil {
		return constantError.MissingHandler
	}
	d.errorHandler(ctx, msg, err)
	return nil
}
