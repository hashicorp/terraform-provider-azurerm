package utils

import (
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestToPtr(t *testing.T) {
	cases := []struct {
		in  interface{}
		out interface{}
	}{
		{
			in:  1,
			out: Int(1),
		},
		{
			in:  "a",
			out: String("a"),
		},
		{
			in:  true,
			out: Bool(true),
		},
	}
	for idx, c := range cases {
		out := ToPtr(c.in)
		if !reflect.DeepEqual(out, c.out) {
			t.Fatalf("%d failed\nexpected:\n%s\nactual:\n%s\n", idx, spew.Sdump(out), spew.Sdump(c.out))
		}
	}

}
