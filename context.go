package subcmd

import (
	"context"
	"flag"
)

type ctxkey int

const (
	fsKey ctxkey = iota + 1
	subcmdPairListKey
	suppressCheckKey
)

func withFlagSet(ctx context.Context, fs *flag.FlagSet) context.Context {
	return context.WithValue(ctx, fsKey, fs)
}

// FlagSet produces the [flag.FlagSet] used in a call to a [Subcmd] function.
func FlagSet(ctx context.Context) *flag.FlagSet {
	val, _ := ctx.Value(fsKey).(*flag.FlagSet)
	return val
}

type subcmdPair struct {
	name   string
	subcmd Subcmd
}

func subcmdPairList(ctx context.Context) []subcmdPair {
	pairListPtr := ctx.Value(subcmdPairListKey)
	if pairListPtr == nil {
		return nil
	}
	return *(pairListPtr.(*[]subcmdPair))
}

func addSubcmdPair(ctx context.Context, name string, subcmd Subcmd) context.Context {
	var pairListPtr *[]subcmdPair
	if pairListPtrVal := ctx.Value(subcmdPairListKey); pairListPtrVal == nil {
		var pairList []subcmdPair
		pairListPtr = &pairList
		ctx = context.WithValue(ctx, subcmdPairListKey, pairListPtr)
	} else {
		pairListPtr = pairListPtrVal.(*[]subcmdPair)
	}
	*pairListPtr = append(*pairListPtr, subcmdPair{name: name, subcmd: subcmd})
	return ctx
}

// WithSuppressCheck creates a new context that tells whether to suppress the [CheckMap] call in [Run].
// This value can be retrieved with [SuppressCheck].
func WithSuppressCheck(ctx context.Context, suppress bool) context.Context {
	return context.WithValue(ctx, suppressCheckKey, suppress)
}

// SuppressCheck tells whether to suppress the [CheckMap] call in [Run].
// It returns the value passed to an earlier call to [WithSuppressCheck],
// or false by default.
func SuppressCheck(ctx context.Context) bool {
	val, _ := ctx.Value(suppressCheckKey).(bool)
	return val
}