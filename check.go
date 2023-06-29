package subcmd

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/pkg/errors"
)

// Check performs type checking on subcmd as follows:
//
//  - It checks that subcmd.F is a function and returns ErrNotAFunction if it isn't.
//  - It checks that subcmd.F returns no more than one value and returns ErrTooManyReturns if it doesn't.
//  - It checks that the type of the value returned by subcmd.F (if any) is error and returns ErrNotError if it isn't.
//  - It checks that subcmd.F takes an initial context.Context parameter and returns ErrNoContext if it doesn't.
//  - It checks that subcmd.F takes a final []string parameter and returns ErrNoStringSlice if it doesn't.
//  - It checks that the length of subcmd.Params matches the number of parameters subcmd.F takes (not counting the initial context.Context and final []string parameters) and returns a NumParamsErr if it doesn't.
//  - It checks that each parameter in subcmd.Params matches the corresponding parameter in subcmd.F and returns a ParamTypeErr if it doesn't.
//  - It checks that the default value of each parameter in subcmd.Params matches the parameter's type and returns a ParamDefaultErr if it doesn't.
//
// Only the first error encountered is returned.
func Check(subcmd Subcmd) error {
	fv := reflect.ValueOf(subcmd.F)
	ft := fv.Type()
	if ft.Kind() != reflect.Func {
		return ErrNotAFunction
	}
	numIn := ft.NumIn()
	if numIn != len(subcmd.Params)+2 {
		return NumParamsErr{Want: len(subcmd.Params) + 2, Got: numIn}
	}
	if err := checkParam(ft, 0, contextType); err != nil {
		return errors.Wrap(err, "checking parameter 0")
	}
	if err := checkParam(ft, numIn-1, stringSliceType); err != nil {
		return errors.Wrap(err, "checking last parameter")
	}

	for i, param := range subcmd.Params {
		if err := checkParam(ft, i+1, param.Type); err != nil {
			return errors.Wrapf(err, "checking parameter %d", i+1)
		}
	}

	numOut := ft.NumOut()
	switch numOut {
	case 0:
		// ok
	case 1:
		if !ft.Out(0).Implements(errType) {
			return ErrNotError
		}
	default:
		return ErrTooManyReturns
	}

	return nil
}

var (
	// ErrNotAFunction is returned by Check if the F field of a Subcmd is not a function.
	ErrNotAFunction = errors.New("not a function")

	// ErrNotError is returned by Check if the F field of a Subcmd returns a value that is not an error.
	ErrNotError = errors.New("function returns non-error")

	// ErrTooManyReturns is returned by Check if the F field of a Subcmd returns more than one value.
	ErrTooManyReturns = errors.New("function returns too many values")

	// ErrNoContext is returned by Check if the F field of a Subcmd does not take an initial context.Context parameter.
	ErrNoContext = errors.New("parameter 0 is not context.Context")

	// ErrNoStringSlice is returned by Check if the F field of a Subcmd does not take a final []string parameter.
	ErrNoStringSlice = errors.New("last parameter is not []string")
)

// CheckMap calls Check on each of the entries in the Map.
func CheckMap(m Map) error {
	for name, subcmd := range m {
		if err := Check(subcmd); err != nil {
			return errors.Wrapf(err, "checking subcommand %s", name)
		}
	}
	return nil
}

func checkParam(ft reflect.Type, n int, want Type) error {
	var (
		paramType = ft.In(n)
		errType   reflect.Type
	)

	switch want {
	case Bool:
		var (
			x  bool
			xt = reflect.TypeOf(x)
		)
		if !xt.AssignableTo(paramType) {
			errType = xt
		}

	case Int:
		var (
			x  int
			xt = reflect.TypeOf(x)
		)
		if !xt.AssignableTo(paramType) {
			errType = xt
		}

	case Int64:
		var (
			x  int64
			xt = reflect.TypeOf(x)
		)
		if !xt.AssignableTo(paramType) {
			errType = xt
		}

	case Uint:
		var (
			x  uint
			xt = reflect.TypeOf(x)
		)
		if !xt.AssignableTo(paramType) {
			errType = xt
		}

	case Uint64:
		var (
			x  uint64
			xt = reflect.TypeOf(x)
		)
		if !xt.AssignableTo(paramType) {
			errType = xt
		}

	case String:
		var (
			x  string
			xt = reflect.TypeOf(x)
		)
		if !xt.AssignableTo(paramType) {
			errType = xt
		}

	case Float64:
		var (
			x  float64
			xt = reflect.TypeOf(x)
		)
		if !xt.AssignableTo(paramType) {
			errType = xt
		}

	case Duration:
		var (
			x  time.Duration
			xt = reflect.TypeOf(x)
		)
		if !xt.AssignableTo(paramType) {
			errType = xt
		}

	case contextType:
		var (
			x  = context.Background() // Need a concrete value, not a nil interface.
			xt = reflect.TypeOf(x)
		)
		if !xt.AssignableTo(paramType) {
			// Special case.
			return ErrNoContext
		}

	case stringSliceType:
		var (
			x  []string
			xt = reflect.TypeOf(x)
		)
		if !xt.AssignableTo(paramType) {
			// Special case.
			return ErrNoStringSlice
		}

	default:
		return fmt.Errorf("unknown type %v", want)
	}

	if errType != nil {
		return ParamTypeErr{N: n, Want: want, Got: errType}
	}

	return nil
}

// ParamTypeErr is returned by Check if the type of a parameter in a Subcmd's function doesn't match the corresponding Param.Type field.
type ParamTypeErr struct {
	// N is the parameter number. The initial context.Context is parameter 0.
	N int

	// Want is the type specified in the Param.Type field.
	Want Type

	// Got is the type of the parameter in the function.
	Got reflect.Type
}

func (e ParamTypeErr) Error() string {
	return fmt.Sprintf("parameter %d has type %v, want %v", e.N, e.Got, e.Want)
}

// NumParamsErr is returned by Check if the number of parameters in a Subcmd's function doesn't match the number of Param fields
// (plus two for the initial context.Context and final []string).
type NumParamsErr struct {
	Want, Got int
}

func (e NumParamsErr) Error() string {
	return fmt.Sprintf("function has %d parameters, want %d", e.Got, e.Want)
}
