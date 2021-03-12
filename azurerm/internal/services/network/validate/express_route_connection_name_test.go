package validate

import (
	"strings"
	"testing"
)

func TestExpressRouteConnectionName(t *testing.T) {
	testCases := []struct {
		Input    string
		Expected bool
	}{
		{
			Input:    "",
			Expected: false,
		},
		{
			Input:    strings.Repeat("s", 79),
			Expected: true,
		},
		{
			Input:    strings.Repeat("s", 80),
			Expected: true,
		},
		{
			Input:    strings.Repeat("s", 81),
			Expected: false,
		},
		{
			Input:    "_",
			Expected: false,
		},
		{
			Input:    "a",
			Expected: true,
		},
		{
			Input:    "a_",
			Expected: true,
		},
		{
			Input:    "ab",
			Expected: true,
		},
		{
			Input:    "abc",
			Expected: true,
		},
	}
	for _, v := range testCases {
		_, errors := ExpressRouteConnectionName(v.Input, "name")
		result := len(errors) == 0
		if result != v.Expected {
			t.Fatalf("Expected the result to be %t but got %t (and %d errors)", v.Expected, result, len(errors))
		}
	}
}
