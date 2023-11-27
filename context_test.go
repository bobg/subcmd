package subcmd

import (
	"context"
	"testing"
)

func TestSubcmdPairs(t *testing.T) {
	ctx := context.Background()

	got := subcmdPairList(ctx)
	if len(got) != 0 {
		t.Errorf("got %d subcmd pair(s), want 0", len(got))
	}

	subcmd1 := Subcmd{Desc: "foo"}
	ctx = addSubcmdPair(ctx, "foo", subcmd1)
	got = subcmdPairList(ctx)
	if len(got) != 1 {
		t.Errorf("got %d subcmd pair(s), want 1", len(got))
	} else if got[0].subcmd.Desc != "foo" {
		t.Errorf(`got subcmd pair with description "%s", want "foo"`, got[0].subcmd.Desc)
	}

	subcmd2 := Subcmd{Desc: "bar"}
	ctx = addSubcmdPair(ctx, "bar", subcmd2)
	got = subcmdPairList(ctx)
	if len(got) != 2 {
		t.Errorf("got %d subcmd pair(s), want 2", len(got))
	} else {
		if got[0].subcmd.Desc != "foo" {
			t.Errorf(`subcmd pair 0 has description "%s", want "foo"`, got[0].subcmd.Desc)
		}
		if got[1].subcmd.Desc != "bar" {
			t.Errorf(`subcmd pair 1 has description "%s", want "var"`, got[1].subcmd.Desc)
		}
	}
}

func TestSuppressCheck(t *testing.T) {
	ctx := context.Background()

	if SuppressCheck(ctx) {
		t.Error("SuppressCheck(ctx) is true, want false")
	}
	ctx2 := WithSuppressCheck(ctx, false)
	if SuppressCheck(ctx2) {
		t.Error("SuppressCheck(ctx2) is true, want false")
	}
	ctx3 := WithSuppressCheck(ctx, true)
	if !SuppressCheck(ctx3) {
		t.Error("SuppressCheck(ctx3) is false, want true")
	}
}
