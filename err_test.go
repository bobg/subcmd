package subcmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestMissingSubcmdErr(t *testing.T) {
	err := Run(context.Background(), errtestcmd{}, nil)
	if err == nil {
		t.Fatal("got no error, want MissingSubcmdErr")
	}

	var merr *MissingSubcmdErr
	if !errors.As(err, &merr) {
		t.Fatalf("got %T, want *MissingSubcmdErr", err)
	}

	got := merr.Error()
	want := `missing subcommand, want one of: a; bb; ccc`
	if got != want {
		t.Errorf(`got "%s", want "%s"`, got, want)
	}

	got = merr.Detail()
	want = `Missing subcommand, want one of:
a    Do a
bb   Do b
ccc  Do c
`

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}

func TestUnknownSubcmdErr(t *testing.T) {
	err := Run(context.Background(), errtestcmd{}, []string{"dddd"})
	if err == nil {
		t.Fatal("got no error, want UnknownSubcmdErr")
	}

	var uerr *UnknownSubcmdErr
	if !errors.As(err, &uerr) {
		t.Fatalf("got %T, want *UnknownSubcmdErr", err)
	}

	got := uerr.Error()
	want := `unknown subcommand "dddd", want one of: a; bb; ccc`
	if got != want {
		t.Errorf(`got "%s", want "%s"`, got, want)
	}

	got = uerr.Detail()
	want = `Unknown subcommand "dddd", want one of:
a    Do a
bb   Do b
ccc  Do c
`

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}

func TestHelpRequestedErr(t *testing.T) {
	t.Run("no subcmd", func(t *testing.T) {
		err := Run(context.Background(), errtestcmd{}, []string{"help"})
		if err == nil {
			t.Fatal("got no error, want HelpRequestedErr")
		}

		var herr *HelpRequestedErr
		if !errors.As(err, &herr) {
			t.Fatalf("got %T, want *HelpRequestedErr", err)
		}

		t.Run("short", func(t *testing.T) {
			got := herr.Error()
			want := "subcommands are: a; bb; ccc"
			if got != want {
				t.Errorf(`got "%s", want "%s"`, got, want)
			}
		})

		t.Run("long", func(t *testing.T) {
			got := herr.Detail()
			want := `Subcommands are:
a    Do a
bb   Do b
ccc  Do c
`

			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	})

	t.Run("with subcmd", func(t *testing.T) {
		err := Run(context.Background(), errtestcmd{}, []string{"help", "a"})
		if err == nil {
			t.Fatal("got no error, want HelpRequestedErr")
		}

		var herr *HelpRequestedErr
		if !errors.As(err, &herr) {
			t.Fatalf("got %T, want *HelpRequestedErr", err)
		}

		t.Run("short", func(t *testing.T) {
			got := herr.Error()
			want := fmt.Sprintf(`usage: %s a [-a1] [-a2 int] [-a3 word] a4 [a5]`, os.Args[0])
			if got != want {
				t.Errorf(`got "%s", want "%s"`, got, want)
			}
		})

		t.Run("long", func(t *testing.T) {
			got := herr.Detail()
			want := fmt.Sprintf(`a: Do a
Usage: %s a [-a1] [-a2 int] [-a3 word] a4 [a5]
-a1       the a1 flag
-a2 int   the a2 flag
-a3 word  a word flag
`, os.Args[0])

			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	})
}

func TestTooFewArgs(t *testing.T) {
	err := Run(context.Background(), errtestcmd{}, []string{"a"})
	if !errors.Is(err, ErrTooFewArgs) {
		t.Errorf("got %v, want %s", err, ErrTooFewArgs)
	}
}

func TestParseErr(t *testing.T) {
	err := Run(context.Background(), errtestcmd{}, []string{"a", "x"})
	var perr ParseErr
	if !errors.As(err, &perr) {
		t.Errorf("got %T, want *ParseErr", err)
	}
}

type errtestcmd struct{}

func (errtestcmd) Subcmds() Map {
	return Commands(
		"a", errtestA, "Do a", Params(
			"-a1", Bool, false, "the a1 flag",
			"-a2", Int, 0, "the a2 flag",
			"-a3", String, "", "a `word` flag",
			"a4", Duration, 0, "positional duration",
			"a5?", Bool, false, "optional positional bool",
		),
		"bb", errtestB, "Do b", nil,
		"ccc", errtestC, "Do c", nil,
	)
}

func errtestA(_ context.Context, _ []string) error { return nil }
func errtestB(_ context.Context, _ []string) error { return nil }
func errtestC(_ context.Context, _ []string) error { return nil }
