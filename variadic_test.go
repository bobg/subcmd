package subcmd

import (
	"context"
	"reflect"
	"strconv"
	"testing"
)

func TestVariadic(t *testing.T) {
	args := []string{"a", "-x", "7", "b", "c", "d"}

	for _, variadic := range []bool{false, true} {
		t.Run(strconv.FormatBool(variadic), func(t *testing.T) {
			c := vtestcmd{variadic: variadic}
			if err := Run(context.Background(), &c, args); err != nil {
				t.Fatal(err)
			}

			want := vtestcmd{
				variadic: variadic,
				x:        7,
				pos:      "b",
				rest:     []string{"c", "d"},
			}
			if !reflect.DeepEqual(c, want) {
				t.Errorf("got %+v, want %+v", c, want)
			}
		})
	}
}

type vtestcmd struct {
	// In.
	variadic bool

	// Out.
	x    int
	pos  string
	rest []string
}

func (c *vtestcmd) Subcmds() Map {
	var f interface{}

	if c.variadic {
		f = c.v
	} else {
		f = c.a
	}

	return Commands(
		"a", f, "", Params(
			"-x", Int, 0, "",
			"pos", String, "", "",
		),
	)
}

func (c *vtestcmd) a(ctx context.Context, x int, pos string, rest []string) error {
	c.x = x
	c.pos = pos
	c.rest = rest
	return nil
}

func (c *vtestcmd) v(ctx context.Context, x int, pos string, rest ...string) error {
	c.x = x
	c.pos = pos
	c.rest = rest
	return nil
}
