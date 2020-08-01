package subcmd

import (
	"context"
	"flag"
	"sort"

	"github.com/pkg/errors"
)

type Subcmd func(context.Context, *flag.FlagSet, []string) error

type Cmd interface {
	Subcmds() map[string]Subcmd
}

var (
	ErrNoSubcommand  = errors.New("no subcommand")
	ErrBadSubcommand = errors.New("bad subcommand")
)

// Run runs the subcommand of c in args[0], passing the rest of args as its arguments.
func Run(ctx context.Context, c Cmd, args []string) error {
	cmds := c.Subcmds()

	var cmdnames sort.StringSlice
	for cmdname := range cmds {
		cmdnames = append(cmdnames, cmdname)
	}
	cmdnames.Sort()

	if len(args) == 0 {
		return errors.Wrapf(ErrNoSubcommand, "want one of %v", cmdnames)
	}
	subcmd := args[0]
	args = args[1:]
	f, ok := cmds[subcmd]
	if !ok {
		return errors.Wrapf(ErrBadSubcommand, "unknown subcommand %s, want one of %v", subcmd, cmdnames)
	}
	err := f(ctx, flag.NewFlagSet("", flag.ContinueOnError), args)
	return errors.Wrapf(err, "running %s", subcmd)
}
