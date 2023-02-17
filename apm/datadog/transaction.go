package datadog

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"

	"github.com/Dert12318/Utilities/apm"
)

type (
	transaction struct {
		app  apm.APM
		span tracer.Span
	}
)

func (t *transaction) Application() apm.APM {
	return t.app
}

func (t *transaction) End() {
	t.span.Finish()
}

func (t *transaction) Ignore() {}

func (t *transaction) SetName(name string) {}

func (t *transaction) NoticeError(err error) {}

func (t *transaction) AddAttribute(key string, value interface{}) {
	t.span.SetTag(key, value)
}

func (t *transaction) SetWebRequestHTTP(r *http.Request) {}

type dummyResponseWriter struct{}

func (rw dummyResponseWriter) Header() http.Header { return nil }

func (rw dummyResponseWriter) Write(b []byte) (int, error) { return 0, nil }

func (rw dummyResponseWriter) WriteHeader(code int) {}

func (t *transaction) SetWebResponse(w http.ResponseWriter) http.ResponseWriter {
	if w == nil {
		return dummyResponseWriter{}
	}
	return w
}

func (t *transaction) StartSegment(name string) apm.Segment {
	childSpan := tracer.StartSpan(name, tracer.ChildOf(t.span.Context()))
	return &segment{
		app:  t.app,
		span: childSpan,
	}
}

func (t *transaction) StartDataStoreSegment(dto apm.DatastoreSegmentDTO) apm.Segment {
	childSpan := tracer.StartSpan(
		fmt.Sprintf("%s#%s", dto.DatabaseName, dto.Operation),
		tracer.ChildOf(t.span.Context()))
	return &segment{
		app:  t.app,
		span: childSpan,
	}
}

func (t *transaction) StartMessageProducerSegment(request apm.MessageProducerSegmentDTO) apm.Segment {
	childSpan := tracer.StartSpan(
		fmt.Sprintf("%s#%s", request.DestinationName, request.DestinationType),
		tracer.ChildOf(t.span.Context()))
	return &segment{
		app:  t.app,
		span: childSpan,
	}
}

func (t *transaction) StartExternalSegment(request *http.Request) apm.Segment {
	childSpan := tracer.StartSpan(
		fmt.Sprintf("%s#%s", request.Method, request.URL),
		tracer.ChildOf(t.span.Context()))
	return &segment{
		app:  t.app,
		span: childSpan,
	}
}

func (t *transaction) InsertDistributedTraceHeaders(header http.Header) {

}

func (t *transaction) NewContext(ctx context.Context) context.Context {
	return tracer.ContextWithSpan(ctx, t.span)
}

func (t *transaction) GetTraceID() string {
	return strconv.Itoa(int(t.span.Context().TraceID()))
}
