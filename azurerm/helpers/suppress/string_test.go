package suppress

import "testing"

func TestHelper_Suppress_CaseDifference(t *testing.T) {
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
			t.Fatalf("Expected CaseDifference to return %t for '%q' == '%q'", tc.Suppress, tc.StringA, tc.StringB)
		}
	}
}
