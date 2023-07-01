package subcmd

import (
	"reflect"

	"github.com/pkg/errors"
)

// Check checks that the type of subcmd.F matches the expectations set by subcmd.Params:
//
//   - It must be a function;
//   - It must return no more than one value;
//   - If it returns a value, that value must be of type error;
//   - It must take an initial context.Context parameter;
//   - It must take a final []string parameter;
//   - The length of subcmd.Params must match the number of parameters subcmd.F takes (not counting the initial context.Context and final []string parameters);
//   - Each parameter in subcmd.Params must match the corresponding parameter in subcmd.F.
//
// It also checks that the default value of each parameter in subcmd.Params matches the parameter's type.
func Check(subcmd Subcmd) error {
	fv := reflect.ValueOf(subcmd.F)
	ft := fv.Type()

	wantFuncTypeErr, wantFuncTypeNoErr := funcTypeForParams(subcmd.Params)
	switch ft {
	case wantFuncTypeErr, wantFuncTypeNoErr:
		// ok

	default:
		return FuncTypeErr{Want: wantFuncTypeErr, Got: ft}
	}

	for i, param := range subcmd.Params {
		if err := checkParam(param); err != nil {
			return errors.Wrapf(err, "checking parameter %d", i+1)
		}
	}

	return nil
}

func checkParam(param Param) error {
	if !reflect.TypeOf(param.Default).AssignableTo(param.Type.reflectType()) {
		return ParamDefaultErr{Param: param}
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

func funcTypeForParams(params []Param) (withErr, withoutErr reflect.Type) {
	in := []reflect.Type{ctxType}
	for _, param := range params {
		in = append(in, param.Type.reflectType())
	}
	in = append(in, reflect.TypeOf([]string(nil)))

	withoutErr = reflect.FuncOf(in, nil, false)

	out := []reflect.Type{errType}

	withErr = reflect.FuncOf(in, out, false)

	return withErr, withoutErr
}
