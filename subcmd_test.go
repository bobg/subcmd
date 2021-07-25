package subcmd

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRun(t *testing.T) {
	var (
		whichsubcommand string
		whichyval       string
	)
	cmd := &command{
		whichsubcommand: &whichsubcommand,
		whichyval:       &whichyval,
	}
	err := Run(context.Background(), cmd, []string{"x", "-y", "z"})
	if err != nil {
		t.Fatal(err)
	}
	if whichsubcommand != "x" {
		t.Errorf(`got subcommand "%s", want "x"`, whichsubcommand)
	}
	if whichyval != "z" {
		t.Errorf(`got y value "%s", want "z"`, whichyval)
	}
}

type command struct {
	whichsubcommand, whichyval *string
}

func (c *command) Subcmds() Map {
	return Commands(
		"x", c.xcmd, Params(
			"y", String, "", "y value",
		),
	)
}

func (c *command) xcmd(_ context.Context, y string, _ []string) {
	*c.whichsubcommand = "x"
	*c.whichyval = y
}

func TestCommands(t *testing.T) {
	got := Commands(
		"foo", foocmd, Params(
			"a", Bool, false, "flag a",
			"b", Int, 0, "flag b",
		),
		"bar", barcmd, nil,
	)
	want := Map{
		"foo": Subcmd{
			F: foocmd,
			Params: []Param{{
				Name:    "a",
				Type:    Bool,
				Default: false,
				Doc:     "flag a",
			}, {
				Name:    "b",
				Type:    Int,
				Default: 0,
				Doc:     "flag b",
			}},
		},
		"bar": Subcmd{
			F: barcmd,
		},
	}
	if diff := cmp.Diff(want, got, fooopt, baropt); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}

// The following is needed because cmp.Diff
// (and reflect.DeepEqual for that matter)
// do not work as expected with function pointers.

func foocmd(context.Context, bool, int, []string) {}
func barcmd(context.Context, []string)            {}

var (
	foocomparer = func(_, _ func(context.Context, bool, int, []string)) bool { return true }
	barcomparer = func(_, _ func(context.Context, []string)) bool { return true }

	fooopt = cmp.FilterValues(foocomparer, cmp.Comparer(foocomparer))
	baropt = cmp.FilterValues(barcomparer, cmp.Comparer(barcomparer))
)
