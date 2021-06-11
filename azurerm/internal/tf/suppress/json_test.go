package suppress

import "testing"

func TestJsonDiff(t *testing.T) {
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
			Name:     "simple same object",
			StringA:  "{\"field\": \"value\"}",
			StringB:  "{\"field\": \"value\"}",
			Suppress: true,
		},
		{
			Name:     "simple object whitespace diff",
			StringA:  "{\n\"field\":      \"value\"\n}",
			StringB:  "{\"field\": \"value\"}",
			Suppress: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			if JsonDiff("test", tc.StringA, tc.StringB, nil) != tc.Suppress {
				t.Fatalf("Expected JsonDiff to return %t for '%q' == '%q'", tc.Suppress, tc.StringA, tc.StringB)
			}
		})
	}
}
