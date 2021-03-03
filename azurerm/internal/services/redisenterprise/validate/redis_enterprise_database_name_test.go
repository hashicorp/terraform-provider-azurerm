package validate

import "testing"

func TestRedisEnterpriseDatabaseName(t *testing.T) {
	testData := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "Invalid Empty database name",
			input:    "",
			expected: false,
		},
		{
			name:     "Invalid database name",
			input:    "My Database",
			expected: false,
		},
		{
			name:     "Valid database name",
			input:    "default",
			expected: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.name)

		_, errors := RedisEnterpriseDatabaseName(v.input, "name")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}
