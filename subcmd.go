// Package subcmd provides types and functions for creating command-line interfaces with subcommands and flags.
package subcmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"sort"

	"github.com/pkg/errors"
)

var errType = reflect.TypeOf((*error)(nil)).Elem()

// Cmd is a command that has subcommands.
// It tells Run how to parse its subcommands, and their flags and positional parameters,
// and how to run them.
type Cmd interface {
	// Subcmds returns this Cmd's subcommands as a map,
	// whose keys are subcommand names and values are Subcmd objects.
	// The Commands() function is useful in building this map.
	Subcmds() Map
}

// Prefixer is an optional additional interface that a Cmd can implement.
// If it does, and a call to Run encounters an unknown subcommand,
// then before returning an error it will look for an executable in $PATH
// whose name is Prefix() plus the subcommand name.
// If it finds one,
// it is executed with the remaining args as arguments,
// and a JSON-marshaled copy of the Cmd in the environment variable SUBCMD_ENV
// (that can be parsed by the subprocess using ParseEnv).
type Prefixer interface {
	Prefix() string
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
	// Name is the flag name for the parameter.
	// Flags must have a leading "-", as in "-verbose".
	// Positional parameters have no leading "-".
	// Optional positional parameters have a trailing "?", as in "optional?".
	Name string

	// Type is the type of the parameter.
	Type Type

	// Default is a default value for the parameter.
	// Its type must be suitable for Type.
	Default interface{}

	// Doc is a docstring for the parameter.
	Doc string
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

	// The following are for internal use only.
	contextType
	stringSliceType
)

// String returns the name of t.
func (t Type) String() string {
	switch t {
	case Bool:
		return "bool"
	case Int:
		return "int"
	case Int64:
		return "int64"
	case Uint:
		return "uint"
	case Uint64:
		return "uint64"
	case String:
		return "string"
	case Float64:
		return "float64"
	case Duration:
		return "time.Duration"
	case contextType:
		return "context.Context"
	case stringSliceType:
		return "[]string"
	default:
		return fmt.Sprintf("unknown type %d", t)
	}
}

