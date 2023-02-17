package disabled

import (
	"context"
	"net/http"

	"github.com/Dert12318/Utilities/apm"
)

type (
	transaction struct {
		app apm.APM
	}
)

func (txn *transaction) Application() apm.APM {
	return txn.app
}

func (txn *transaction) End() {}

func (txn *transaction) Ignore() {}

func (txn *transaction) SetName(name string) {}

func (txn *transaction) NoticeError(err error) {}

func (txn *transaction) AddAttribute(key string, value interface{}) {}

func (txn *transaction) SetWebRequestHTTP(r *http.Request) {}

type dummyResponseWriter struct{}

func (rw dummyResponseWriter) Header() http.Header { return nil }

func (rw dummyResponseWriter) Write(b []byte) (int, error) { return 0, nil }

func (rw dummyResponseWriter) WriteHeader(code int) {}

func (txn *transaction) SetWebResponse(w http.ResponseWriter) http.ResponseWriter {
	if w == nil {
		return dummyResponseWriter{}
	}
	return w
}

func (txn *transaction) StartSegment(name string) apm.Segment {
	return &segment{}
}

func (txn *transaction) StartDataStoreSegment(dto apm.DatastoreSegmentDTO) apm.Segment {
	return &segment{}
}

func (txn *transaction) StartMessageProducerSegment(request apm.MessageProducerSegmentDTO) apm.Segment {
	return &segment{}
}

func (txn *transaction) StartExternalSegment(request *http.Request) apm.Segment {
	return &segment{}
}

func (txn *transaction) InsertDistributedTraceHeaders(header http.Header) {}

func (txn *transaction) NewContext(ctx context.Context) context.Context {
	return ctx
}
