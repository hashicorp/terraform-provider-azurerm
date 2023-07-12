// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import "testing"

func TestElasticMonitorsName(t *testing.T) {
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
			// proper string
			input:    "hello",
			expected: true,
		},
		{
			// end with exclamation
			input:    "hello!",
			expected: false,
		},
		{
			// with hypen
			input:    "malcolm-in-the-middle",
			expected: true,
		},
		{
			// end with fullstop
			input:    "hello.",
			expected: false,
		},
		{
			// less than 32
			input:    "qwertyuioplkjhgfdsazxcv",
			expected: true,
		},
		{
			// with underscore
			input:    "qwertyuiop_jhgfdsazxcva",
			expected: true,
		},
		{
			// more than 32
			input:    "qwertyuioplkjhgfdsazxcvggeitofkjhkt",
			expected: false,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.input)

		_, errors := ElasticsearchName(v.input, "name")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}
