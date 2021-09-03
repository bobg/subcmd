package subcmd_test

import (
	"context"
	"database/sql"
	"flag"

	"github.com/bobg/subcmd/v2"
)

func Example() {
	// First, parse global flags normally.
	var (
		driver = flag.String("driver", "sqlite3", "database driver")
		dbname = flag.String("db", "", "database name")
	)
	flag.Parse()

	db, err := sql.Open(*driver, *dbname)
	if err != nil {
		panic(err)
	}

	// Stash the relevant info into an object that implements the subcmd.Cmd interface.
	cmd := command{db: db}

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
	db *sql.DB
}

func (c command) Subcmds() subcmd.Map {
	return subcmd.Commands(
		// The "list" subcommand takes one flag, -reverse.
		"list", c.list, "list employees", subcmd.Params(
			"-reverse", subcmd.Bool, false, "reverse order of list",
		),

		// The "add" subcommand takes no flags but one positional argument.
		"add", c.add, "add new employee", subcmd.Params(
			"name", subcmd.String, "", "employee name",
		),
	)
}

// The implementation of a subcommand takes a context object,
// the values of its parsed flags and positional arguments,
// and a slice of remaining command-line arguments
// (which could be used in another call to subcmd.Run
// to implement a sub-subcommand).
// It can optionally return an error.
func (c command) list(ctx context.Context, reverse bool, _ []string) error {
	query := "SELECT name FROM employees ORDER BY name"
	if reverse {
		query += " DESC"
	}
	rows, err := c.db.QueryContext(ctx, query)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		// ...print the employee name contained in this row...
	}
	return rows.Err()
}

// Implementation of the "add" subcommand.
func (c command) add(ctx context.Context, name string, _ []string) error {
	_, err := c.db.ExecContext(ctx, "INSERT INTO employees (name) VALUES ($1)", name)
	return err
}
