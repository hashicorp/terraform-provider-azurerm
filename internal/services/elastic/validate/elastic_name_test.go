// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import "testing"

func TestElasticName(t *testing.T) {
	testData := []struct {
		name     string
		input    interface{}
		expected bool
	}{
		{
			name:     "minimum length",
			input:    "ab",
			expected: true,
		},
		{
			name:     "less than minimum length",
			input:    "a",
			expected: false,
		},
		{
			name:     "maximum length",
			input:    "qwertyuioplkjhgfdsazxcvbnmqwerty",
			expected: true,
		},
		{
			name:     "more than maximum length",
			input:    "qwertyuioplkjhgfdsazxcvbnmqwertyx",
			expected: false,
		},
		{
			name:     "hyphen and underscore",
			input:    "elastic_name-01",
			expected: true,
		},
		{
			name:     "invalid character",
			input:    "elastic.name",
			expected: false,
		},
		{
			name:     "non-string value",
			input:    123,
			expected: false,
		},
	}

	for _, test := range testData {
		t.Run(test.name, func(t *testing.T) {
			_, errors := ElasticName(test.input, "name")
			actual := len(errors) == 0
			if test.expected != actual {
				t.Fatalf("expected %t but got %t", test.expected, actual)
			}
		})
	}
}
