package suppress

import "testing"

func TestHelper_Supress_Rfc3339Time(t *testing.T) {
	cases := []struct {
		TimeA    string
		TimeB    string
		Suppress bool
	}{
		{
			TimeA:    "",
			TimeB:    "",
			Suppress: false,
		},
		{
			TimeA:    "this is not a time",
			TimeB:    "neither is this",
			Suppress: false,
		},
		{
			TimeA:    "2000-01-01T01:23:45+00:00",
			TimeB:    "that is a valid time",
			Suppress: false,
		},
		{
			TimeA:    "2000-01-01T01:23:45+00:00",
			TimeB:    "1984-07-07T01:23:45+00:00",
			Suppress: false,
		},
		{
			TimeA:    "2000-01-01T01:23:45+00:00",
			TimeB:    "2000-01-01T01:23:45Z",
			Suppress: true,
		},
		{
			TimeA:    "2000-01-01T01:23:45-08:00",
			TimeB:    "2000-01-01T09:23:45Z",
			Suppress: true,
		},
	}

	for _, tc := range cases {
		if Rfc3339Time("test", tc.TimeA, tc.TimeB, nil) != tc.Suppress {
			t.Fatalf("Expected Rfc3339Time to return %t for '%q' == '%q'", tc.Suppress, tc.TimeA, tc.TimeB)
		}
	}
}
