package subcmd

import (
	"context"
	"fmt"
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

func TestDifferentValues(t *testing.T) {
	dflt := &valuetestvalue{result: []string{"dflt"}}
	cmd := new(valuetestcmd)
	cmd.subcmds = Commands(
		"b", cmd.differentValues, "", Params(
			"x", Value, dflt, "",
			"y", Value, dflt, "",
		),
	)
	if err := Run(context.Background(), cmd, []string{"b", "res1", "res2"}); err != nil {
		t.Fatal(err)
	}
}

type valuetestcmd struct {
	x, y, pos1, pos2 valuetestvalue
	rest             []string

	subcmds Map
}

func (c *valuetestcmd) Subcmds() Map {
	if c.subcmds != nil {
		return c.subcmds
	}

	return Commands(
		"a", c.a, "", Params(
			"-x", Value, &c.x, "",
			"-y", Value, &c.y, "",
			"pos1", Value, &c.pos1, "",
			"pos2", Value, &c.pos2, "",
		),
	)
}

func (c *valuetestcmd) a(_ context.Context, x, y, pos1, pos2 ValueType, rest []string) error {
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

func (c *valuetestcmd) differentValues(_ context.Context, xv, yv ValueType, _ []string) error {
	x, ok := xv.(*valuetestvalue)
	if !ok {
		return fmt.Errorf("unexpected type %T for x", xv)
	}
	y, ok := yv.(*valuetestvalue)
	if !ok {
		return fmt.Errorf("unexpected type %T for y", yv)
	}
	if len(x.result) != 1 || x.result[0] != "res1" {
		return fmt.Errorf("unexpected x value %v", x.result)
	}
	if len(y.result) != 1 || y.result[0] != "res2" {
		return fmt.Errorf("unexpected y value %v", y.result)
	}
	return nil
}

type valuetestvalue struct {
	result []string
}

var _ ValueType = &valuetestvalue{}

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

func (v *valuetestvalue) Copy() ValueType {
	result := &valuetestvalue{
		result: make([]string, len(v.result)),
	}
	copy(result.result, v.result)
	return result
}
