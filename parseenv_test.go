package subcmd

import (
	"encoding/json"
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
	t.Setenv(EnvVar, string(j))

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
