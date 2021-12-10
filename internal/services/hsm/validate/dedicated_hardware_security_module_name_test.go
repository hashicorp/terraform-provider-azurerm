package validate

import (
	"strings"
	"testing"
)

func TestDedicatedHardwareSecurityModuleName(t *testing.T) {
	testCases := []struct {
		Input    string
		Expected bool
	}{
		{
			Input:    "hello-world",
			Expected: true,
		},
		{
			Input:    "-hello-world",
			Expected: false,
		},
		{
			Input:    "hello-world-",
			Expected: false,
		},
		{
			Input:    "9hello-world",
			Expected: false,
		},
		{
			Input:    "hello-world9",
			Expected: true,
		},
		{
			Input:    "hello-world-test",
			Expected: true,
		},
		{
			Input:    "hello--world",
			Expected: false,
		},
		{
			Input:    strings.Repeat("a", 23),
			Expected: true,
		},
		{
			Input:    strings.Repeat("a", 24),
			Expected: true,
		},
		{
			Input:    strings.Repeat("a", 25),
			Expected: false,
		},
	}

	for _, v := range testCases {
		_, errors := DedicatedHardwareSecurityModuleName(v.Input, "name")
		result := len(errors) == 0
		if result != v.Expected {
			t.Fatalf("Expected the result to be %t but got %t (and %d errors)", v.Expected, result, len(errors))
		}
	}
}
