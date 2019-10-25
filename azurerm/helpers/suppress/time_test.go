package suppress

import "testing"

func TestRFC3339Time(t *testing.T) {
	cases := []struct {
		Name     string
		TimeA    string
		TimeB    string
		Suppress bool
	}{
		{
			Name:     "empty",
			TimeA:    "",
			TimeB:    "",
			Suppress: false,
		},
		{
			Name:     "neither are time",
			TimeA:    "this is not a time",
			TimeB:    "neither is this",
			Suppress: false,
		},
		{
			Name:     "time vs text",
			TimeA:    "2000-01-01T01:23:45+00:00",
			TimeB:    "that is a valid time",
			Suppress: false,
		},
		{
			Name:     "two different times",
			TimeA:    "2000-01-01T01:23:45+00:00",
			TimeB:    "1984-07-07T01:23:45+00:00",
			Suppress: false,
		},
		{
			Name:     "same time, different zone 1",
			TimeA:    "2000-01-01T01:23:45+00:00",
			TimeB:    "2000-01-01T01:23:45Z",
			Suppress: true,
		},
		{
			Name:     "same time, different zone 2",
			TimeA:    "2000-01-01T01:23:45-08:00",
			TimeB:    "2000-01-01T09:23:45Z",
			Suppress: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			if RFC3339Time("test", tc.TimeA, tc.TimeB, nil) != tc.Suppress {
				t.Fatalf("Expected RFC3339Time to return %t for '%q' == '%q'", tc.Suppress, tc.TimeA, tc.TimeB)
			}
		})
	}
}
