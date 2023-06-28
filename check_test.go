package subcmd

import (
	"context"
	"testing"
)

func TestCheck(t *testing.T) {
	cases := []struct {
		name    string
		f       interface{}
		ptypes  []Type
		wantErr bool
	}{{
		name: "noParamsOK",
		f:    checkNoParams,
	}, {
		name:    "noParamsBad",
		f:       checkNoParams,
		ptypes:  []Type{Bool},
		wantErr: true,
	}, {
		name: "noParamsNoErrOK",
		f:    checkNoParamsNoErr,
	}, {
		name:    "noParamsNoErrBad",
		f:       checkNoParamsNoErr,
		ptypes:  []Type{Bool},
		wantErr: true,
	}, {
		name:   "oneBoolFlagOK",
		f:      checkOneBoolFlag,
		ptypes: []Type{Bool},
	}, {
		name:    "oneBoolFlagTooMany",
		f:       checkOneBoolFlag,
		ptypes:  []Type{Bool, Bool},
		wantErr: true,
	}, {
		name:    "oneBoolFlagTooFew",
		f:       checkOneBoolFlag,
		wantErr: true,
	}, {
		name:    "oneBoolFlagWrongType",
		f:       checkOneBoolFlag,
		ptypes:  []Type{Int},
		wantErr: true,
	}}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var params []Param
			for _, ptype := range tc.ptypes {
				params = append(params, Param{Type: ptype})
			}
			err := Check(Subcmd{F: tc.f, Params: params})
			switch {
			case err == nil && tc.wantErr:
				t.Error("got no error but want one")
			case err != nil && !tc.wantErr:
				t.Error(err)
			}
		})
	}
}

func checkNoParamsNoErr(_ context.Context, _ []string) {}

func checkNoParams(_ context.Context, _ []string) error {
	return nil
}

func checkOneBoolFlag(_ context.Context, _ bool, _ []string) {}
