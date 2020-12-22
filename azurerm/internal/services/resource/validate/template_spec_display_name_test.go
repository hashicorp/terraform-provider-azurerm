package validate

import (
	"strings"
	"testing"
)

func TestTemplateSpecDisplayName(t *testing.T) {
	testCases := []struct {
		Input    string
		Expected bool
	}{
		{
			Input:    "",
			Expected: false,
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
		{
			Input:    "Test Display Name",
			Expected: true,
		},
	}
	for _, v := range testCases {
		_, errors := TemplateSpecDisplayName(v.Input, "display_name")
		result := len(errors) == 0
		if result != v.Expected {
			t.Fatalf("Expected the result to be %t but got %t (and %d errors)", v.Expected, result, len(errors))
		}
	}
}
