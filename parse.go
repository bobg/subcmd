package subcmd

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

func parseArgs(ctx context.Context, params []Param, args []string) ([]reflect.Value, error) {
	fs, ptrs, positional, err := ToFlagSet(params)
	if err != nil {
		return nil, err
	}

	err = fs.Parse(args)
	if err != nil {
		return nil, errors.Wrap(err, "parsing args")
	}

	args = fs.Args()
	ctx = withFlagSet(ctx, fs)

	argvals := []reflect.Value{reflect.ValueOf(ctx)}
	for _, ptr := range ptrs {
		argvals = append(argvals, ptr.Elem())
	}

	for _, p := range positional {
		err = parsePositionalArg(p, &args, &argvals)
		if err != nil {
			return nil, err
		}
	}

	argvals = append(argvals, reflect.ValueOf(args))

	return argvals, nil
}

func parsePositionalArg(p Param, args *[]string, argvals *[]reflect.Value) error {
	if len(*args) == 0 && !strings.HasSuffix(p.Name, "?") {
		return ErrTooFewArgs
	}

	switch p.Type {
	case Bool:
		return parseBoolPos(args, argvals)

	case Int:
		return parseIntPos(args, argvals)

	case Int64:
		return parseInt64Pos(args, argvals)

	case Uint:
		return parseUintPos(args, argvals)

	case Uint64:
		return parseUint64Pos(args, argvals)

	case String:
		return parseStringPos(args, argvals)

	case Float64:
		return parseFloat64Pos(args, argvals)

	case Duration:
		return parseDurationPos(args, argvals)

	default:
		return fmt.Errorf("unknown arg type %v", p.Type)
	}
}

func parseBoolPos(args *[]string, argvals *[]reflect.Value) error {
	var val bool
	if len(*args) > 0 {
		var err error
		val, err = strconv.ParseBool((*args)[0])
		if err != nil {
			return ParseErr{Err: err}
		}
		*args = (*args)[1:]
	}
	*argvals = append(*argvals, reflect.ValueOf(val))
	return nil
}

func parseIntPos(args *[]string, argvals *[]reflect.Value) error {
	var val int64
	if len(*args) > 0 {
		var err error
		val, err = strconv.ParseInt((*args)[0], 10, 32)
		if err != nil {
			return ParseErr{Err: err}
		}
		*args = (*args)[1:]
	}
	*argvals = append(*argvals, reflect.ValueOf(int(val)))
	return nil
}

func parseInt64Pos(args *[]string, argvals *[]reflect.Value) error {
	var val int64
	if len(*args) > 0 {
		var err error
		val, err = strconv.ParseInt((*args)[0], 10, 64)
		if err != nil {
			return ParseErr{Err: err}
		}
		*args = (*args)[1:]
	}
	*argvals = append(*argvals, reflect.ValueOf(val))
	return nil
}

func parseUintPos(args *[]string, argvals *[]reflect.Value) error {
	var val uint64
	if len(*args) > 0 {
		var err error
		val, err = strconv.ParseUint((*args)[0], 10, 32)
		if err != nil {
			return ParseErr{Err: err}
		}
		*args = (*args)[1:]
	}
	*argvals = append(*argvals, reflect.ValueOf(uint(val)))
	return nil
}

func parseUint64Pos(args *[]string, argvals *[]reflect.Value) error {
	var val uint64
	if len(*args) > 0 {
		var err error
		val, err = strconv.ParseUint((*args)[0], 10, 64)
		if err != nil {
			return ParseErr{Err: err}
		}
		*args = (*args)[1:]
	}
	*argvals = append(*argvals, reflect.ValueOf(val))
	return nil
}

func parseStringPos(args *[]string, argvals *[]reflect.Value) error {
	var val string
	if len(*args) > 0 {
		val = (*args)[0]
		*args = (*args)[1:]
	}
	*argvals = append(*argvals, reflect.ValueOf(val))
	return nil
}

func parseFloat64Pos(args *[]string, argvals *[]reflect.Value) error {
	var val float64
	if len(*args) > 0 {
		var err error
		val, err = strconv.ParseFloat((*args)[0], 64)
		if err != nil {
			return ParseErr{Err: err}
		}
		*args = (*args)[1:]
	}
	*argvals = append(*argvals, reflect.ValueOf(val))
	return nil
}

func parseDurationPos(args *[]string, argvals *[]reflect.Value) error {
	var val time.Duration
	if len(*args) > 0 {
		var err error
		val, err = time.ParseDuration((*args)[0])
		if err != nil {
			return ParseErr{Err: err}
		}
		*args = (*args)[1:]
	}
	*argvals = append(*argvals, reflect.ValueOf(val))
	return nil
}
