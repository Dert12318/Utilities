package apm

import (
	"context"
	"net/http"

	relic "github.com/newrelic/go-agent/v3/newrelic"
)

const (
	MessageQueue    MessageDestinationType = "Queue"
	MessageTopic    MessageDestinationType = "Topic"
	MessageExchange MessageDestinationType = "Exchange"
)

type (
	MessageDestinationType string

	MessageProducerSegmentDTO struct {
		Library              string
		DestinationType      MessageDestinationType
		DestinationName      string
		DestinationTemporary bool
	}

	DatastoreSegmentDTO struct {
		Collection         string
		Operation          string
		ParameterizedQuery string
		QueryParameters    []interface{}
		DatabaseName       string
		DatastoreProduct   relic.DatastoreProduct
	}

	Transaction interface {
		Application() APM
		End()
		Ignore()
		SetName(name string)
		NoticeError(err error)
		AddAttribute(key string, value interface{})
		SetWebRequestHTTP(r *http.Request)
		SetWebResponse(w http.ResponseWriter) http.ResponseWriter
		StartSegment(name string) Segment
		StartDataStoreSegment(dto DatastoreSegmentDTO) Segment
		StartMessageProducerSegment(request MessageProducerSegmentDTO) Segment
		StartExternalSegment(request *http.Request) Segment
		InsertDistributedTraceHeaders(header http.Header)
		NewContext(ctx context.Context) context.Context
		GetTraceID() string
	}
)
