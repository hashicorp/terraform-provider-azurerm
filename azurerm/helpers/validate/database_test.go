package validate

import "testing"

func TestDatabaseCollation(t *testing.T) {
	cases := []struct {
		Value  string
		Errors int
	}{
		{
			Value:  "en@US",
			Errors: 1,
		},
		{
			Value:  "en-US",
			Errors: 0,
		},
		{
			Value:  "en_US",
			Errors: 0,
		},
		{
			Value:  "en US",
			Errors: 0,
		},
		{
			Value:  "English_United States.1252",
			Errors: 0,
		},
	}

	for _, tc := range cases {
		_, errors := PostgresDatabaseCollation(tc.Value, "collation")
		if len(errors) != tc.Errors {
			t.Fatalf("Expected DatabaseCollation to trigger '%d' errors for '%s' - got '%d'", tc.Errors, tc.Value, len(errors))
		}
	}
}
