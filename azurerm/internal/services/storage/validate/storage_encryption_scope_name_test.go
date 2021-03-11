package validate

import "testing"

func TestStorageEncryptionScopeName(t *testing.T) {
	testCases := []struct {
		input       string
		shouldError bool
	}{
		{"ad", true},
		{"a", true},
		{"", true},
		{"asdf_12", true},
		{"asdf-12", true},
		{"asdf.12", true},
		{"asdf", false},
		{" aSDf12", false},
	}

	for _, test := range testCases {
		_, es := StorageEncryptionScopeName(test.input, "name")

		if test.shouldError && len(es) == 0 {
			t.Fatalf("Expected validating name %q to fail", test.input)
		}
	}
}
