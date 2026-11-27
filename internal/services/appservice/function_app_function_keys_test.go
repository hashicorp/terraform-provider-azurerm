// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package appservice

import (
	"reflect"
	"testing"
)

func TestParseFunctionAppFunctionKeys(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected map[string]string
		error    bool
	}{
		{
			name:     "empty body",
			input:    "",
			expected: map[string]string{},
		},
		{
			name:     "empty object",
			input:    `{}`,
			expected: map[string]string{},
		},
		{
			name:  "flat object with default and custom keys",
			input: `{"default": "abc", "ci": "def"}`,
			expected: map[string]string{
				"default": "abc",
				"ci":      "def",
			},
		},
		{
			name:  "documented StringDictionary shape",
			input: `{"id": "/subscriptions/xxx", "name": "keys", "properties": {"default": "abc", "ci": "def"}}`,
			expected: map[string]string{
				"default": "abc",
				"ci":      "def",
			},
		},
		{
			name:     "documented StringDictionary shape without keys",
			input:    `{"id": "/subscriptions/xxx", "properties": {}}`,
			expected: map[string]string{},
		},
		{
			name:  "flat object ignores envelope fields and non-string values",
			input: `{"id": null, "name": null, "type": null, "kind": null, "default": "abc", "other": 1}`,
			expected: map[string]string{
				"default": "abc",
			},
		},
		{
			name:  "invalid json",
			input: `not-json`,
			error: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := parseFunctionAppFunctionKeys([]byte(tc.input))
			if tc.error {
				if err == nil {
					t.Fatalf("expected an error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Fatalf("expected %+v, got %+v", tc.expected, actual)
			}
		})
	}
}
