package databaselink_test

import (
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedredis/databaselink"
)

func TestDbIdsToUnlink(t *testing.T) {
	tests := []struct {
		name     string
		from     []string
		to       []string
		expected []string
	}{
		{
			name:     "from: a,b,c -- to: a,c -- expected: b",
			from:     []string{"a", "b", "c"},
			to:       []string{"a", "c"},
			expected: []string{"b"},
		},
		{
			name:     "from: a,b,c -- to: b,d,e,f -- expected: a,c",
			from:     []string{"a", "b", "c"},
			to:       []string{"b", "d", "e", "f"},
			expected: []string{"a", "c"},
		},
		{
			name:     "from: a,b,c -- to: a,b,c -- expected: []",
			from:     []string{"a", "b", "c"},
			to:       []string{"a", "b", "c"},
			expected: []string{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := databaselink.DbIdsToUnlink(tc.from, tc.to)
			if len(actual) != len(tc.expected) {
				t.Errorf("\nexpected length %d but got %d", len(tc.expected), len(actual))
			}
			matches := true
			for i, v := range actual {
				if v != tc.expected[i] {
					matches = false
					break
				}
			}
			if !matches {
				t.Errorf("\nexpected: %v\nactual: %v", tc.expected, actual)
			}
		})
	}
}
