package subcmd_test

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/bobg/subcmd"
)

func Example() {
	// First, parse global flags normally.
	var (
		verbose = flag.Bool("verbose", false, "be verbose")
		config  = flag.String("config", ".config", "path to config file")
	)
	flag.Parse()

	confData, err := ioutil.ReadFile(*config)
	if err != nil {
		panic(err)
	}

	// Stash the relevant info into an object that implements the subcmd.Cmd interface.
	cmd := command{
		conf:    string(confData),
		verbose: *verbose,
	}

	// Pass that object to subcmd.Run,
	// which will parse a subcommand and its flags
	// from the remaining command-line arguments
	// and run them.
	err = subcmd.Run(context.Background(), cmd, flag.Args())
	if err != nil {
		panic(err)
	}
}

// Type command implements the subcmd.Cmd interface,
// meaning that it can report its subcommands,
// their names,
// and their parameters and types
// via the Subcmds method.
type command struct {
	conf    string
	verbose bool
}

func (c command) Subcmds() subcmd.Map {
	// There are two subcommands:
	// hello, which takes -name and -spanish flags,
	// and add, which takes no flags.
	return subcmd.Commands(
		"hello", hello, subcmd.Params(
			"name", subcmd.String, "", "name to greet",
			"spanish", subcmd.Bool, false, "greet in Spanish",
		),
		"add", c.add, nil,
	)
}

// The implementation of a subcommand takes a context object,
// the values of its parsed flags,
// and a slice of remaining command-line arguments
// (which could be used in another call to subcmd.Run
// to implement a sub-subcommand).
// It can optionally return an error.
func hello(_ context.Context, name string, spanish bool, _ []string) {
	if spanish {
		fmt.Print("Hola")
	} else {
		fmt.Print("Hello")
	}
	if name != "" {
		fmt.Printf(" %s", name)
	}
	fmt.Print("\n")
}

func (c command) add(_ context.Context, args []string) error {
	if c.verbose {
		fmt.Printf("adding %d value(s)\n", len(args))
	}
	var result float64
	for _, arg := range args {
		aval, err := strconv.ParseFloat(arg, 64)
		if err != nil {
			return err
		}
		result += aval
	}
	fmt.Println(result)
	return nil
}
