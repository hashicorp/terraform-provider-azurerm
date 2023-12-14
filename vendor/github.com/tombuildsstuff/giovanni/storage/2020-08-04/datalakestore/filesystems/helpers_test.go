package filesystems

import (
	"reflect"
	"testing"
)

func TestParseProperties(t *testing.T) {
	testData := []struct {
		name        string
		input       string
		expected    map[string]string
		expectError bool
	}{
		{
			name:        "no items",
			input:       "",
			expected:    map[string]string{},
			expectError: false,
		},
		{
			name:        "invalid item",
			input:       "hello",
			expectError: true,
		},
		{
			name:  "single item",
			input: "hello=world",
			expected: map[string]string{
				"hello": "world",
			},
		},
		{
			name:  "single-item-base64",
			input: "hello=aGVsbG8=",
			expected: map[string]string{
				"hello": "aGVsbG8=",
			},
			expectError: false,
		},
		{
			name:  "single-item-base64-multipleequals",
			input: "hello=d29uZGVybGFuZA==",
			expected: map[string]string{
				"hello": "d29uZGVybGFuZA==",
			},
			expectError: false,
		},
		{
			name:  "multiple-items-base64",
			input: "hello=d29uZGVybGFuZA==,private=ZXll",

			expected: map[string]string{
				"hello":   "d29uZGVybGFuZA==",
				"private": "ZXll",
			},
			expectError: false,
		},
	}

	for _, testCase := range testData {
		t.Logf("[DEBUG] Test %q", testCase.name)

		actual, err := parseProperties(testCase.input)
		if err != nil {
			if testCase.expectError {
				continue
			}

			t.Fatalf("[DEBUG] Didn't expect an error but got %s", err)
		}
		if !reflect.DeepEqual(testCase.expected, *actual) {
			t.Fatalf("Expected %+v but got %+v", testCase.expected, *actual)
		}
	}
}
