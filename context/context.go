package context

import (
	"context"

	"github.com/labstack/echo/v4"

	"github.com/Dert12318/Utilities/apm"
)

type (
	Context struct {
		ec          echo.Context
		Ctx         context.Context
		Transaction apm.Transaction
		mandatory   MandatoryRequest
	}
)

func New() *Context {
	return &Context{
		Ctx: context.Background(),
	}
}

func NewWithContext(ctx context.Context) *Context {
	ec := echo.New().NewContext(nil, nil)
	return &Context{
		Ctx: ctx,
		ec:  ec,
	}
}

func NewWithEcho(ec echo.Context) *Context {
	ctx := &Context{
		Ctx: context.Background(),
		ec:  ec,
	}
	return ctx
}

func NewWithEchoAndContext(ec echo.Context, ctx context.Context) *Context {
	c := &Context{
		Ctx: ctx,
		ec:  ec,
	}
	return c
}

func (c Context) MandatoryRequest() MandatoryRequest {
	return c.mandatory
}
