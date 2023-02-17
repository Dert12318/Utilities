package disabled

import "github.com/Dert12318/Utilities/apm"

type (
	segment struct{}
)

func (s *segment) AddAttribute(key string, val interface{}) {}

func (s *segment) End() {}

func NewSegment() apm.Segment {
	return &segment{}
}
