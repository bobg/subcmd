package subcmd

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestNoArgs(t *testing.T) {
	err := Run(context.Background(), command{}, []string{"y"})
	if err != nil {
		t.Error(err)
	}
}

func TestRun(t *testing.T) {
	var (
		ctx  = context.Background()
		args = []string{"x"}
		cmd  command
	)
	cmd.booltest(ctx, t, args)
}

func (cmd command) booltest(ctx context.Context, t *testing.T, args []string) {
	for _, boolopt := range []bool{false, true} {
		cmd.boolopt = boolopt
		t.Run(fmt.Sprintf("bool_%v", boolopt), func(t *testing.T) {
			if boolopt {
				args = append(args, "-boolopt")
			}
			cmd.inttest(ctx, t, args)
		})
	}
}

func (cmd command) inttest(ctx context.Context, t *testing.T, args []string) {
	for _, intopt := range []int{0, 1} {
		cmd.intopt = intopt
		t.Run(fmt.Sprintf("int_%d", intopt), func(t *testing.T) {
			if intopt != 0 {
				args = append(args, "-intopt", strconv.Itoa(intopt))
			}
			cmd.int64test(ctx, t, args)
		})
	}
}

func (cmd command) int64test(ctx context.Context, t *testing.T, args []string) {
	for _, int64opt := range []int64{0, 2} {
		cmd.int64opt = int64opt
		t.Run(fmt.Sprintf("int64_%d", int64opt), func(t *testing.T) {
			if int64opt != 0 {
				args = append(args, "-int64opt", strconv.FormatInt(int64opt, 10))
			}
			cmd.uinttest(ctx, t, args)
		})
	}
}

func (cmd command) uinttest(ctx context.Context, t *testing.T, args []string) {
	for _, uintopt := range []uint{0, 3} {
		cmd.uintopt = uintopt
		t.Run(fmt.Sprintf("uint_%d", uintopt), func(t *testing.T) {
			if uintopt != 0 {
				args = append(args, "-uintopt", strconv.FormatUint(uint64(uintopt), 10))
			}
			cmd.uint64test(ctx, t, args)
		})
	}
}

func (cmd command) uint64test(ctx context.Context, t *testing.T, args []string) {
	for _, uint64opt := range []uint64{0, 4} {
		cmd.uint64opt = uint64opt
		t.Run(fmt.Sprintf("uint64_%d", uint64opt), func(t *testing.T) {
			if uint64opt != 0 {
				args = append(args, "-uint64opt", strconv.FormatUint(uint64opt, 10))
			}
			cmd.strtest(ctx, t, args)
		})
	}
}

func (cmd command) strtest(ctx context.Context, t *testing.T, args []string) {
	for _, stropt := range []string{"", "foo"} {
		cmd.stropt = stropt
		t.Run(fmt.Sprintf("str%s", stropt), func(t *testing.T) {
			if stropt != "" {
				args = append(args, "-stropt", stropt)
			}
			cmd.float64test(ctx, t, args)
		})
	}
}

func (cmd command) float64test(ctx context.Context, t *testing.T, args []string) {
	for _, float64opt := range []float64{0, 3.14} {
		cmd.float64opt = float64opt
		t.Run(fmt.Sprintf("float%f", float64opt), func(t *testing.T) {
			if float64opt != 0 {
				args = append(args, "-float64opt", strconv.FormatFloat(float64opt, 'f', -1, 64))
			}
			cmd.durtest(ctx, t, args)
		})
	}
}

func (cmd command) durtest(ctx context.Context, t *testing.T, args []string) {
	for _, duropt := range []time.Duration{0, time.Minute} {
		cmd.duropt = duropt
		t.Run(fmt.Sprintf("dur%s", duropt), func(t *testing.T) {
			if duropt != 0 {
				args = append(args, "-duropt", fmt.Sprintf("%s", duropt))
			}
			cmd.t = t
			err := Run(ctx, cmd, args)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

type command struct {
	t          *testing.T
	boolopt    bool
	intopt     int
	int64opt   int64
	uintopt    uint
	uint64opt  uint64
	stropt     string
	float64opt float64
	duropt     time.Duration
}

func (c command) Subcmds() Map {
	return Commands(
		"x", c.xcmd, "x", Params(
			"-boolopt", Bool, false, "bool",
			"-intopt", Int, 0, "int",
			"-int64opt", Int64, 0, "int64",
			"-uintopt", Uint, 0, "uint",
			"-uint64opt", Uint64, 0, "uint64",
			"-stropt", String, "", "str",
			"-float64opt", Float64, 0, "float64",
			"-duropt", Duration, 0, "dur",
		),
		"y", c.ycmd, "y", nil,
	)
}

func (c command) xcmd(
	_ context.Context,
	boolopt bool,
	intopt int,
	int64opt int64,
	uintopt uint,
	uint64opt uint64,
	stropt string,
	float64opt float64,
	duropt time.Duration,
	_ []string,
) error {
	if boolopt != c.boolopt {
		c.t.Errorf("boolopt: got %v, want %v", boolopt, c.boolopt)
	}
	if intopt != c.intopt {
		c.t.Errorf("intopt: got %d, want %d", intopt, c.intopt)
	}
	if int64opt != c.int64opt {
		c.t.Errorf("int64opt: got %d, want %d", int64opt, c.int64opt)
	}
	if uintopt != c.uintopt {
		c.t.Errorf("uintopt: got %d, want %d", uintopt, c.uintopt)
	}
	if uint64opt != c.uint64opt {
		c.t.Errorf("uint64opt: got %d, want %d", uint64opt, c.uint64opt)
	}
	if stropt != c.stropt {
		c.t.Errorf(`stropt: got "%s", want "%s"`, stropt, c.stropt)
	}
	if float64opt != c.float64opt {
		c.t.Errorf("float64opt: got %f, want %f", float64opt, c.float64opt)
	}
	if duropt != c.duropt {
		c.t.Errorf("duropt: got %s, want %s", duropt, c.duropt)
	}
	return nil
}

func TestCommands(t *testing.T) {
	got := Commands(
		"foo", foocmd, "foo command", Params(
			"a", Bool, false, "flag a",
			"b", Int, 0, "flag b",
		),
		"bar", barcmd, "bar command", nil,
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
			Desc: "foo command",
		},
		"bar": Subcmd{
			F:    barcmd,
			Desc: "bar command",
		},
	}
	if diff := cmp.Diff(want, got, fooopt, baropt); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}

func (command) ycmd(_ context.Context, _ []string) error { return nil }

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

func TestUsage(t *testing.T) {
	c := command{}.Subcmds()["x"]

	t.Run("long", func(t *testing.T) {
		const want = `x
-boolopt           bool
-duropt duration   dur
-float64opt float  float64
-int64opt int      int64
-intopt int        int
-stropt string     str
-uint64opt uint    uint64
-uintopt uint      uint
`

		got, err := c.Usage(true)
		if err != nil {
			t.Fatal()
		}
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("short", func(t *testing.T) {
		got, err := c.Usage(false)
		if err != nil {
			t.Fatal()
		}
		const want = "[-boolopt] [-duropt duration] [-float64opt float] [-int64opt int] [-intopt int] [-stropt string] [-uint64opt uint] [-uintopt uint]"
		if got != want {
			t.Errorf(`got "%s", want "%s"`, got, want)
		}
	})
}
