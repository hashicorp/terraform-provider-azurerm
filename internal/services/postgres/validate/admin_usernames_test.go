package validate

import (
	"testing"
)

func TestValidateAdminUsernames(t *testing.T) {
	testData := []struct {
		input    string
		expected bool
	}{
		{
			// guest
			input:    "guest",
			expected: false,
		},
		{
			// basic example
			input:    "blah",
			expected: true,
		},
		{
			// contains pg_
			input:    "pg_blah",
			expected: false,
		},
		{
			// azure_pg_admin
			input:    "azure_pg_admin",
			expected: false,
		},
		{
			// Capitalised example
			input:    "Azure_superuser",
			expected: false,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.input)

		_, errors := AdminUsernames(v.input, "name")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}
