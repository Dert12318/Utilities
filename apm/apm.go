package apm

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/event"
)

type (
	APM interface {
		FromContext(ctx context.Context) Transaction
		StartTransaction(transactionName string) Transaction
		RecordCustomEvent(eventType string, params map[string]interface{})
		Shutdown(duration time.Duration)
		CommandMonitor() *event.CommandMonitor
	}
)
