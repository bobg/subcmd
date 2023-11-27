package subcmd

import (
	"flag"
	"testing"
)

func TestToFlagSet(t *testing.T) {
	cases := []struct {
		name   string
		params []Param
		wantfs map[string]interface{}
	}{{
		name: "float64_with_int_default",
		params: []Param{{
			Name:    "-float64",
			Type:    Float64,
			Default: 1,
		}},
		wantfs: map[string]interface{}{
			"float64": 1.0,
		},
	}}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			fs, _, _, err := ToFlagSet(c.params)
			if err != nil {
				t.Fatal(err)
			}

			for k, v := range c.wantfs {
				val := fs.Lookup(k)
				if val == nil {
					t.Fatalf("flag %s not found", k)
				}
				getter := val.Value.(flag.Getter)
				got := getter.Get()
				if got != v {
					t.Errorf("got %v (%T), want %v (%T)", got, got, v, v)
				}
			}
		})
	}
}
