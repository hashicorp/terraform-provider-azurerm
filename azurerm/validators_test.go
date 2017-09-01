package azurerm

import (
	"testing"
)

func TestValidateRFC3339Date(t *testing.T) {
	cases := []struct {
		Value    string
		ErrCount int
	}{
		{
			Value:    "",
			ErrCount: 1,
		},
		{
			Value:    "Random",
			ErrCount: 1,
		},
		{
			Value:    "2017-01-01",
			ErrCount: 1,
		},
		{
			Value:    "2017-01-01T01:23:45",
			ErrCount: 1,
		},
		{
			Value:    "2017-01-01T01:23:45+00:00",
			ErrCount: 0,
		},
		{
			Value:    "2017-01-01T01:23:45Z",
			ErrCount: 0,
		},
	}

	for _, tc := range cases {
		_, errors := validateRFC3339Date(tc.Value, "example")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected validateRFC3339Date to trigger '%d' errors for '%s' - got '%d'", tc.ErrCount, tc.Value, len(errors))
		}
	}
}
