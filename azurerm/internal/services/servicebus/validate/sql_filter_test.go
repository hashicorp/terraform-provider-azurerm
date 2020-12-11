package validate

import (
	"strings"
	"testing"
)

func TestValidateSqlFilter(t *testing.T) {
	cases := []struct {
		Value       string
		ShouldError bool
	}{
		{
			Value:       "to='terraform'",
			ShouldError: false,
		},
		{
			Value:       "",
			ShouldError: true,
		},
		{
			Value:       strings.Repeat("user.foo='bar' AND ", 55)[:1040],
			ShouldError: true,
		},
	}

	for _, tc := range cases {
		_, errors := SqlFilter(tc.Value, "sql_filter")

		hasErrors := len(errors) > 0
		if hasErrors && !tc.ShouldError {
			t.Fatalf("Expected no errors but got %d for %q", errors[0], tc.Value)
		}

		if !hasErrors && tc.ShouldError {
			t.Fatalf("Expected no errors but got %d for %q", len(errors), tc.Value)
		}
	}
}
