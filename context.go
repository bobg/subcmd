package subcmd

import (
	"context"
	"flag"
)

type ctxkey int

const (
	fsKey ctxkey = iota + 1
	subcmdPairListKey
)

func withFlagSet(ctx context.Context, fs *flag.FlagSet) context.Context {
	return context.WithValue(ctx, fsKey, fs)
}

// FlagSet produces the [flag.FlagSet] used in a call to a [Subcmd] function.
func FlagSet(ctx context.Context) *flag.FlagSet {
	val := ctx.Value(fsKey)
	return val.(*flag.FlagSet)
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
