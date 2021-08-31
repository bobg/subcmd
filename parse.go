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
		return parseBoolPos(args, argvals, p)

	case Int:
		return parseIntPos(args, argvals, p)

	case Int64:
		return parseInt64Pos(args, argvals, p)

	case Uint:
		return parseUintPos(args, argvals, p)

	case Uint64:
		return parseUint64Pos(args, argvals, p)

	case String:
		return parseStringPos(args, argvals, p)

	case Float64:
		return parseFloat64Pos(args, argvals, p)

	case Duration:
		return parseDurationPos(args, argvals, p)

	default:
		return fmt.Errorf("unknown arg type %v", p.Type)
	}
}

func parseBoolPos(args *[]string, argvals *[]reflect.Value, p Param) error {
	val, _ := p.Default.(bool)
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

func parseIntPos(args *[]string, argvals *[]reflect.Value, p Param) error {
	var val int64
	if dflt, ok := p.Default.(int); ok {
		val = int64(dflt)
	} else if dflt, ok := p.Default.(int64); ok {
		val = dflt
	}

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

func parseInt64Pos(args *[]string, argvals *[]reflect.Value, p Param) error {
	var val int64
	if dflt, ok := p.Default.(int); ok {
		val = int64(dflt)
	} else if dflt, ok := p.Default.(int64); ok {
		val = dflt
	}

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

func parseUintPos(args *[]string, argvals *[]reflect.Value, p Param) error {
	var val uint64
	if dflt, ok := p.Default.(uint); ok {
		val = uint64(dflt)
	} else if dflt, ok := p.Default.(uint64); ok {
		val = dflt
	}

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

func parseUint64Pos(args *[]string, argvals *[]reflect.Value, p Param) error {
	var val uint64
	if dflt, ok := p.Default.(uint); ok {
		val = uint64(dflt)
	} else if dflt, ok := p.Default.(uint64); ok {
		val = dflt
	}

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

func parseStringPos(args *[]string, argvals *[]reflect.Value, p Param) error {
	val, _ := p.Default.(string)
	if len(*args) > 0 {
		val = (*args)[0]
		*args = (*args)[1:]
	}
	*argvals = append(*argvals, reflect.ValueOf(val))
	return nil
}

func parseFloat64Pos(args *[]string, argvals *[]reflect.Value, p Param) error {
	var val float64
	switch dflt := p.Default.(type) {
	case int:
		val = float64(dflt)
	case int64:
		val = float64(dflt)
	case uint:
		val = float64(dflt)
	case uint64:
		val = float64(dflt)
	case float32:
		val = float64(dflt)
	case float64:
		val = dflt
	}

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

func parseDurationPos(args *[]string, argvals *[]reflect.Value, p Param) error {
	var val time.Duration
	switch dflt := p.Default.(type) {
	case int:
		val = time.Duration(dflt)
	case int64:
		val = time.Duration(dflt)
	case time.Duration:
		val = dflt
	}

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
