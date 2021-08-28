package subcmd

import (
	"errors"
	"fmt"
	"strings"
)

var (
	// ErrNoArgs is the error when Run is called with an empty list of args.
	ErrNoArgs = errors.New("no arguments")

	// ErrUnknown is the error when Run is called with an unknown subcommand as args[0].
	ErrUnknown = errors.New("unknown subcommand")

	// ErrTooFewArgs is the error when not enough arguments are supplied to populate required positional parameters.
	ErrTooFewArgs = errors.New("too few arguments")
)

type UsageErr struct {
	names    []string
	defaults string
}

func (e *UsageErr) Error() string {
	names := make([]string, len(e.names))
	for i := 0; i < len(e.names); i++ {
		names[i] = e.names[len(e.names)-1-i]
	}
	return fmt.Sprintf("usage: %s [flags] ...\n%s", strings.Join(names, " "), e.defaults)
}

type ParseErr struct {
	Err error
}

func (e ParseErr) Error() string {
	return "parse error: " + e.Err.Error()
}

func (e ParseErr) Unwrap() error {
	return e.Err
}
