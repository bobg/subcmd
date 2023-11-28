package subcmd

import (
	"flag"
	"reflect"

	"github.com/pkg/errors"
)

// Check checks that the type of subcmd.F matches the expectations set by subcmd.Params:
//
//   - It must be a function;
//   - It must return no more than one value;
//   - If it returns a value, that value must be of type error;
//   - It must take an initial context.Context parameter;
//   - It must take a final []string or ...string parameter;
//   - The length of subcmd.Params must match the number of parameters subcmd.F takes (not counting the initial context.Context and final []string parameters);
//   - Each parameter in subcmd.Params must match the corresponding parameter in subcmd.F.
//
// It also checks that the default value of each parameter in subcmd.Params matches the parameter's type.
func Check(subcmd Subcmd) error {
	fv := reflect.ValueOf(subcmd.F)
	ft := fv.Type()

	if err := checkFuncType(ft, subcmd.Params); err != nil {
		return err
	}

	for i, param := range subcmd.Params {
		if err := checkParam(param); err != nil {
			return errors.Wrapf(err, "checking parameter %d", i+1)
		}
	}

	return nil
}

func checkFuncType(ft reflect.Type, params []Param) error {
	in := make([]reflect.Type, 0, 2+len(params))
	in = append(in, ctxType)
	for _, param := range params {
		in = append(in, param.Type.reflectType())
	}
	in = append(in, strSliceType)

	out := []reflect.Type{errType}

	if ft == reflect.FuncOf(in, nil, true) {
		return nil
	}
	if ft == reflect.FuncOf(in, out, true) {
		return nil
	}
	if ft == reflect.FuncOf(in, nil, false) {
		return nil
	}
	if want := reflect.FuncOf(in, out, false); ft != want {
		return FuncTypeErr{Got: ft, Want: want}
	}
	return nil
}

func checkParam(param Param) error {
	switch param.Type {
	case Bool:
		if _, ok := param.Default.(bool); !ok {
			return ParamDefaultErr{Param: param}
		}

	case Int:
		if _, err := asInt(param.Default); err != nil {
			return ParamDefaultErr{Param: param}
		}

	case Int64:
		if _, err := asInt64(param.Default); err != nil {
			return ParamDefaultErr{Param: param}
		}

	case Uint:
		if _, err := asUint(param.Default); err != nil {
			return ParamDefaultErr{Param: param}
		}

	case Uint64:
		if _, err := asUint64(param.Default); err != nil {
			return ParamDefaultErr{Param: param}
		}

	case String:
		if _, ok := param.Default.(string); !ok {
			return ParamDefaultErr{Param: param}
		}

	case Float64:
		if _, err := asFloat64(param.Default); err != nil {
			return ParamDefaultErr{Param: param}
		}

	case Duration:
		if _, err := asDuration(param.Default); err != nil {
			return ParamDefaultErr{Param: param}
		}

	case Value:
		if _, ok := param.Default.(flag.Value); !ok {
			return ParamDefaultErr{Param: param}
		}
	}

	return nil
}

// CheckMap calls [Check] on each of the entries in the Map.
func CheckMap(m Map) error {
	for name, subcmd := range m {
		if err := Check(subcmd); err != nil {
			return errors.Wrapf(err, "checking subcommand %s", name)
		}
	}
	return nil
}
