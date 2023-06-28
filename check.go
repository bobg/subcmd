package subcmd

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/pkg/errors"
)

// Check checks that the type of subcmd.F matches the list of parameters in subcmd.Params.
func Check(subcmd Subcmd) error {
	fv := reflect.ValueOf(subcmd.F)
	ft := fv.Type()
	if ft.Kind() != reflect.Func {
		return fmt.Errorf("F is not a function")
	}
	numIn := ft.NumIn()
	if numIn != len(subcmd.Params)+2 {
		return fmt.Errorf("F has %d parameters, want %d", numIn, len(subcmd.Params)+2)
	}

	if err := checkParam(ft, 0, contextType); err != nil {
		return errors.Wrap(err, "checking first parameter is a context.Context")
	}
	if err := checkParam(ft, numIn-1, stringSliceType); err != nil {
		return errors.Wrap(err, "checking last parameter is []string")
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
			return fmt.Errorf("return type is not error")
		}
	default:
		return fmt.Errorf("F returns %d values, want 0 or 1", numOut)
	}

	return nil
}

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
	paramType := ft.In(n)

	switch want {
	case Bool:
		var x bool
		if !reflect.TypeOf(x).AssignableTo(paramType) {
			return fmt.Errorf("parameter %d is %s, want bool", n, paramType)
		}

	case Int:
		var x int
		if !reflect.TypeOf(x).AssignableTo(paramType) {
			return fmt.Errorf("parameter %d is %s, want int", n, paramType)
		}

	case Int64:
		var x int64
		if !reflect.TypeOf(x).AssignableTo(paramType) {
			return fmt.Errorf("parameter %d is %s, want int64", n, paramType)
		}

	case Uint:
		var x uint
		if !reflect.TypeOf(x).AssignableTo(paramType) {
			return fmt.Errorf("parameter %d is %s, want uint", n, paramType)
		}

	case Uint64:
		var x uint64
		if !reflect.TypeOf(x).AssignableTo(paramType) {
			return fmt.Errorf("parameter %d is %s, want uint64", n, paramType)
		}

	case String:
		var x string
		if !reflect.TypeOf(x).AssignableTo(paramType) {
			return fmt.Errorf("parameter %d is %s, want string", n, paramType)
		}

	case Float64:
		var x float64
		if !reflect.TypeOf(x).AssignableTo(paramType) {
			return fmt.Errorf("parameter %d is %s, want float64", n, paramType)
		}

	case Duration:
		var x time.Duration
		if !reflect.TypeOf(x).AssignableTo(paramType) {
			return fmt.Errorf("parameter %d is %s, want time.Duration", n, paramType)
		}

	case contextType:
		x := context.Background() // Need a concrete value, not a nil interface.
		if !reflect.TypeOf(x).AssignableTo(paramType) {
			return fmt.Errorf("parameter %d is %s, want context.Context", n, paramType)
		}

	case stringSliceType:
		var x []string
		if !reflect.TypeOf(x).AssignableTo(paramType) {
			return fmt.Errorf("parameter %d is %s, want []string", n, paramType)
		}

	default:
		return fmt.Errorf("unknown type %v", want)
	}

	return nil
}
