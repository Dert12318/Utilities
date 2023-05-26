package errors

import "github.com/pkg/errors"

var (
	MissingHandler = errors.New("handler is missing")
)
