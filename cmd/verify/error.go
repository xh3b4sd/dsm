package verify

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

var invalidFlagError = &tracer.Error{
	Kind: "invalidFlagError",
}

func IsInvalidFlag(err error) bool {
	return errors.Is(err, invalidFlagError)
}

var invalidValueError = &tracer.Error{
	Kind: "invalidValueError",
}

func IsInvalidValue(err error) bool {
	return errors.Is(err, invalidValueError)
}

var notFoundError = &tracer.Error{
	Kind: "notFoundError",
	Desc: "When verifying the consistency of values across multiple files, there must be at least one file found. This error is caused by no files being found given the provided flags. Check if there are typos in the query and that the command is being executed against the correct directory.",
}

func IsNotFound(err error) bool {
	return errors.Is(err, notFoundError)
}
