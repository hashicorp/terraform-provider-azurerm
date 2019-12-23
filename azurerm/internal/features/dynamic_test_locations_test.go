package features

import (
	"os"
	"testing"
)

func TestDynamicTestLocations(t *testing.T) {
	testData := []struct {
		name     string
		value    string
		expected bool
	}{
		{
			name:     "unset",
			value:    "",
			expected: false,
		},
		{
			name:     "disabled lower-case",
			value:    "false",
			expected: false,
		},
		{
			name:     "disabled upper-case",
			value:    "FALSE",
			expected: false,
		},
		{
			name:     "enabled lower-case",
			value:    "true",
			expected: true,
		},
		{
			name:     "enabled upper-case",
			value:    "TRUE",
			expected: true,
		},
		{
			name:     "invalid",
			value:    "pandas",
			expected: false,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Test %q..", v.name)

		os.Setenv("ARM_PROVIDER_DYNAMIC_TEST", v.value)
		actual := UseDynamicTestLocations()
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}
