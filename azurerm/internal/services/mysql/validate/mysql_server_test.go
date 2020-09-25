package validate

import (
	"testing"
)

func TestValidateMysqlServerServerID(t *testing.T) {
	testData := []struct {
		input    string
		expected bool
	}{
		{
			// empty
			input:    "",
			expected: false,
		},
		{
			// invalid
			input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg",
			expected: false,
		},
		{
			// valid
			input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.DBforMySQL/servers/test-mysql",
			expected: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.input)

		_, errors := MySQLServerID(v.input, "name")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}

func TestValidateMysqlServerServerName(t *testing.T) {
	testData := []struct {
		input    string
		expected bool
	}{
		{
			// empty
			input:    "",
			expected: false,
		},
		{
			// basic example
			input:    "ab-c",
			expected: true,
		},
		{
			// can't contain upper case letter
			input:    "AbcD",
			expected: false,
		},
		{
			// can't start with a hyphen
			input:    "-abc",
			expected: false,
		},
		{
			// can't contain underscore
			input:    "ab_c",
			expected: false,
		},
		{
			// can't end with hyphen
			input:    "abc-",
			expected: false,
		},
		{
			// can not be shorter than 3 characters
			input:    "ab",
			expected: false,
		},
		{
			// can not be shorter than 3 characters (catching bad regex)
			input:    "a",
			expected: false,
		},
		{
			// 63 chars
			input:    "abcdefghijklmnopqrstuvwxyzabcdefabcdefghijklmnopqrstuvwxyzabcde",
			expected: true,
		},
		{
			// 64 chars
			input:    "abcdefghijklmnopqrstuvwxyzabcdefabcdefghijklmnopqrstuvwxyzabcdef",
			expected: false,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.input)

		_, errors := MySQLServerName(v.input, "name")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}
