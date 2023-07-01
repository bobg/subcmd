package subcmd

import (
	"context"
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
		if err := Run(ctx, testPrefixMainCmd{}, []string{"subcmd", "a", "b", "c"}); err != nil {
			t.Error(err)
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

type testPrefixMainCmd struct{}

func (testPrefixMainCmd) Subcmds() Map   { return nil }
func (testPrefixMainCmd) Prefix() string { return "foo-" }
