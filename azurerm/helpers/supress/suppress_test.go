package supress

import (
	"testing"
)

func TestInternal_Suppress_CaseDifference(t *testing.T) {
	cases := []struct {
		StringA  string
		StringB  string
		Suppress bool
	}{
		{
			StringA:  "",
			StringB:  "",
			Suppress: true,
		},
		{
			StringA:  "ye old text",
			StringB:  "",
			Suppress: false,
		},
		{
			StringA:  "ye old text?",
			StringB:  "ye different text",
			Suppress: false,
		},
		{
			StringA:  "ye same text!",
			StringB:  "ye same text!",
			Suppress: true,
		},
		{
			StringA:  "ye old text?",
			StringB:  "Ye OLD texT?",
			Suppress: true,
		},
	}

	for _, tc := range cases {
		if CaseDifference("test", tc.StringA, tc.StringB, nil) != tc.Suppress {
			t.Fatalf("Expected CaseDifference to return %t for '%s' == '%s'", tc.Suppress, tc.StringA, tc.StringB)
		}
	}
}

func TestInternal_Supress_Rfc3339Time(t *testing.T) {
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
			t.Fatalf("Expected Rfc3339Time to return %t for '%s' == '%s'", tc.Suppress, tc.TimeA, tc.TimeB)
		}
	}
}
