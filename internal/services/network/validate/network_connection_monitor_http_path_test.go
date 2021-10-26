package validate

import "testing"

func TestNetworkConnectionMonitorHttpPath(t *testing.T) {
	cases := []struct {
		Value  string
		Errors int
	}{
		{
			Value:  "",
			Errors: 1,
		},
		{
			Value:  "a/b",
			Errors: 1,
		},
		{
			Value:  "/ab/b1/",
			Errors: 0,
		},
		{
			Value:  "/a/b",
			Errors: 0,
		},
		{
			Value:  "http://www.pluginsdk.io",
			Errors: 1,
		},
		{
			Value:  "/a/b/",
			Errors: 0,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Value, func(t *testing.T) {
			_, errors := NetworkConnectionMonitorHttpPath(tc.Value, "path")

			if len(errors) != tc.Errors {
				t.Fatalf("Expected Path to return %d error(s) not %d", tc.Errors, len(errors))
			}
		})
	}
}
