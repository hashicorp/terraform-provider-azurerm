package validate

import (
	"strings"
	"testing"
)

func TestSAPFQDN(t *testing.T) {
	cases := []struct {
		Input string
		Valid bool
	}{
		{
			Input: "",
			Valid: false,
		},
		{
			Input: "sap.contoso.com",
			Valid: true,
		},
		{
			Input: strings.Repeat("s", 33),
			Valid: true,
		},
		{
			Input: strings.Repeat("s", 34),
			Valid: true,
		},
		{
			Input: strings.Repeat("s", 35),
			Valid: false,
		},
	}

	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := SAPFQDN(tc.Input, "sap_fqdn")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}
