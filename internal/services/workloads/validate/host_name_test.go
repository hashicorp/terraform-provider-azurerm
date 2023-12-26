package validate

import (
	"strings"
	"testing"
)

func TestHostName(t *testing.T) {
	cases := []struct {
		Input string
		Valid bool
	}{
		{
			Input: "",
			Valid: false,
		},
		{
			Input: "testhost",
			Valid: true,
		},
		{
			Input: strings.Repeat("s", 12),
			Valid: true,
		},
		{
			Input: strings.Repeat("s", 13),
			Valid: true,
		},
		{
			Input: strings.Repeat("s", 14),
			Valid: false,
		},
	}

	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := HostName(tc.Input, "host_name")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}
