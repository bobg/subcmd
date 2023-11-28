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
	val, ok := p.Default.(bool)
	if !ok {
		return ParseErr{Err: fmt.Errorf("param %s has type Bool but default value %v has type %T", p.Name, p.Default, p.Default)}
	}
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
	valInt, err := asInt(p.Default)
	if err != nil {
		return ParseErr{Err: err}
	}
	val := int64(valInt) // has to be int64 to receive the result of ParseInt below

	if len(*args) > 0 {
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
	val, err := asInt64(p.Default)
	if err != nil {
		return ParseErr{Err: err}
	}

	if len(*args) > 0 {
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
	valUint, err := asUint(p.Default)
	if err != nil {
		return ParseErr{Err: err}
	}
	val := uint64(valUint) // has to be uint64 to receive the result of ParseUint below

	if len(*args) > 0 {
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
	val, err := asUint64(p.Default)
	if err != nil {
		return ParseErr{Err: err}
	}

	if len(*args) > 0 {
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
	val, ok := p.Default.(string)
	if !ok {
		return ParseErr{Err: fmt.Errorf("param %s has type String but default value %v has type %T", p.Name, p.Default, p.Default)}
	}
	if len(*args) > 0 {
		val = (*args)[0]
		*args = (*args)[1:]
	}
	*argvals = append(*argvals, reflect.ValueOf(val))
	return nil
}

func parseFloat64Pos(args *[]string, argvals *[]reflect.Value, p Param) error {
	val, err := asFloat64(p.Default)
	if err != nil {
		return ParseErr{Err: err}
	}

	if len(*args) > 0 {
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
	val, err := asDuration(p.Default)
	if err != nil {
		return ParseErr{Err: err}
	}

	if len(*args) > 0 {
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
	val, ok := p.Default.(flag.Value)
	if !ok {
		return ParseErr{Err: fmt.Errorf("param %s is not a flag.Value", p.Name)}
	}
	if len(*args) > 0 {
		if err := val.Set((*args)[0]); err != nil {
			return ParseErr{Err: err}
		}
		*args = (*args)[1:]
	}
	*argvals = append(*argvals, reflect.ValueOf(val))
	return nil
}

func asInt(val interface{}) (int, error) {
	switch v := val.(type) {
	case int:
		return v, nil

	case int8:
		return int(v), nil
	case int16:
		return int(v), nil
	case int32:
		return int(v), nil

	case uint8:
		return int(v), nil
	case uint16:
		return int(v), nil
	}

	return 0, fmt.Errorf("cannot convert %v (type %T) to int", val, val)
}

func asInt64(val interface{}) (int64, error) {
	switch v := val.(type) {
	case int:
		return int64(v), nil
	case int64:
		return v, nil

	case int8:
		return int64(v), nil
	case int16:
		return int64(v), nil
	case int32:
		return int64(v), nil

	case uint8:
		return int64(v), nil
	case uint16:
		return int64(v), nil
	case uint32:
		return int64(v), nil
	}

	return 0, fmt.Errorf("cannot convert %v (type %T) to int64", val, val)
}

func asUint(val interface{}) (uint, error) {
	switch v := val.(type) {
	case int:
		if v >= 0 {
			return uint(v), nil
		}

	case uint:
		return v, nil

	case int8:
		if v >= 0 {
			return uint(v), nil
		}

	case int16:
		if v >= 0 {
			return uint(v), nil
		}

	case int32:
		if v >= 0 {
			return uint(v), nil
		}

	case uint8:
		return uint(v), nil
	case uint16:
		return uint(v), nil
	case uint32:
		return uint(v), nil
	}

	return 0, fmt.Errorf("cannot convert %v (type %T) to uint", val, val)
}

func asUint64(val interface{}) (uint64, error) {
	switch v := val.(type) {
	case int:
		if v >= 0 {
			return uint64(v), nil
		}

	case uint:
		return uint64(v), nil
	case uint64:
		return v, nil

	case int8:
		if v >= 0 {
			return uint64(v), nil
		}

	case int16:
		if v >= 0 {
			return uint64(v), nil
		}

	case int32:
		if v >= 0 {
			return uint64(v), nil
		}

	case int64:
		if v >= 0 {
			return uint64(v), nil
		}

	case uint8:
		return uint64(v), nil
	case uint16:
		return uint64(v), nil
	case uint32:
		return uint64(v), nil
	}

	return 0, fmt.Errorf("cannot convert %v (type %T) to uint64", val, val)
}

func asFloat64(val interface{}) (float64, error) {
	switch v := val.(type) {
	case int:
		return float64(v), nil
	case uint:
		return float64(v), nil

	case int8:
		return float64(v), nil
	case int16:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case int64:
		return float64(v), nil

	case uint8:
		return float64(v), nil
	case uint16:
		return float64(v), nil
	case uint32:
		return float64(v), nil
	case uint64:
		return float64(v), nil

	case float32:
		return float64(v), nil
	case float64:
		return v, nil
	}

	return 0, fmt.Errorf("cannot convert %v (type %T) to float64", val, val)
}

func asDuration(val interface{}) (time.Duration, error) {
	if v, ok := val.(time.Duration); ok {
		return v, nil
	}

	v, err := asInt64(val)
	if err != nil {
		return 0, fmt.Errorf("cannot convert %v (type %T) to time.Duration", val, val)
	}
	return time.Duration(v), nil
}

// ToFlagSet takes a slice of [Param] and produces:
//
//   - a [flag.FlagSet],
//   - a list of properly typed pointers (or in the case of a [Value]-typed Param, a [flag.Value]) in which to store the results of calling Parse on the FlagSet,
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
			dflt, err := asInt(p.Default)
			if err != nil {
				return nil, nil, nil, err
			}
			v = fs.Int(name, dflt, p.Doc)

		case Int64:
			dflt, err := asInt64(p.Default)
			if err != nil {
				return nil, nil, nil, err
			}
			v = fs.Int64(name, dflt, p.Doc)

		case Uint:
			dflt, err := asUint(p.Default)
			if err != nil {
				return nil, nil, nil, err
			}
			v = fs.Uint(name, dflt, p.Doc)

		case Uint64:
			dflt, err := asUint64(p.Default)
			if err != nil {
				return nil, nil, nil, err
			}
			v = fs.Uint64(name, dflt, p.Doc)

		case String:
			dflt, ok := p.Default.(string)
			if !ok {
				return nil, nil, nil, fmt.Errorf("param %s has type String but default value %v has type %T", p.Name, p.Default, p.Default)
			}
			v = fs.String(name, dflt, p.Doc)

		case Float64:
			dflt, err := asFloat64(p.Default)
			if err != nil {
				return nil, nil, nil, err
			}
			v = fs.Float64(name, dflt, p.Doc)

		case Duration:
			dflt, err := asDuration(p.Default)
			if err != nil {
				return nil, nil, nil, err
			}
			v = fs.Duration(name, dflt, p.Doc)

		case Value:
			val, ok := p.Default.(flag.Value)
			if !ok {
				err = fmt.Errorf("param %s has type Value but default value %v has type %T", p.Name, p.Default, p.Default)
				return
			}
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
