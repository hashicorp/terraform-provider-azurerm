package validate

import (
	"testing"
)

func TestTagsWithLowerCaseKey(t *testing.T) {
	testData := []struct {
		input    map[string]interface{}
		expected bool
	}{
		{
			input:    map[string]interface{}{},
			expected: true,
		},
		{
			// basic example
			input: map[string]interface{}{
				"key1": "Value1",
				"key2": "VALUE",
			},
			expected: true,
		},
		{
			// contains upper case key
			input: map[string]interface{}{
				"KEY":  "value",
				"key2": "VALUE",
			},
			expected: false,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.input)

		_, errors := TagsWithLowerCaseKey(v.input, "tags")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}
