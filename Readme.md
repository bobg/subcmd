# Subcmd - command-line interfaces with subcommands and flags

[![Go Reference](https://pkg.go.dev/badge/github.com/bobg/subcmd.svg)](https://pkg.go.dev/github.com/bobg/subcmd)
[![Go Report Card](https://goreportcard.com/badge/github.com/bobg/subcmd)](https://goreportcard.com/report/github.com/bobg/subcmd)
![Tests](https://github.com/bobg/subcmd/actions/workflows/go.yml/badge.svg)
[![Coverage Status](https://coveralls.io/repos/github/bobg/subcmd/badge.svg?branch=master)](https://coveralls.io/github/bobg/subcmd?branch=master)

This is subcmd,
a Go package for writing command-line programs that require flag parsing
and that have “subcommands” that also require flag parsing.

Use it when you want your program to parse command lines that look like this:

```
command -globalopt subcommand -subopt1 FOO -subopt2 ARG1 ARG2
```

Subcommands may have sub-subcommands and so on.

This is a layer on top of the standard Go `flag` package.

## Usage

```go
func main() {
  // Parse global flags normally.
  dbname := flag.String("db", "", "database connection string")
  flag.Parse()

  db, err := sql.Open(dbdriver, *dbname)
  if err != nil { ... }

  // Stash global options in a top-level command object.
  c := command{db: db}

  // Run the subcommand given in the remainder of the command line.
  err = subcmd.Run(context.Background(), c, flag.Args())
  if err != nil { ... }
}

// The top-level command object.
type command struct {
  db *sql.DB
}

// To be used in subcmd.Run above, `command` must implement this method.
func (c command) Subcmds() subcmd.Map {
  return subcmd.Commands(
    // The "list" subcommand takes one flag, -reverse.
    "list", c.list, subcmd.Params(
      "reverse", subcmd.Bool, false, "reverse order of list",
    ),

    // The "add" subcommand takes no flags.
    "add", c.add, nil,
  )
}

// Implementation of the "list" subcommand.
// The value of the -reverse flag is passed as an argument.
func (c command) list(ctx context.Context, reverse bool, _ []string) error {
  query := "SELECT name FROM employees ORDER BY name"
  if reverse {
    query += " DESC"
  }
  rows, err := c.db.QueryContext(ctx, query)
  if err != nil { ... }
  defer rows.Close()
  for rows.Next() { ... }
  return rows.Err()
}

// Implementation of the "add" subcommand.
func (c command) add(ctx context.Context, args []string) error {
  if len(args) != 1 { ...usage error... }
  _, err := c.db.ExecContext(ctx, "INSERT INTO employees (name) VALUES ($1)", args[0])
  return err
}
```
