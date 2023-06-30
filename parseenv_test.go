package subcmd

import (
	"encoding/json"
	"os"
	"testing"
)

func TestParseEnv(t *testing.T) {
	j, err := json.Marshal(map[string]interface{}{
		"foo": "bar",
		"baz": 42,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Can't use t.Setenv, introduced in Go 1.17,
	// because we're on Go 1.14 and don't want to break callers unnecessarily.
	oldval, ok := os.LookupEnv(EnvVar)
	if ok {
		defer os.Setenv(EnvVar, oldval)
	} else {
		defer os.Unsetenv(EnvVar)
	}
	os.Setenv(EnvVar, string(j))

	type s struct {
		Foo string `json:"foo"`
		Baz int    `json:"baz"`
	}

	var (
		got  s
		want = s{
			Foo: "bar",
			Baz: 42,
		}
	)

	if err := ParseEnv(&got); err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Errorf("got %#v, want %#v", got, want)
	}
}
