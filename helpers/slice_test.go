package helpers

import (
	"reflect"
	"testing"
)

func TestRemoveFromStringArray(t *testing.T) {
	testCases := []struct {
		name     string
		elements []string
		remove   string
		expected []string
	}{
		{
			name:     "element not found",
			elements: []string{"a", "b", "c"},
			remove:   "z",
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "element found once",
			elements: []string{"a", "b", "c"},
			remove:   "b",
			expected: []string{"a", "c"},
		},
		{
			name:     "element found multiple times",
			elements: []string{"a", "b", "b", "c"},
			remove:   "b",
			expected: []string{"a", "c"},
		},
		{
			name:     "element at beginning",
			elements: []string{"a", "b", "c"},
			remove:   "a",
			expected: []string{"b", "c"},
		},
		{
			name:     "element at end",
			elements: []string{"a", "b", "c"},
			remove:   "c",
			expected: []string{"a", "b"},
		},
		{
			name:     "empty slice",
			elements: []string{},
			remove:   "a",
			expected: []string{},
		},
		{
			name:     "nil slice",
			elements: nil,
			remove:   "a",
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := RemoveFromStringArray(tc.elements, tc.remove)
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Fatalf("expected: %v, got: %v", tc.expected, actual)
			}
		})
	}
}

func TestSliceContainsValue(t *testing.T) {
	testCases := []struct {
		name     string
		input    []string
		value    string
		expected bool
	}{
		{
			name:     "value is present",
			input:    []string{"a", "b", "c"},
			value:    "b",
			expected: true,
		},
		{
			name:     "value is not present",
			input:    []string{"a", "b", "c"},
			value:    "z",
			expected: false,
		},
		{
			name:     "value is present multiple times",
			input:    []string{"a", "b", "b", "c"},
			value:    "b",
			expected: true,
		},
		{
			name:     "empty slice",
			input:    []string{},
			value:    "a",
			expected: false,
		},
		{
			name:     "nil slice",
			input:    nil,
			value:    "a",
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := SliceContainsValue(tc.input, tc.value)
			if actual != tc.expected {
				t.Fatalf("expected: %v, got: %v", tc.expected, actual)
			}
		})
	}
}
