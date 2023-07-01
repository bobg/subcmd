package subcmd

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"reflect"
	"strings"
)

// ErrTooFewArgs is the error when not enough arguments are supplied for required positional parameters.
var ErrTooFewArgs = errors.New("too few arguments")

// ParseErr is the type of error returned when parsing a positional parameter according to its type fails.
type ParseErr struct {
	Err error
}

func (e ParseErr) Error() string {
	return "parse error: " + e.Err.Error()
}

// Unwrap unwraps the nested error in e.
func (e ParseErr) Unwrap() error {
	return e.Err
}

// UsageErr is the type of errors that give usage information.
// Such errors have the usual Error() method producing a one-line string,
// but also a Detail() method producing a multiline string with more detail.
type UsageErr interface {
	error
	Detail() string
}

// MissingSubcmdErr is a usage error returned when Run is called with an empty `args` list.
type MissingSubcmdErr struct {
	pairs []subcmdPair
	cmd   Cmd
}

func (e *MissingSubcmdErr) Error() string {
	return fmt.Sprintf("missing subcommand, want one of: %s", strings.Join(subcmdNames(e.cmd), "; "))
}

// Detail implements Usage.
func (e *MissingSubcmdErr) Detail() string {
	return missingUnknownSubcmd("Missing subcommand, want one of:", e.cmd)
}

// HelpRequestedErr is a usage error returned when the "help" pseudo-subcommand-name is used.
type HelpRequestedErr struct {
	pairs []subcmdPair
	cmd   Cmd
	name  string
}

func (e *HelpRequestedErr) Error() string {
	if e.name != "" {
		// foo bar help baz
		subcmds := e.cmd.Subcmds()
		subcmd, ok := subcmds[e.name]
		if !ok {
			return fmt.Sprintf(`unknown subcommand "%s", want one of: %s`, e.name, strings.Join(subcmdNames(e.cmd), "; "))
		}

		fs, _, positional, err := ToFlagSet(subcmd.Params)
		if err != nil {
			return fmt.Sprintf("error constructing usage string: %s", err.Error())
		}

		b := new(strings.Builder)
		fmt.Fprintf(b, "usage: %s", os.Args[0])
		for _, pair := range e.pairs {
			fmt.Fprint(b, " ", pair.name)
		}
		fmt.Fprintf(b, " %s", e.name)

		fs.VisitAll(func(f *flag.Flag) {
			if name, _ := flag.UnquoteUsage(f); name == "" {
				fmt.Fprintf(b, " [-%s]", f.Name)
			} else {
				fmt.Fprintf(b, " [-%s %s]", f.Name, name)
			}
		})
		for _, p := range positional {
			name := p.Name
			if strings.HasSuffix(name, "?") {
				fmt.Fprintf(b, " [%s]", name[:len(name)-1])
			} else {
				fmt.Fprint(b, " ", name)
			}
		}
		return b.String()
	}

	// foo bar help
	return fmt.Sprintf("subcommands are: %s", strings.Join(subcmdNames(e.cmd), "; "))
}

