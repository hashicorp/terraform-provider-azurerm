// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package databaselink_test

import (
	"slices"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedredis/databaselink"
)

func TestForceUnlinkInvocations(t *testing.T) {
	tests := []struct {
		name            string
		intermediateIds []string
		idsToUnlink     []string
		expected        []databaselink.ForceUnlinkInvocation
	}{
		{
			name:            "intermediateIds: [a] idsToUnlink: [b,c]",
			intermediateIds: []string{"a"},
			idsToUnlink:     []string{"b", "c"},
			expected: []databaselink.ForceUnlinkInvocation{
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
			expected:        []databaselink.ForceUnlinkInvocation{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := databaselink.ForceUnlinkInvocations(tc.intermediateIds, tc.idsToUnlink)
			if !slices.EqualFunc(actual, tc.expected, func(a, b databaselink.ForceUnlinkInvocation) bool {
				return a.Id == b.Id && slices.Equal(a.Ids, b.Ids)
			}) {
				t.Errorf("expected %+v, got %+v", tc.expected, actual)
			}
		})
	}
}
