package datadog

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/event"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"

	"github.com/Dert12318/Utilities/apm"
	"github.com/Dert12318/Utilities/apm/disabled"
	"github.com/Dert12318/Utilities/logs"
)

type (
	Option struct {
		AppName    string
		Logger     logs.Logger
		ActiveSpan bool
		DebugMode  bool
		Env        string
	}
	datadog struct {
		option Option
	}
)

func (dd *datadog) CommandMonitor() *event.CommandMonitor {
	return nil
}

func (dd *datadog) FromContext(ctx context.Context) apm.Transaction {
	span, _ := tracer.SpanFromContext(ctx)
	return &transaction{
		app:  dd,
		span: span,
	}
}

func (dd *datadog) RecordCustomEvent(eventType string, params map[string]interface{}) {

}

func (dd *datadog) StartTransaction(transactionName string) apm.Transaction {
	span := tracer.StartSpan(transactionName)
	return &transaction{
		app:  dd,
		span: span,
	}
}

func (dd *datadog) Shutdown(duration time.Duration) {
	tracer.Stop()
}

func New(option Option) (apm.APM, error) {
	tracer.Start(
		tracer.WithServiceName(option.AppName),
		tracer.WithLogger(option.Logger),
		tracer.WithDebugMode(option.DebugMode),
		tracer.WithEnv(option.Env),
	)

	if option.ActiveSpan {
		return &datadog{
			option: option,
		}, nil
	}

	return disabled.New()
}
