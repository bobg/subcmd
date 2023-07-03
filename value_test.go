package subcmd

import (
	"context"
	"flag"
	"reflect"
	"strings"
	"testing"
)

func TestValue(t *testing.T) {
	c := valuetestcmd{
		x:    valuetestvalue{result: []string{"x"}},
		y:    valuetestvalue{result: []string{"y"}},
		pos1: valuetestvalue{result: []string{"pos1"}},
		pos2: valuetestvalue{result: []string{"pos2"}},
	}

	if err := Run(context.Background(), &c, []string{"a", "-x", "1,2,3", "4,5,6", "7,8", "9,10"}); err != nil {
		t.Fatal(err)
	}

	want := valuetestcmd{
		x:    valuetestvalue{result: []string{"1", "2", "3"}},
		y:    valuetestvalue{result: []string{"y"}},
		pos1: valuetestvalue{result: []string{"4", "5", "6"}},
		pos2: valuetestvalue{result: []string{"7", "8"}},
		rest: []string{"9,10"},
	}
	if !reflect.DeepEqual(c, want) {
		t.Errorf("got %+v, want %+v", c, want)
	}
}

type valuetestcmd struct {
	x, y, pos1, pos2 valuetestvalue
	rest             []string
}

func (c *valuetestcmd) Subcmds() Map {
	return Commands(
		"a", c.a, "", Params(
			"-x", Value, &c.x, "",
			"-y", Value, &c.y, "",
			"pos1", Value, &c.pos1, "",
			"pos2", Value, &c.pos2, "",
		),
	)
}

func (c *valuetestcmd) a(_ context.Context, x, y, pos1, pos2 flag.Value, rest []string) error {
	if x, _ := x.(*valuetestvalue); x != nil {
		c.x = *x
	}
	if y, _ := y.(*valuetestvalue); y != nil {
		c.y = *y
	}
	if pos1, _ := pos1.(*valuetestvalue); pos1 != nil {
		c.pos1 = *pos1
	}
	if pos2, _ := pos2.(*valuetestvalue); pos2 != nil {
		c.pos2 = *pos2
	}
	c.rest = rest
	return nil
}

type valuetestvalue struct {
	result []string
}

var _ flag.Value = &valuetestvalue{}

func (v *valuetestvalue) String() string {
	if v == nil {
		return ""
	}
	return strings.Join(v.result, ",")
}

func (v *valuetestvalue) Set(s string) error {
	v.result = strings.Split(s, ",")
	return nil
}
