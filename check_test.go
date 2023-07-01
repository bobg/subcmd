package subcmd

import (
	"context"
	"errors"
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

			t.Run("one", func(t *testing.T) {
				if err := Check(Subcmd{F: fOK, Params: []Param{{Type: ptyp, Default: dflts[ptyp]}}}); err != nil {
					t.Error(err)
				}
			})

			t.Run("toomany", func(t *testing.T) {
				err := Check(Subcmd{F: fTooMany, Params: []Param{{Type: ptyp, Default: dflts[ptyp]}}})
				var e FuncTypeErr
				if !errors.As(err, &e) {
					t.Errorf("got %v, want FuncTypeErr", err)
				}
			})

			t.Run("toofew", func(t *testing.T) {
				err := Check(Subcmd{F: func(context.Context, []string) {}, Params: []Param{{Type: ptyp, Default: dflts[ptyp]}}})
				var e FuncTypeErr
				if !errors.As(err, &e) {
					t.Errorf("got %v, want FuncTypeErr", err)
				}
			})

			t.Run("wrongtype", func(t *testing.T) {
				for ptyp2 := Bool; ptyp2 <= Duration; ptyp2++ {
					if ptyp2 == ptyp {
						continue
					}

					t.Run(ptyp2.String(), func(t *testing.T) {
						err := Check(Subcmd{F: fOK, Params: []Param{{Type: ptyp2, Default: dflts[ptyp2]}}})
						var e FuncTypeErr
						if !errors.As(err, &e) {
							t.Errorf("got %v, want FuncTypeErr", err)
						}
					})
				}
			})
		})
	}
}

func TestCheckNotFunc(t *testing.T) {
	var e FuncTypeErr
	if err := Check(Subcmd{F: 42}); !errors.As(err, &e) {
		t.Errorf("got %v, want FuncTypeErr", err)
	}
}

func TestCheckNoContext(t *testing.T) {
	var e FuncTypeErr
	if err := Check(Subcmd{F: func(int, []string) {}}); !errors.As(err, &e) {
		t.Errorf("got %v, want ErrNoContext", err)
	}
}

func TestCheckNoStringSlice(t *testing.T) {
	var e FuncTypeErr
	if err := Check(Subcmd{F: func(context.Context, int) {}}); !errors.As(err, &e) {
		t.Errorf("got %v, want ErrNoStringSlice", err)
	}
}

func TestCheckNoError(t *testing.T) {
	var e FuncTypeErr
	if err := Check(Subcmd{F: func(context.Context, []string) int { return 0 }}); !errors.As(err, &e) {
		t.Errorf("got %v, want ErrNotError", err)
	}
}

func TestTooManyReturns(t *testing.T) {
	var e FuncTypeErr
	if err := Check(Subcmd{F: func(context.Context, []string) (int, int) { return 0, 0 }}); !errors.As(err, &e) {
		t.Errorf("got %v, want ErrTooManyReturns", err)
	}
}

func TestCheckParam(t *testing.T) {
	for ptyp := Bool; ptyp <= Duration; ptyp++ {
		t.Run(ptyp.String(), func(t *testing.T) {
			for ptyp2 := Bool; ptyp2 <= Duration; ptyp2++ {
				t.Run(ptyp2.String(), func(t *testing.T) {
					err := checkParam(Param{Type: ptyp, Default: dflts[ptyp2]})
					if ptyp == ptyp2 {
						if err != nil {
							t.Error(err)
						}
					} else {
						var e ParamDefaultErr
						if !errors.As(err, &e) {
							t.Errorf("got %v, want ParamDefaultErr", err)
						}
					}
				})
			}
		})
	}
}

var dflts = map[Type]interface{}{
	Bool:     false,
	Int:      0,
	Int64:    int64(0),
	Uint:     uint(0),
	Uint64:   uint64(0),
	String:   "",
	Float64:  float64(0),
	Duration: time.Duration(0),
}
