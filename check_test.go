package subcmd

import (
	"context"
	"testing"
	"time"
)

func TestCheckZeroArgs(t *testing.T) {
	cases := []struct {
		name    string
		f       interface{}
		wantErr bool
	}{{
		name: "noErr",
		f:    func(context.Context, []string) {},
	}, {
		name: "err",
		f:    func(context.Context, []string) error { return nil },
	}, {
		name:    "noContext",
		f:       func([]string) {},
		wantErr: true,
	}, {
		name:    "noArgs",
		f:       func(context.Context) {},
		wantErr: true,
	}, {
		name:    "extraArg",
		f:       func(context.Context, int, []string) {},
		wantErr: true,
	}}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := Check(Subcmd{F: tc.f})
			if gotErr := err != nil; gotErr != tc.wantErr {
				t.Errorf("got err %v, wantErr is %v", err, tc.wantErr)
			}
		})
	}
}

func TestCheckOneArg(t *testing.T) {
	for ptyp := Bool; ptyp <= Duration; ptyp++ {
		t.Run(ptyp.String(), func(t *testing.T) {
			var fOK, fTooMany interface{}
			switch ptyp {
			case Bool:
				fOK = func(context.Context, bool, []string) {}
				fTooMany = func(context.Context, bool, bool, []string) {}
			case Int:
				fOK = func(context.Context, int, []string) {}
				fTooMany = func(context.Context, int, int, []string) {}
			case Int64:
				fOK = func(context.Context, int64, []string) {}
				fTooMany = func(context.Context, int64, int64, []string) {}
			case Uint:
				fOK = func(context.Context, uint, []string) {}
				fTooMany = func(context.Context, uint, uint, []string) {}
			case Uint64:
				fOK = func(context.Context, uint64, []string) {}
				fTooMany = func(context.Context, uint64, uint64, []string) {}
			case String:
				fOK = func(context.Context, string, []string) {}
				fTooMany = func(context.Context, string, string, []string) {}
			case Float64:
				fOK = func(context.Context, float64, []string) {}
				fTooMany = func(context.Context, float64, float64, []string) {}
			case Duration:
				fOK = func(context.Context, time.Duration, []string) {}
				fTooMany = func(context.Context, time.Duration, time.Duration, []string) {}
			}

			if err := Check(Subcmd{F: fOK, Params: []Param{{Type: ptyp}}}); err != nil {
				t.Error(err)
			}

			_ = fTooMany

		})

	}
}
