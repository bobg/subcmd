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
	err := Run(context.Background(), &command{}, []string{"y"})
	if err != nil {
		t.Error(err)
	}
}

func TestOptionalArg(t *testing.T) {
	cmd := new(command)

	err := Run(context.Background(), cmd, []string{"z"})
	if err != nil {
		t.Fatal(err)
	}
	if cmd.zgotopt != 7 {
		t.Errorf("got %d, want 7", cmd.zgotopt)
	}

	err = Run(context.Background(), cmd, []string{"z", "17"})
	if err != nil {
		t.Fatal(err)
	}
	if cmd.zgotopt != 17 {
		t.Errorf("got %d, want 17", cmd.zgotopt)
	}
}

func TestRun(t *testing.T) {
	var (
		ctx  = context.Background()
		args = []string{"x"}
		cmd  command
	)
	cmd.boolflagtest(ctx, t, args)
}

func (cmd *command) boolflagtest(ctx context.Context, t *testing.T, args []string) {
	for _, boolopt := range []bool{false, true} {
		cmd.boolopt = boolopt
		t.Run(fmt.Sprintf("bool_%v", boolopt), func(t *testing.T) {
			if boolopt {
				args = append(args, "-boolopt")
			}
			cmd.intflagtest(ctx, t, args)
		})
	}
}

func (cmd *command) intflagtest(ctx context.Context, t *testing.T, args []string) {
	for _, intopt := range []int{0, 1} {
		cmd.intopt = intopt
		t.Run(fmt.Sprintf("int_%d", intopt), func(t *testing.T) {
			if intopt != 0 {
				args = append(args, "-intopt", strconv.Itoa(intopt))
			}
			cmd.int64flagtest(ctx, t, args)
		})
	}
}

func (cmd *command) int64flagtest(ctx context.Context, t *testing.T, args []string) {
	for _, int64opt := range []int64{0, 2} {
		cmd.int64opt = int64opt
		t.Run(fmt.Sprintf("int64_%d", int64opt), func(t *testing.T) {
			if int64opt != 0 {
				args = append(args, "-int64opt", strconv.FormatInt(int64opt, 10))
			}
			cmd.uintflagtest(ctx, t, args)
		})
	}
}

func (cmd *command) uintflagtest(ctx context.Context, t *testing.T, args []string) {
	for _, uintopt := range []uint{0, 3} {
		cmd.uintopt = uintopt
		t.Run(fmt.Sprintf("uint_%d", uintopt), func(t *testing.T) {
			if uintopt != 0 {
				args = append(args, "-uintopt", strconv.FormatUint(uint64(uintopt), 10))
			}
			cmd.uint64flagtest(ctx, t, args)
		})
	}
}

func (cmd *command) uint64flagtest(ctx context.Context, t *testing.T, args []string) {
	for _, uint64opt := range []uint64{0, 4} {
		cmd.uint64opt = uint64opt
		t.Run(fmt.Sprintf("uint64_%d", uint64opt), func(t *testing.T) {
			if uint64opt != 0 {
				args = append(args, "-uint64opt", strconv.FormatUint(uint64opt, 10))
			}
			cmd.strflagtest(ctx, t, args)
		})
	}
}

func (cmd *command) strflagtest(ctx context.Context, t *testing.T, args []string) {
	for _, stropt := range []string{"", "foo"} {
		cmd.stropt = stropt
		t.Run(fmt.Sprintf("str%s", stropt), func(t *testing.T) {
			if stropt != "" {
				args = append(args, "-stropt", stropt)
			}
			cmd.float64flagtest(ctx, t, args)
		})
	}
}

func (cmd *command) float64flagtest(ctx context.Context, t *testing.T, args []string) {
	for _, float64opt := range []float64{0, 3.14} {
		cmd.float64opt = float64opt
		t.Run(fmt.Sprintf("float%f", float64opt), func(t *testing.T) {
			if float64opt != 0 {
				args = append(args, "-float64opt", strconv.FormatFloat(float64opt, 'f', -1, 64))
			}
			cmd.durflagtest(ctx, t, args)
		})
	}
}

