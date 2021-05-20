package validate

import (
	"strings"
	"testing"
)

func TestCommunicationServiceName(t *testing.T) {
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
			Input:    "-a",
			Expected: false,
		},
		{
			Input:    "1a",
			Expected: false,
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
			Input:    "ab",
			Expected: true,
		},
		{
			Input:    "Test",
			Expected: true,
		},
		{
			Input:    strings.Repeat("s", 63),
			Expected: true,
		},
		{
			Input:    strings.Repeat("s", 64),
			Expected: true,
		},
		{
			Input:    strings.Repeat("s", 65),
			Expected: false,
		},
	}

	for _, v := range testCases {
		_, errors := CommunicationServiceName(v.Input, "name")
		result := len(errors) == 0
		if result != v.Expected {
			t.Fatalf("Expected the result to be %t but got %t (and %d errors)", v.Expected, result, len(errors))
		}
	}
}
