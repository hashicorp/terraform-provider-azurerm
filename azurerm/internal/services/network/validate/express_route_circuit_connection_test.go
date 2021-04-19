package validate

import (
	"strings"
	"testing"
)

func TestExpressRouteCircuitConnectionName(t *testing.T) {
	testCases := []struct {
		Input    string
		Expected bool
	}{
		{
			Input:    "",
			Expected: false,
		},
		{
			Input:    "a",
			Expected: true,
		},
		{
			Input:    "2",
			Expected: true,
		},
		{
			Input:    "_",
			Expected: false,
		},
		{
			Input:    "a_",
			Expected: true,
		},
		{
			Input:    "2_",
			Expected: true,
		},
		{
			Input:    "_a",
			Expected: false,
		},
		{
			Input:    "1a",
			Expected: true,
		},
		{
			Input:    "a2",
			Expected: true,
		},
		{
			Input:    "a-",
			Expected: false,
		},
		{
			Input:    "a-b",
			Expected: true,
		},
		{
			Input:    "a.b",
			Expected: true,
		},
		{
			Input:    "Test",
			Expected: true,
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
	}

	for _, v := range testCases {
		_, errors := ExpressRouteCircuitConnectionName(v.Input, "name")
		result := len(errors) == 0
		if result != v.Expected {
			t.Fatalf("Expected the result to be %t but got %t (and %d errors)", v.Expected, result, len(errors))
		}
	}
}
