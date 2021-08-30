// Package subcmd provides types and functions for creating command-line interfaces with subcommands and flags.
package subcmd

import (
	"context"
	"flag"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/pkg/errors"
)

var errType = reflect.TypeOf((*error)(nil)).Elem()

// Cmd is the way a command tells Run how to parse and run its subcommands.
type Cmd interface {
	// Subcmds returns this Cmd's subcommands as a map,
	// whose keys are subcommand names and values are Subcmd objects.
	// Implementations may use the Commands function to build this map.
	Subcmds() Map
}

// Map is the type of the data structure returned by Cmd.Subcmds and by Commands.
// It maps a subcommand name to its Subcmd structure.
type Map = map[string]Subcmd

// Returns c's subcommand names as a sorted slice.
func subcmdNames(c Cmd) []string {
	var result []string
	for cmdname := range c.Subcmds() {
		result = append(result, cmdname)
	}
	sort.Strings(result)
	return result
}

// Subcmd is one subcommand of a Cmd.
type Subcmd struct {
	// F is the function implementing the subcommand.
	// Its signature must be func(context.Context, ..., []string) error,
	// where the number and types of parameters between the context and the string slice
	// is given by Params.
	// The error return is optional.
	F interface{}

	// Params describes the parameters to F
	// (excluding the initial context.Context that F takes, and the final []string).
	Params []Param

	// Desc is a one-line description of this subcommand.
	Desc string
}

// Param is one parameter of a Subcmd.
type Param struct {
	// Name is the flag name for the parameter
	// (e.g., "-verbose" for a -verbose flag,
	// "filename" for a positional parameter,
	// "optional?" for an optional positional parameter).
	Name string

	// Type is the type of the parameter.
	Type Type

	// Default is a default value for the parameter.
	// Its type must be suitable for Type.
	Default interface{}

	// Doc is a docstring for the parameter.
	Doc string
}

// Commands is a convenience function for producing the map needed by a Cmd.
// It takes 4n arguments,
// where n is the number of subcommands.
// Each group of three is:
// - the subcommand name, a string;
// - the function implementing the subcommand;
// - a short description of the subcommand;
// - the list of parameters for the function, a slice of Param (which can be produced with Params).
//
// A call like this:
//
//   Commands(
//     "foo", foo, "is the foo subcommand", Params(
//       "verbose", Bool, false, "be verbose",
//     ),
//     "bar", bar, "is the bar subcommand", Params(
//       "level", Int, 0, "barness level",
//     ),
//   )
//
// is equivalent to:
//
//   Map{
//     "foo": Subcmd{
//       F: foo,
//       Params: []Param{
//         {
//           Name: "verbose",
//           Type: Bool,
//           Default: false,
//           Doc: "be verbose",
//         },
//       },
//       Desc: "is the foo subcommand",
//     },
//     "bar": Subcmd{
//       F: bar,
//       Params: []Param{
//         {
//           Name: "level",
//           Type: Int,
//           Default: 0,
//           Doc: "barness level",
//         },
//       },
//       Desc: "is the bar subcommand",
//     },
//  }
//
// This function panics if the number or types of the arguments are wrong.
func Commands(args ...interface{}) Map {
	if len(args)%4 != 0 {
		panic(fmt.Sprintf("S has %d arguments, which is not divisible by 4", len(args)))
	}

	result := make(Map)

	for len(args) > 0 {
		var (
			name = args[0].(string)
			f    = args[1]
			d    = args[2].(string)
			p    = args[3]
		)
		subcmd := Subcmd{F: f, Desc: d}
		if p != nil {
			subcmd.Params = p.([]Param)
		}
		result[name] = subcmd

		args = args[4:]
	}

	return result
}

// Params is a convenience function for producing the list of parameters needed by a Subcmd.
// It takes 4n arguments,
// where n is the number of parameters.
// Each group of four is:
// - the flag name for the parameter, a string (e.g. "verbose" for a -verbose flag);
// - the type of the parameter, a Type constant;
// - the default value of the parameter,
// - the doc string for the parameter.
//
// This function panics if the number or types of the arguments are wrong.
func Params(a ...interface{}) []Param {
	if len(a)%4 != 0 {
		panic(fmt.Sprintf("Params has %d arguments, which is not divisible by 4", len(a)))
	}
	var result []Param
	for len(a) > 0 {
		var (
			name = a[0].(string)
			typ  = a[1].(Type)
			dflt = a[2]
			doc  = a[3].(string)
		)
		result = append(result, Param{Name: name, Type: typ, Default: dflt, Doc: doc})
		a = a[4:]
	}
	return result
}

