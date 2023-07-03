package subcmd

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestPrefix(t *testing.T) {
	ctx := context.Background()

	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	path := os.Getenv("PATH")
	path += ":" + filepath.Join(wd, "testdata")

	restoreEnv := testSetenv("PATH", path)
	defer restoreEnv()

	t.Run("subcmd", func(t *testing.T) {
		oldStdout := os.Stdout

		f, err := os.CreateTemp("", "subcmd")
		if err != nil {
			t.Fatal(err)
		}
		tmpname := f.Name()
		defer os.Remove(tmpname)
		defer f.Close()

		os.Stdout = f
		defer func() { os.Stdout = oldStdout }()

		c := testPrefixMainCmd{Data: "xyz"}

		if err := Run(ctx, c, []string{"subcmd", "a", "b", "c"}); err != nil {
			t.Error(err)
		}
		if err := f.Close(); err != nil {
			t.Fatal(err)
		}

		f, err = os.Open(tmpname)
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()

		var got testPrefixMainCmd
		if err = json.NewDecoder(f).Decode(&got); err != nil {
			t.Fatal(err)
		}

		if got != c {
			t.Errorf("got %+v, want %+v", got, c)
		}
	})

	t.Run("nosubcmd", func(t *testing.T) {
		err := Run(ctx, testPrefixMainCmd{}, []string{"nosubcmd", "a", "b", "c"})
		var u *UnknownSubcmdErr
		switch {
		case errors.As(err, &u):
			// ok

		case err != nil:
			t.Error(err)

		default:
			t.Errorf("expected error of type %T, got nil", u)
		}
	})
}

type testPrefixMainCmd struct {
	Data string `json:"data"`
}

func (testPrefixMainCmd) Subcmds() Map   { return nil }
func (testPrefixMainCmd) Prefix() string { return "foo-" }
