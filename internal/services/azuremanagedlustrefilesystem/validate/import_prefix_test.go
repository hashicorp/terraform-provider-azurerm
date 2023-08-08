package validate

import (
	"testing"
)

func TestImportPrefix(t *testing.T) {
	testCases := []struct {
		Input    string
		Expected bool
	}{
		{
			Input:    "",
			Expected: false,
		},
		{
			Input:    "test",
			Expected: false,
		},
		{
			Input:    "/",
			Expected: true,
		},
		{
			Input:    "/example",
			Expected: true,
		},
	}

	for _, v := range testCases {
		_, errors := ImportPrefix(v.Input, "import_prefix")
		result := len(errors) == 0
		if result != v.Expected {
			t.Fatalf("Expected the result to be %t but got %t (and %d errors)", v.Expected, result, len(errors))
		}
	}
}