// Run runs the subcommand of c named in args[0].
// The remaining args are parsed with a new flag.FlagSet,
// populated according to the parameters of the named Subcmd.
// The Subcmd's function is invoked with a context object,
// the parameter values parsed by the FlagSet,
// and a slice of the args left over after FlagSet parsing.
// The FlagSet is placed in the context object that's passed to the Subcmd's function,
// and can be retrieved if needed with the FlagSet function.
// No FlagSet is present if the subcommand takes no parameters.
func Run(ctx context.Context, c Cmd, args []string) error {
	if len(args) == 0 {
		return &MissingSubcmdErr{
			pairs: subcmdPairList(ctx),
			cmd:   c,
		}
	}

	cmds := c.Subcmds()

	name := args[0]
	args = args[1:]
	subcmd, ok := cmds[name]

	if !ok && name == "help" {
		e := &HelpRequestedErr{
			pairs: subcmdPairList(ctx),
			cmd:   c,
		}
		if len(args) > 0 {
			e.name = args[0]
		}
		return e
	}
	if !ok {
		return &UnknownSubcmdErr{
			pairs: subcmdPairList(ctx),
			cmd:   c,
			name:  name,
		}
	}

	ctx = addSubcmdPair(ctx, name, subcmd)

	argvals, err := parseArgs(ctx, subcmd.Params, args)
	if err != nil {
		return err
	}

	fv := reflect.ValueOf(subcmd.F)
	ft := fv.Type()
	if ft.Kind() != reflect.Func {
		return fmt.Errorf("implementation for subcommand %s is a %s, want a function", name, ft.Kind())
	}
	if numIn := ft.NumIn(); numIn != len(argvals) {
		return fmt.Errorf("function for subcommand %s takes %d arg(s), want %d", name, numIn, len(argvals))
	}
	for i, argval := range argvals {
		if !argval.Type().AssignableTo(ft.In(i)) {
			return fmt.Errorf("type of arg %d is %s, want %s", i, ft.In(i), argval.Type())
		}
	}

	numOut := ft.NumOut()
	if numOut > 1 {
		return fmt.Errorf("function for subcommand %s returns %d values, want 0 or 1", name, numOut)
	}
	if numOut == 1 && !ft.Out(0).Implements(errType) {
		return fmt.Errorf("return type is not error")
	}

	rv := fv.Call(argvals)

	if numOut == 1 {
		err, _ = rv[0].Interface().(error)
	}

	return errors.Wrapf(err, "running %s", name)
}

// ToFlagSet produces a *flag.FlagSet from the given params,
// plus a list of properly typed pointers in which to store the result of calling Parse on the FlagSet.
func ToFlagSet(params []Param) (*flag.FlagSet, []reflect.Value, []Param, error) {
	var (
		fs   = flag.NewFlagSet("", flag.ContinueOnError)
		ptrs []reflect.Value
	)

	for len(params) > 0 {
		p := params[0]
		if !strings.HasPrefix(p.Name, "-") {
			break
		}
		name := p.Name[1:]
		params = params[1:]

		var v interface{}

		switch p.Type {
		case Bool:
			dflt, _ := p.Default.(bool)
			v = fs.Bool(name, dflt, p.Doc)

		case Int:
			dflt, _ := p.Default.(int)
			v = fs.Int(name, dflt, p.Doc)

		case Int64:
			dflt, _ := p.Default.(int64)
			v = fs.Int64(name, dflt, p.Doc)

		case Uint:
			dflt, _ := p.Default.(uint)
			v = fs.Uint(name, dflt, p.Doc)

		case Uint64:
			dflt, _ := p.Default.(uint64)
			v = fs.Uint64(name, dflt, p.Doc)

		case String:
			dflt, _ := p.Default.(string)
			v = fs.String(name, dflt, p.Doc)

		case Float64:
			dflt, _ := p.Default.(float64)
			v = fs.Float64(name, dflt, p.Doc)

		case Duration:
			dflt, _ := p.Default.(time.Duration)
			v = fs.Duration(name, dflt, p.Doc)

		default:
			return nil, nil, nil, fmt.Errorf("unknown arg type %v", p.Type)
		}

		ptrs = append(ptrs, reflect.ValueOf(v))
	}

	return fs, ptrs, params, nil
}

// Type is the type of a Param.
type Type int

// Possible Param types.
// These correspond with the types in the standard flag package.
const (
	Bool Type = iota + 1
	Int
	Int64
	Uint
	Uint64
	String
	Float64
	Duration
)
