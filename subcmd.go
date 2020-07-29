package subcmd

import (
	"context"
	"flag"

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
	if len(args) == 0 {
		return ErrNoSubcommand
	}
	subcmd := args[0]
	args = args[1:]
	cmds := c.Subcmds()
	f, ok := cmds[subcmd]
	if !ok {
		return errors.Wrap(ErrBadSubcommand, subcmd)
	}
	err := f(ctx, flag.NewFlagSet("", flag.ContinueOnError), args)
	return errors.Wrapf(err, "running %s", subcmd)
}
