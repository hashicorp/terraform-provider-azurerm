package databaselink_test

import (
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedredis/databaselink"
)

func TestHasDbToLink(t *testing.T) {
	tests := []struct {
		name     string
		from     []string
		to       []string
		expected bool
	}{
		{
			name:     "from: a,b -- to: a,b,c -- expected: true",
			from:     []string{"a", "b"},
			to:       []string{"a", "b", "c"},
			expected: true,
		},
		{
			name:     "from: a,b -- to: a,b -- expected: false",
			from:     []string{"a", "b"},
			to:       []string{"a", "b"},
			expected: false,
		},
		{
			name:     "from: a,b,c -- to: a,b -- expected: false",
			from:     []string{"a", "b", "c"},
			to:       []string{"a", "b"},
			expected: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := databaselink.HasDbToLink(tc.from, tc.to)
			if actual != tc.expected {
				t.Errorf("\nexpected: %v but got: %v", tc.expected, actual)
			}
		})
	}
}
