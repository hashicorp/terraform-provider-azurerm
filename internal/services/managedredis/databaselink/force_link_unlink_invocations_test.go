// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

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
		expected []databaselink.LinkUnlinkInvocation
	}{
		{
			name: "from: [a] to: [b]",
			from: []string{"a"},
			to:   []string{"b"},
			expected: []databaselink.LinkUnlinkInvocation{
				{
					Id:  "b",
					Ids: []string{"a", "b"},
				},
			},
		},
		{
			name: "from: [a] to: [b,c]",
			from: []string{"a"},
			to:   []string{"b", "c"},
			expected: []databaselink.LinkUnlinkInvocation{
				{
					Id:  "b",
					Ids: []string{"a", "b"},
				},
				{
					Id:  "c",
					Ids: []string{"a", "b", "c"},
				},
			},
		},
		{
			name: "from: [a,b] to: [c]",
			from: []string{"a", "b"},
			to:   []string{"c"},
			expected: []databaselink.LinkUnlinkInvocation{
				{
					Id:  "c",
					Ids: []string{"a", "b", "c"},
				},
			},
		},
		{
			name: "from: [a,b,c] to: [d,e]",
			from: []string{"a", "b", "c"},
			to:   []string{"d", "e"},
			expected: []databaselink.LinkUnlinkInvocation{
				{
					Id:  "d",
					Ids: []string{"a", "b", "c", "d"},
				},
				{
					Id:  "e",
					Ids: []string{"a", "b", "c", "d", "e"},
				},
			},
		},
		{
			name:     "from: [a,b,c] to: []",
			from:     []string{"a", "b", "c"},
			to:       []string{},
			expected: []databaselink.LinkUnlinkInvocation{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := databaselink.ForceLinkInvocations(tc.from, tc.to)
			if !slices.EqualFunc(actual, tc.expected, func(a, b databaselink.LinkUnlinkInvocation) bool {
				if a.Id != b.Id {
					return false
				}
				return slices.Equal(a.Ids, b.Ids)
			}) {
				t.Errorf("\nexpected: %v but got: %v", tc.expected, actual)
			}
		})
	}
}

func TestForceUnlinkInvocations(t *testing.T) {
	tests := []struct {
		name            string
		intermediateIds []string
		idsToUnlink     []string
		expected        []databaselink.LinkUnlinkInvocation
	}{
		{
			name:            "intermediateIds: [a] idsToUnlink: [b,c]",
			intermediateIds: []string{"a"},
			idsToUnlink:     []string{"b", "c"},
			expected: []databaselink.LinkUnlinkInvocation{
				{
					Id:  "a",
					Ids: []string{"b"},
				},
				{
					Id:  "a",
					Ids: []string{"c"},
				},
			},
		},
		{
			name:            "intermediateIds: [a] idsToUnlink: []",
			intermediateIds: []string{"a"},
			idsToUnlink:     []string{},
			expected:        []databaselink.LinkUnlinkInvocation{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := databaselink.ForceUnlinkInvocations(tc.intermediateIds, tc.idsToUnlink)
			if !slices.EqualFunc(actual, tc.expected, func(a, b databaselink.LinkUnlinkInvocation) bool {
				return a.Id == b.Id && slices.Equal(a.Ids, b.Ids)
			}) {
				t.Errorf("expected %+v, got %+v", tc.expected, actual)
			}
		})
	}
}
