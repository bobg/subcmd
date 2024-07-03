package subcmd

import (
	"context"
	"flag"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// If variadic is false, the length of the resulting slice is len(params)+2.
// If it's true, the length is >= len(params)+1.
func parseArgs(ctx context.Context, params []Param, args []string, variadic bool) ([]reflect.Value, error) {
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
		if ptr.Type().Implements(valueType) {
			argvals = append(argvals, ptr)
		} else {
			argvals = append(argvals, ptr.Elem())
		}
	}

	for _, p := range positional {
		err = parsePositionalArg(p, &args, &argvals)
		if err != nil {
			return nil, err
		}
	}

	if variadic {
		for _, arg := range args {
			argvals = append(argvals, reflect.ValueOf(arg))
		}
	} else {
		argvals = append(argvals, reflect.ValueOf(args))
	}

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

	case Value:
		return parseValuePos(args, argvals, p)

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
	val := int64(asInt(p.Default)) // has to be int64 to receive the result of ParseInt below

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
	val := asInt64(p.Default)

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
	val := uint64(asUint(p.Default)) // has to be uint64 to receive the result of ParseUint below

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
	val := asUint64(p.Default)

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
	val := asFloat64(p.Default)

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
	val := asDuration(p.Default)

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

func parseValuePos(args *[]string, argvals *[]reflect.Value, p Param) error {
	val, ok := p.Default.(ValueType)
	if !ok {
		return ParseErr{Err: fmt.Errorf("param %s is not a ValueType", p.Name)}
	}
	val = val.Copy()
	if len(*args) > 0 {
		if err := val.Set((*args)[0]); err != nil {
			return ParseErr{Err: err}
		}
		*args = (*args)[1:]
	}
	*argvals = append(*argvals, reflect.ValueOf(val))
	return nil
}

func asInt(val interface{}) int {
	switch v := val.(type) {
	case int:
		return v

	case int8:
		return int(v)
	case int16:
		return int(v)
	case int32:
		return int(v)

	case uint8:
		return int(v)
	case uint16:
		return int(v)
	}
	return 0
}

func asInt64(val interface{}) int64 {
	switch v := val.(type) {
	case int:
		return int64(v)
	case int64:
		return v

	case int8:
		return int64(v)
	case int16:
		return int64(v)
	case int32:
		return int64(v)

	case uint8:
		return int64(v)
	case uint16:
		return int64(v)
	case uint32:
		return int64(v)
	}
	return 0
}

func asUint(val interface{}) uint {
	switch v := val.(type) {
	case uint:
		return v

	case int8:
		return uint(v)
	case int16:
		return uint(v)

	case uint8:
		return uint(v)
	case uint16:
		return uint(v)
	case uint32:
		return uint(v)
	}
	return 0
}

func asUint64(val interface{}) uint64 {
	switch v := val.(type) {
	case uint:
		return uint64(v)
	case uint64:
		return v

	case int8:
		return uint64(v)
	case int16:
		return uint64(v)
	case int32:
		return uint64(v)

	case uint8:
		return uint64(v)
	case uint16:
		return uint64(v)
	case uint32:
		return uint64(v)
	}
	return 0
}

func asFloat64(val interface{}) float64 {
	switch v := val.(type) {
	case int:
		return float64(v)
	case uint:
		return float64(v)

	case int8:
		return float64(v)
	case int16:
		return float64(v)
	case int32:
		return float64(v)
	case int64:
		return float64(v)

	case uint8:
		return float64(v)
	case uint16:
		return float64(v)
	case uint32:
		return float64(v)
	case uint64:
		return float64(v)

	case float32:
		return float64(v)
	case float64:
		return v
	}
	return 0
}

func asDuration(val interface{}) time.Duration {
	switch v := val.(type) {
	case int:
		return time.Duration(v)
	case int64:
		return time.Duration(v)
	case time.Duration:
		return v
	}
	return 0
}

// ToFlagSet takes a slice of [Param] and produces:
//
//   - a [flag.FlagSet],
//   - a list of properly typed pointers (or in the case of a [Value]-typed Param, a [ValueType]) in which to store the results of calling Parse on the FlagSet,
//   - a list of positional Params that are not part of the resulting FlagSet.
//
// On a successful return, len(ptrs)+len(positional) == len(params).
func ToFlagSet(params []Param) (fs *flag.FlagSet, ptrs []reflect.Value, positional []Param, err error) {
	fs = flag.NewFlagSet("", flag.ContinueOnError)

	for _, p := range params {
		if !strings.HasPrefix(p.Name, "-") {
			positional = append(positional, p)
			continue
		}

		var (
			name = strings.TrimLeft(p.Name, "-")
			v    interface{}
		)

		switch p.Type {
		case Bool:
			dflt, _ := p.Default.(bool)
			v = fs.Bool(name, dflt, p.Doc)

		case Int:
			v = fs.Int(name, asInt(p.Default), p.Doc)

		case Int64:
			v = fs.Int64(name, asInt64(p.Default), p.Doc)

		case Uint:
			v = fs.Uint(name, asUint(p.Default), p.Doc)

		case Uint64:
			v = fs.Uint64(name, asUint64(p.Default), p.Doc)

		case String:
			dflt, _ := p.Default.(string)
			v = fs.String(name, dflt, p.Doc)

		case Float64:
			v = fs.Float64(name, asFloat64(p.Default), p.Doc)

		case Duration:
			v = fs.Duration(name, asDuration(p.Default), p.Doc)

		case Value:
			val, ok := p.Default.(ValueType)
			if !ok {
				err = fmt.Errorf("param %s has type Value but default value %v is not a ValueType", p.Name, p.Default)
				return
			}
			val = val.Copy()
			fs.Var(val, name, p.Doc)
			v = val

		default:
			err = fmt.Errorf("unknown arg type %v", p.Type)
			return
		}

		ptrs = append(ptrs, reflect.ValueOf(v))
	}

	return fs, ptrs, positional, nil
}