func (cmd *command) durflagtest(ctx context.Context, t *testing.T, args []string) {
	for _, duropt := range []time.Duration{0, time.Minute} {
		cmd.duropt = duropt
		t.Run(fmt.Sprintf("dur%s", duropt), func(t *testing.T) {
			if duropt != 0 {
				args = append(args, "-duropt", duropt.String())
			}
			cmd.boolpostest(ctx, t, args)
		})
	}
}

func (cmd *command) boolpostest(ctx context.Context, t *testing.T, args []string) {
	cmd.boolpos = true
	t.Run(fmt.Sprintf("boolpos_%v", cmd.boolpos), func(t *testing.T) {
		args = append(args, "true")
		cmd.intpostest(ctx, t, args)
	})
}

func (cmd *command) intpostest(ctx context.Context, t *testing.T, args []string) {
	cmd.intpos = 412
	t.Run(fmt.Sprintf("intpos_%d", cmd.intpos), func(t *testing.T) {
		args = append(args, "412")
		cmd.int64postest(ctx, t, args)
	})
}

func (cmd *command) int64postest(ctx context.Context, t *testing.T, args []string) {
	cmd.int64pos = 733
	t.Run(fmt.Sprintf("int64pos_%d", cmd.int64pos), func(t *testing.T) {
		args = append(args, "733")
		cmd.uintpostest(ctx, t, args)
	})
}

func (cmd *command) uintpostest(ctx context.Context, t *testing.T, args []string) {
	cmd.uintpos = 31178
	t.Run(fmt.Sprintf("uintpos_%d", cmd.uintpos), func(t *testing.T) {
		args = append(args, "31178")
		cmd.uint64postest(ctx, t, args)
	})
}

func (cmd *command) uint64postest(ctx context.Context, t *testing.T, args []string) {
	cmd.uint64pos = 2134
	t.Run(fmt.Sprintf("uint64pos_%d", cmd.uint64pos), func(t *testing.T) {
		args = append(args, "2134")
		cmd.strpostest(ctx, t, args)
	})
}

func (cmd *command) strpostest(ctx context.Context, t *testing.T, args []string) {
	cmd.strpos = "plugh"
	t.Run(fmt.Sprintf("strpos%s", cmd.strpos), func(t *testing.T) {
		args = append(args, "plugh")
		cmd.float64postest(ctx, t, args)
	})
}

func (cmd *command) float64postest(ctx context.Context, t *testing.T, args []string) {
	cmd.float64pos = 2.718
	t.Run(fmt.Sprintf("floatpos%f", cmd.float64pos), func(t *testing.T) {
		args = append(args, "2.718")
		cmd.durpostest(ctx, t, args)
	})
}

func (cmd *command) durpostest(ctx context.Context, t *testing.T, args []string) {
	cmd.durpos = 7 * time.Second
	t.Run(fmt.Sprintf("durpos%s", cmd.durpos), func(t *testing.T) {
		args = append(args, "7s")

		cmd.t = t
		err := Run(ctx, cmd, args)
		if err != nil {
			t.Fatal(err)
		}
	})
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
	boolpos    bool
	intpos     int
	int64pos   int64
	uintpos    uint
	uint64pos  uint64
	strpos     string
	float64pos float64
	durpos     time.Duration

	zgotopt int
}

func (cmd *command) Subcmds() Map {
	return Commands(
		"x", cmd.xcmd, "x", Params(
			"-boolopt", Bool, false, "bool flag",
			"-intopt", Int, 0, "int flag",
			"-int64opt", Int64, 0, "int64 flag",
			"-uintopt", Uint, 0, "uint flag",
			"-uint64opt", Uint64, 0, "uint64 flag",
			"-stropt", String, "", "str flag",
			"-float64opt", Float64, 0, "float64 flag",
			"-duropt", Duration, 0, "dur flag",
			"boolpos", Bool, false, "bool pos",
			"intpos", Int, 0, "int pos",
			"int64pos", Int64, 0, "int64 pos",
			"uintpos", Uint, 0, "uint pos",
			"uint64pos", Uint64, 0, "uint64 pos",
			"strpos", String, "", "str pos",
			"float64pos", Float64, 0, "float64 pos",
			"durpos", Duration, 0, "dur pos",
		),
		"y", cmd.ycmd, "y", nil,
		"z", cmd.zcmd, "z", Params(
			"opt?", Int, 7, "optional int",
		),
	)
}

