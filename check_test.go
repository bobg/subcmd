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
				if err := Check(Subcmd{F: fOK, Params: []Param{{Type: ptyp}}}); err != nil {
					t.Error(err)
				}
			})

			t.Run("toomany", func(t *testing.T) {
				err := Check(Subcmd{F: fTooMany, Params: []Param{{Type: ptyp}}})
				var npErr NumParamsErr
				if errors.As(err, &npErr) {
					if npErr.Got != 4 || npErr.Want != 3 {
						t.Errorf("got got=%d, want=%d, want got=4, want=3", npErr.Got, npErr.Want)
					}
				} else {
					t.Errorf("got %v, want NumParamsErr", err)
				}
			})

			t.Run("toofew", func(t *testing.T) {
				err := Check(Subcmd{F: func(context.Context, []string) {}, Params: []Param{{Type: ptyp}}})
				var npErr NumParamsErr
				if errors.As(err, &npErr) {
					if npErr.Got != 2 || npErr.Want != 3 {
						t.Errorf("got got=%d, want=%d, want got=2, want=3", npErr.Got, npErr.Want)
					}
				} else {
					t.Errorf("got %v, want NumParamsErr", err)
				}
			})

			t.Run("wrongtype", func(t *testing.T) {
				for ptyp2 := Bool; ptyp2 <= Duration; ptyp2++ {
					if ptyp2 == ptyp {
						continue
					}
					t.Run(ptyp2.String(), func(t *testing.T) {
						err := Check(Subcmd{F: fOK, Params: []Param{{Type: ptyp2}}})
						var ptErr ParamTypeErr
						if errors.As(err, &ptErr) {
							if ptErr.N != 1 {
								t.Errorf("got N=%d, want N=1", ptErr.N)
							}
							if ptErr.Want != ptyp2 {
								t.Errorf("got Want=%v, want Want=%v", ptErr.Want, ptyp2)
							}
							if rt := ptyp.reflectType(); ptErr.Got != rt {
								t.Errorf("got Got=%v, want Got=%v", ptErr.Got, rt)
							}
						} else {
							t.Errorf("got %v, want ParamTypeErr", err)
						}
					})
				}
			})
		})
	}
}

func TestCheckNotFunc(t *testing.T) {
	if err := Check(Subcmd{F: 42}); !errors.Is(err, ErrNotAFunction) {
		t.Errorf("got %v, want ErrNotAFunction", err)
	}
}

func TestCheckNoContext(t *testing.T) {
	if err := Check(Subcmd{F: func(int, []string) {}}); !errors.Is(err, ErrNoContext) {
		t.Errorf("got %v, want ErrNoContext", err)
	}
}

func TestCheckNoStringSlice(t *testing.T) {
	if err := Check(Subcmd{F: func(context.Context, int) {}}); !errors.Is(err, ErrNoStringSlice) {
		t.Errorf("got %v, want ErrNoStringSlice", err)
	}
}

func TestCheckNoError(t *testing.T) {
	if err := Check(Subcmd{F: func(context.Context, []string) int { return 0 }}); !errors.Is(err, ErrNotError) {
		t.Errorf("got %v, want ErrNotError", err)
	}
}

func TestTooManyReturns(t *testing.T) {
	if err := Check(Subcmd{F: func(context.Context, []string) (int, int) { return 0, 0 }}); !errors.Is(err, ErrTooManyReturns) {
		t.Errorf("got %v, want ErrTooManyReturns", err)
	}
}