// Commands is a convenience function for producing the Map
// needed by an implementation of Cmd.Subcmd.
// It takes arguments in groups of two or four,
// one group per subcommand.
//
// The first argument of a group is the subcommand's name, a string.
// The second argument of a group may be a Subcmd,
// making this a two-argument group.
// If it's not a Subcmd,
// then this is a four-argument group,
// whose second through fourth arguments are:
//
//   - the function implementing the subcommand;
//   - a short description of the subcommand;
//   - the list of parameters for the function, a slice of Param (which can be produced with the Params function).
//
// These are used to populate a Subcmd.
//
// A call like this:
//
//	Commands(
//	  "foo", foo, "is the foo subcommand", Params(
//	    "-verbose", Bool, false, "be verbose",
//	  ),
//	  "bar", bar, "is the bar subcommand", Params(
//	    "-level", Int, 0, "barness level",
//	  ),
//	)
//
// is equivalent to:
//
//	 Map{
//	   "foo": Subcmd{
//	     F:      foo,
//	     Desc:   "is the foo subcommand",
//	     Params: []Param{
//	       {
//	         Name:    "-verbose",
//	         Type:    Bool,
//	         Default: false,
//	         Doc:     "be verbose",
//	       },
//	     },
//	   },
//	   "bar": Subcmd{
//	     F:      bar,
//	     Desc:   "is the bar subcommand",
//	     Params: []Param{
//	       {
//	         Name:    "-level",
//	         Type:    Int,
//	         Default: 0,
//	         Doc:     "barness level",
//	       },
//	     },
//	   },
//	}
//
// This function panics if the number or types of the arguments are wrong.
func Commands(args ...interface{}) Map {
	result := make(Map)

	for len(args) > 0 {
		if len(args) < 2 {
			panic(fmt.Errorf("too few arguments to Commands"))
		}

		name := args[0].(string)
		if subcmd, ok := args[1].(Subcmd); ok {
			result[name] = subcmd
			args = args[2:]
			continue
		}

		if len(args) < 4 {
			panic(fmt.Errorf("too few arguments to Commands"))
		}

		var (
			f = args[1]
			d = args[2].(string)
			p = args[3]
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
// - the name for the parameter, a string (e.g. "-verbose" for a -verbose flag);
// - the type of the parameter, a Type constant;
// - the default value of the parameter,
// - the doc string for the parameter.
//
// This function panics if the number or types of the arguments are wrong.
func Params(a ...interface{}) []Param {
	if len(a)%4 != 0 {
		panic(fmt.Sprintf("Params called with %d arguments, which is not divisible by 4", len(a)))
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
//
// That subcommand specifies zero or more flags and zero or more positional parameters.
// The remaining values in args are parsed to populate those.
//
// After argument parsing,
// the subcommand's function is invoked with a context object,
// the flag and parameter values,
// and a slice of the args remaining.
//
// Flags are parsed using a new flag.FlagSet,
// which is placed into the context object passed to the subcommand's function.
// The FlagSet can be retrieved if needed with the FlagSet function.
// No flag.FlagSet is present if the subcommand has no flags.
//
// Flags are always optional, and have names beginning with "-".
// Positional parameters may be required or optional.
// Optional positional parameters have a trailing "?" in their names.
//
// Calling Run with an empty args slice produces a MissingSubcmdErr error.
//
// Calling Run with an unknown subcommand name in args[0] produces an UnknownSubcmdErr error,
// unless the unknown subcommand is "help",
// in which case the result is a HelpRequestedErr,
// or unless c is also a Prefixer.
//
// If c is a Prefixer and the subcommand name is both unknown and not "help",
// then an executable is sought in $PATH with c's prefix plus the subcommand name.
// If one is found,
// it is executed with the remaining args as arguments,
// and a JSON-marshaled copy of c in the environment variable SUBCMD_ENV
// (that can be parsed by the subprocess using ParseEnv).
//
// If there are not enough values in args to populate the subcommand's required positional parameters,
// the result is ErrTooFewArgs.
//
// If argument parsing succeeds,
// Run returns the error produced by calling the subcommand's function, if any.
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
		unknownSubcmdErr := &UnknownSubcmdErr{
			pairs: subcmdPairList(ctx),
			cmd:   c,
			name:  name,
		}

		if p, ok := c.(Prefixer); ok {
			// The cmds map does not contain name,
			// but c is a Prefixer so look for the executable prefix+name to run instead.

			prefix := p.Prefix()
			path, err := exec.LookPath(prefix + name)
			if errors.Is(err, exec.ErrNotFound) {
				return unknownSubcmdErr
			}
			if err != nil {
				return errors.Wrapf(err, "looking for %s%s", prefix, name)
			}

			execCmd := exec.CommandContext(ctx, path, args...)
			execCmd.Stdin, execCmd.Stdout, execCmd.Stderr = os.Stdin, os.Stdout, os.Stderr

			j, err := json.Marshal(c)
			if err != nil {
				return errors.Wrap(err, "marshaling Cmd")
			}
			execCmd.Env = append(os.Environ(), EnvVar+"="+string(j))

			return execCmd.Run()
		}

		return unknownSubcmdErr
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

// EnvVar is the name of the environment variable used by Run to pass the JSON-encoded Cmd to a subprocess.
// Use ParseEnv to decode it.
// See Prefixer.
const EnvVar = "SUBCMD_ENV"

// ParseEnv parses the value of the SUBCMD_ENV environment variable,
// placing the result in the value pointed to by ptr,
// which must be a pointer of a suitable type.
// Executables that implement subcommands should run this at startup.
func ParseEnv(ptr interface{}) error {
	val := os.Getenv(EnvVar)
	if val == "" {
		return nil
	}
	return json.Unmarshal([]byte(val), ptr)
}
