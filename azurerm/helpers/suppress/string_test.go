package suppress

import "testing"

func TestCaseDifference(t *testing.T) {
	cases := []struct {
		Name     string
		StringA  string
		StringB  string
		Suppress bool
	}{
		{
			Name:     "empty",
			StringA:  "",
			StringB:  "",
			Suppress: true,
		},
		{
			Name:     "empty vs text",
			StringA:  "ye old text",
			StringB:  "",
			Suppress: false,
		},
		{
			Name:     "different text",
			StringA:  "ye old text?",
			StringB:  "ye different text",
			Suppress: false,
		},
		{
			Name:     "same text",
			StringA:  "ye same text!",
			StringB:  "ye same text!",
			Suppress: true,
		},
		{
			Name:     "same text different case",
			StringA:  "ye old text?",
			StringB:  "Ye OLD texT?",
			Suppress: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			if CaseDifference("test", tc.StringA, tc.StringB, nil) != tc.Suppress {
				t.Fatalf("Expected CaseDifference to return %t for '%q' == '%q'", tc.Suppress, tc.StringA, tc.StringB)
			}
		})
	}
}
