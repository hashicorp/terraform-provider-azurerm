package databaselink_test

import (
	"slices"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedredis/databaselink"
)

func TestForceLinkInvocations(t *testing.T) {
	tests := []struct {
		name     string
		from     []string
		to       []string
		expected []databaselink.ForceLinkInvocation
	}{
		{
			name: "from: [a] to: [b]",
			from: []string{"a"},
			to:   []string{"b"},
			expected: []databaselink.ForceLinkInvocation{
				{
					Id:                "b",
					LinkedDatabaseIds: []string{"a", "b"},
				},
			},
		},
		{
			name: "from: [a] to: [b,c]",
			from: []string{"a"},
			to:   []string{"b", "c"},
			expected: []databaselink.ForceLinkInvocation{
				{
					Id:                "b",
					LinkedDatabaseIds: []string{"a", "b"},
				},
				{
					Id:                "c",
					LinkedDatabaseIds: []string{"a", "b", "c"},
				},
			},
		},
		{
			name: "from: [a,b] to: [c]",
			from: []string{"a", "b"},
			to:   []string{"c"},
			expected: []databaselink.ForceLinkInvocation{
				{
					Id:                "c",
					LinkedDatabaseIds: []string{"a", "b", "c"},
				},
			},
		},
		{
			name: "from: [a,b,c] to: [d,e]",
			from: []string{"a", "b", "c"},
			to:   []string{"d", "e"},
			expected: []databaselink.ForceLinkInvocation{
				{
					Id:                "d",
					LinkedDatabaseIds: []string{"a", "b", "c", "d"},
				},
				{
					Id:                "e",
					LinkedDatabaseIds: []string{"a", "b", "c", "d", "e"},
				},
			},
		},
		{
			name:     "from: [a,b,c] to: []",
			from:     []string{"a", "b", "c"},
			to:       []string{},
			expected: []databaselink.ForceLinkInvocation{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := databaselink.ForceLinkInvocations(tc.from, tc.to)
			if !slices.EqualFunc(actual, tc.expected, func(a, b databaselink.ForceLinkInvocation) bool {
				if a.Id != b.Id {
					return false
				}
				return slices.Equal(a.LinkedDatabaseIds, b.LinkedDatabaseIds)
			}) {
				t.Errorf("\nexpected: %v but got: %v", tc.expected, actual)
			}
		})
	}
}