// Detail implements Usage.
func (e *HelpRequestedErr) Detail() string {
	if e.name != "" {
		// foo bar help baz
		subcmds := e.cmd.Subcmds()
		subcmd, ok := subcmds[e.name]
		if !ok {
			return fmt.Sprintf(`unknown subcommand "%s", want one of: %s`, e.name, strings.Join(subcmdNames(e.cmd), "; "))
		}

		fs, _, positional, err := ToFlagSet(subcmd.Params)
		if err != nil {
			return fmt.Sprintf("error constructing usage string: %s", err.Error())
		}

		b := new(strings.Builder)

		if subcmd.Desc != "" {
			fmt.Fprintf(b, "%s: %s\n", e.name, subcmd.Desc)
		}

		fmt.Fprintf(b, "Usage: %s", os.Args[0])
		for _, pair := range e.pairs {
			fmt.Fprint(b, " ", pair.name)
		}
		fmt.Fprintf(b, " %s", e.name)

		var maxlen int
		fs.VisitAll(func(f *flag.Flag) {
			var l int
			if name, _ := flag.UnquoteUsage(f); name == "" {
				fmt.Fprintf(b, " [-%s]", f.Name)
				l = len(f.Name)
			} else {
				fmt.Fprintf(b, " [-%s %s]", f.Name, name)
				l = 1 + len(f.Name) + len(name)
			}
			if l > maxlen {
				maxlen = l
			}
		})
		for _, p := range positional {
			name := p.Name
			if strings.HasSuffix(name, "?") {
				fmt.Fprintf(b, " [%s]", name[:len(name)-1])
			} else {
				fmt.Fprint(b, " ", name)
			}
		}
		fmt.Fprintln(b)

		format := fmt.Sprintf("-%%-%d.%ds  %%s\n", maxlen, maxlen)

		fs.VisitAll(func(f *flag.Flag) {
			if name, u := flag.UnquoteUsage(f); name == "" {
				fmt.Fprintf(b, format, f.Name, u)
			} else {
				fmt.Fprintf(b, format, f.Name+" "+name, u)
			}
		})

		return b.String()
	}

	// foo bar help
	b := new(strings.Builder)
	fmt.Fprintln(b, "Subcommands are:")
	cmdnames := subcmdNames(e.cmd)
	subcmds := e.cmd.Subcmds()
	var maxlen int
	for _, name := range cmdnames {
		if len(name) > maxlen {
			maxlen = len(name)
		}
	}
	format := fmt.Sprintf("%%-%d.%ds  %%s\n", maxlen, maxlen)
	for _, name := range cmdnames {
		fmt.Fprintf(b, format, name, subcmds[name].Desc)
	}

	return b.String()
}

// UnknownSubcmdErr is a usage error returned when an unknown subcommand name is passed to Run as args[0].
type UnknownSubcmdErr struct {
	pairs []subcmdPair
	cmd   Cmd
	name  string
}

func (e *UnknownSubcmdErr) Error() string {
	return fmt.Sprintf(`unknown subcommand "%s", want one of: %s`, e.name, strings.Join(subcmdNames(e.cmd), "; "))
}

// Detail implements Usage.
func (e *UnknownSubcmdErr) Detail() string {
	return missingUnknownSubcmd(fmt.Sprintf(`Unknown subcommand "%s", want one of:`, e.name), e.cmd)
}

func missingUnknownSubcmd(line1 string, cmd Cmd) string {
	b := new(strings.Builder)
	fmt.Fprintln(b, line1)
	cmdnames := subcmdNames(cmd)
	subcmds := cmd.Subcmds()
	var maxlen int
	for _, name := range cmdnames {
		if len(name) > maxlen {
			maxlen = len(name)
		}
	}
	format := fmt.Sprintf("%%-%d.%ds  %%s\n", maxlen, maxlen)
	for _, name := range cmdnames {
		fmt.Fprintf(b, format, name, subcmds[name].Desc)
	}
	return b.String()
}

// FuncTypeErr means a Subcmd's F field has a type that does not match the function signature implied by its Params field.
type FuncTypeErr struct {
	// Got is the type of the F field.
	Got reflect.Type

	// Want is the expected function type implied by the Params field.
	// Note: for simplicity, this includes the optional error return,
	// even if the type in Got does not
	// (which is not, in itself, an error).
	Want reflect.Type
}

func (e FuncTypeErr) Error() string {
	return fmt.Sprintf("function has type %v, want %v", e.Got, e.Want)
}

// NumArgsErr is the error when too many or too few arguments are supplied to Run for a Subcmd's function.
type NumArgsErr struct {
	// Got is the number of arguments supplied.
	Got int

	// Want is the number of function parameters expected.
	Want int
}

func (e NumArgsErr) Error() string {
	return fmt.Sprintf("got %d arguments but function takes %d parameters", e.Got, e.Want)
}

// ParamDefaultErr is the error when a Param has a default value that is not of the correct type.
type ParamDefaultErr struct {
	Param Param
}

func (e ParamDefaultErr) Error() string {
	return fmt.Sprintf("default value %v is not of type %v", e.Param.Default, e.Param.Type)
}
