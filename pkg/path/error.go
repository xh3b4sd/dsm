package path

import (
	"errors"

	"github.com/xh3b4sd/tracer"
)

var invalidConfigError = &tracer.Error{
	Kind: "invalidConfigError",
}

func IsInvalidConfig(err error) bool {
	return errors.Is(err, invalidConfigError)
}

var invalidFormatError = &tracer.Error{
	Kind: "invalidFormatError",
}

func IsInvalidFormat(err error) bool {
	return errors.Is(err, invalidFormatError)
}

var notFoundError = &tracer.Error{
	Kind: "notFoundError",
}

func IsNotFound(err error) bool {
	return errors.Is(err, notFoundError)
}
