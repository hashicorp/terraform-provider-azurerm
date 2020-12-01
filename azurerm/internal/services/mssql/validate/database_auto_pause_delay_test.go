package validate

import "testing"

func TestDatabaseAutoPauseDelay(t *testing.T) {
	testCases := []struct {
		input       string
		shouldError bool
	}{
		{"-1", false},
		{"-2", true},
		{"30", true},
		{"60", false},
		{"65", true},
		{"360", false},
		{"19900", true},
	}

	for _, test := range testCases {
		_, es := MsSqlDatabaseAutoPauseDelay(test.input, "name")

		if test.shouldError && len(es) == 0 {
			t.Fatalf("Expected validating name %q to fail", test.input)
		}
	}
}