func (cmd *command) xcmd(
	_ context.Context,
	boolopt bool,
	intopt int,
	int64opt int64,
	uintopt uint,
	uint64opt uint64,
	stropt string,
	float64opt float64,
	duropt time.Duration,
	boolpos bool,
	intpos int,
	int64pos int64,
	uintpos uint,
	uint64pos uint64,
	strpos string,
	float64pos float64,
	durpos time.Duration,
	_ []string,
) error {
	if boolopt != cmd.boolopt {
		cmd.t.Errorf("boolopt: got %v, want %v", boolopt, cmd.boolopt)
	}
	if intopt != cmd.intopt {
		cmd.t.Errorf("intopt: got %d, want %d", intopt, cmd.intopt)
	}
	if int64opt != cmd.int64opt {
		cmd.t.Errorf("int64opt: got %d, want %d", int64opt, cmd.int64opt)
	}
	if uintopt != cmd.uintopt {
		cmd.t.Errorf("uintopt: got %d, want %d", uintopt, cmd.uintopt)
	}
	if uint64opt != cmd.uint64opt {
		cmd.t.Errorf("uint64opt: got %d, want %d", uint64opt, cmd.uint64opt)
	}
	if stropt != cmd.stropt {
		cmd.t.Errorf(`stropt: got "%s", want "%s"`, stropt, cmd.stropt)
	}
	if float64opt != cmd.float64opt {
		cmd.t.Errorf("float64opt: got %f, want %f", float64opt, cmd.float64opt)
	}
	if duropt != cmd.duropt {
		cmd.t.Errorf("duropt: got %s, want %s", duropt, cmd.duropt)
	}
	if boolpos != cmd.boolpos {
		cmd.t.Errorf("boolpos: got %v, want %v", boolpos, cmd.boolpos)
	}
	if intpos != cmd.intpos {
		cmd.t.Errorf("intpos: got %d, want %d", intpos, cmd.intpos)
	}
	if int64pos != cmd.int64pos {
		cmd.t.Errorf("int64pos: got %d, want %d", int64pos, cmd.int64pos)
	}
	if uintpos != cmd.uintpos {
		cmd.t.Errorf("uintpos: got %d, want %d", uintpos, cmd.uintpos)
	}
	if uint64pos != cmd.uint64pos {
		cmd.t.Errorf("uint64pos: got %d, want %d", uint64pos, cmd.uint64pos)
	}
	if strpos != cmd.strpos {
		cmd.t.Errorf(`strpos: got "%s", want "%s"`, strpos, cmd.strpos)
	}
	if float64pos != cmd.float64pos {
		cmd.t.Errorf("float64pos: got %f, want %f", float64pos, cmd.float64pos)
	}
	if durpos != cmd.durpos {
		cmd.t.Errorf("durpos: got %s, want %s", durpos, cmd.durpos)
	}
	return nil
}

func TestCommands(t *testing.T) {
	baz := Subcmd{
		F:    bazcmd,
		Desc: "baz command",
	}

	got := Commands(
		"foo", foocmd, "foo command", Params(
			"a", Bool, false, "flag a",
			"b", Int, 0, "flag b",
		),
		"bar", barcmd, "bar command", nil,
		"baz", baz,
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
		"baz": Subcmd{
			F:    bazcmd,
			Desc: "baz command",
		},
	}
	if diff := cmp.Diff(want, got, fooopt, baropt); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}

func (*command) ycmd(_ context.Context, _ []string) error { return nil }

func (cmd *command) zcmd(_ context.Context, opt int, _ []string) error {
	cmd.zgotopt = opt
	return nil
}

// The following is needed because cmp.Diff
// (and reflect.DeepEqual for that matter)
// do not work as expected with function pointers.

func foocmd(context.Context, bool, int, []string) {}
func barcmd(context.Context, []string)            {}
func bazcmd(context.Context, []string)            {}

var (
	foocomparer = func(_, _ func(context.Context, bool, int, []string)) bool { return true }
	barcomparer = func(_, _ func(context.Context, []string)) bool { return true }

	fooopt = cmp.FilterValues(foocomparer, cmp.Comparer(foocomparer))
	baropt = cmp.FilterValues(barcomparer, cmp.Comparer(barcomparer))
)
