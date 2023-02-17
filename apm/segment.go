package apm

type (
	Segment interface {
		AddAttribute(key string, val interface{})
		End()
	}
)
