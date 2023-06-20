package validate

import (
	"testing"
)

func TestGoogleClientID(t *testing.T) {
	cases := []struct {
		Value    string
		ErrCount int
	}{
		{
			Value:    "",
			ErrCount: 1,
		},
		{
			Value:    "invalid client id",
			ErrCount: 1,
		},
		{
			Value:    "00000000-0000-0000-0000-000000000000",
			ErrCount: 1,
		},
		{
			Value:    "123456789000-abcd12ef34gh56ij78900klmnopqrst0.apps.googleusercontent.org",
			ErrCount: 1,
		},
		{
			Value:    "123456789000-abcd12ef34gh56ij78900klmnopqrst0.apps.googleusercontent.com",
			ErrCount: 0,
		},
	}

	for _, tc := range cases {
		_, errors := GoogleClientID(tc.Value, "client_id")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected the Google Client ID %s to trigger a validation error", tc.Value)
		}
	}
}
