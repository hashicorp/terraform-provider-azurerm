package validate

import "testing"

func TestStorageSyncName(t *testing.T) {
	testCases := []struct {
		input       string
		shouldError bool
	}{
		{"Abd1", false},
		{"hu ub", false},
		{"dfj_ df-.Hj12", false},
		{"ui.", true},
		{"76jhu#", true},
		{"df_ *-", true},
		{"dfAd1 ", true},
		{" dfAd1", false},
	}

	for _, test := range testCases {
		_, es := StorageSyncName(test.input, "name")

		if test.shouldError && len(es) == 0 {
			t.Fatalf("Expected validating name %q to fail", test.input)
		}
	}
}
