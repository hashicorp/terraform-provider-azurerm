package validate

import (
	"strings"
	"testing"
)

func TestTemplateSpecDescription(t *testing.T) {
	testCases := []struct {
		Input    string
		Expected bool
	}{
		{
			Input:    "",
			Expected: false,
		},
		{
			Input:    strings.Repeat("s", 4095),
			Expected: true,
		},
		{
			Input:    strings.Repeat("s", 4096),
			Expected: true,
		},
		{
			Input:    strings.Repeat("s", 4097),
			Expected: false,
		},
		{
			Input:    "Test Description",
			Expected: true,
		},
	}
	for _, v := range testCases {
		_, errors := TemplateSpecDescription(v.Input, "description")
		result := len(errors) == 0
		if result != v.Expected {
			t.Fatalf("Expected the result to be %t but got %t (and %d errors)", v.Expected, result, len(errors))
		}
	}
}
