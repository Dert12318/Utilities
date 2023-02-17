package disabled

import (
	"context"
	"time"

	"github.com/Dert12318/Utilities/apm"
	"go.mongodb.org/mongo-driver/event"
)

type (
	disabled struct{}
)

func (nr *disabled) CommandMonitor() *event.CommandMonitor {
	return nil
}

func (nr *disabled) FromContext(ctx context.Context) apm.Transaction {
	return &transaction{}
}

func (nr *disabled) StartTransaction(transactionName string) apm.Transaction {
	return &transaction{}
}

func (nr *disabled) RecordCustomEvent(eventType string, params map[string]interface{}) {}

func (nr *disabled) Shutdown(duration time.Duration) {}

func New() (apm.APM, error) {
	return &disabled{}, nil
}

func (t *transaction) GetTraceID() string {
	return ""
}
