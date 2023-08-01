package validate

import (
	"strings"
	"testing"
)

func TestManagedLustreFileSystemName(t *testing.T) {
	testCases := []struct {
		Input    string
		Expected bool
	}{
		{
			Input:    "",
			Expected: false,
		},
		{
			Input:    "t",
			Expected: false,
		},
		{
			Input:    "test",
			Expected: true,
		},
		{
			Input:    "test_123-test",
			Expected: true,
		},
		{
			Input:    "_123test",
			Expected: false,
		},
		{
			Input:    "test123_",
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
	}

	for _, v := range testCases {
		_, errors := ManagedLustreFileSystemName(v.Input, "name")
		result := len(errors) == 0
		if result != v.Expected {
			t.Fatalf("Expected the result to be %t but got %t (and %d errors)", v.Expected, result, len(errors))
		}
	}
}
