package validate

import "testing"

func TestRemoteFilePath(t *testing.T) {
	cases := []struct {
		Value  string
		Errors int
	}{
		{
			Value:  "bad",
			Errors: 1,
		},
		{
			Value:  "/good/file/path",
			Errors: 0,
		},
	}

	for _, tc := range cases {
		_, errors := RemoteFilePath(tc.Value, "unittest")

		if len(errors) != tc.Errors {
			t.Fatalf("Expected validateDataLakeStoreRemoteFilePath to trigger '%d' errors for '%s' - got '%d'", tc.Errors, tc.Value, len(errors))
		}
	}
}
