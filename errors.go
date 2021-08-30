package subcmd

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

var (
	// ErrNoArgs is the error when Run is called with an empty list of args.
	ErrNoArgs = errors.New("no arguments")

	// ErrUnknown is the error when Run is called with an unknown subcommand as args[0].
	ErrUnknown = errors.New("unknown subcommand")

	// ErrTooFewArgs is the error when not enough arguments are supplied to populate required positional parameters.
	ErrTooFewArgs = errors.New("too few arguments")
)

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

// Usage is the type of errors that give usage information.
// Such errors have the usual Error() method producing a one-line string,
// but also a Long() method producing a multiline string with more detail.
type Usage interface {
	Long() string
}

// MissingSubcmdErr is a usage error returned when Run is called with an empty `args` list.
type MissingSubcmdErr struct {
	pairs []subcmdPair
	cmd   Cmd
}

func (e *MissingSubcmdErr) Error() string {
	return fmt.Sprintf("missing subcommand, want one of: %s", strings.Join(subcmdNames(e.cmd), "; "))
}

// Long implements Usage.
func (e *MissingSubcmdErr) Long() string {
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

// Long implements Usage.
func (e *HelpRequestedErr) Long() string {
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

// UnknownSubcmdErr is a usage error returned when an unknown subcommand name is passed to Run.
type UnknownSubcmdErr struct {
	pairs []subcmdPair
	cmd   Cmd
	name  string
}

func (e *UnknownSubcmdErr) Error() string {
	return fmt.Sprintf(`unknown subcommand "%s", want one of: %s`, e.name, strings.Join(subcmdNames(e.cmd), "; "))
}

// Long implements Usage.
func (e *UnknownSubcmdErr) Long() string {
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
