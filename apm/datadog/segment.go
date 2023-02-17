package datadog

import (
	"github.com/Dert12318/Utilities/apm"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

type (
	segment struct {
		app  apm.APM
		span tracer.Span
	}
)

func (s *segment) AddAttribute(key string, val interface{}) {
	s.span.SetTag(key, val)
}

func (s *segment) End() {
	s.span.Finish()
}
