// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package databaselink_test

import (
	"slices"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedredis/databaselink"
)

func TestLinkUnlink(t *testing.T) {
	tests := []struct {
		name                    string
		fromIds                 []string
		toIds                   []string
		expectedIdsToUnlink     []string
		expectedIntermediateIds []string
		expectedIdsToLink       []string
	}{
		{
			name:                    "fromIds: [a,b,c], toIds: [a,c], expected: ([b], [a,c], [])",
			fromIds:                 []string{"a", "b", "c"},
			toIds:                   []string{"a", "c"},
			expectedIdsToUnlink:     []string{"b"},
			expectedIntermediateIds: []string{"a", "c"},
			expectedIdsToLink:       []string{},
		},
		{
			name:                    "fromIds: [a], toIds: [a,b], expected: ([], [a], [b])",
			fromIds:                 []string{"a"},
			toIds:                   []string{"a", "b"},
			expectedIdsToUnlink:     []string{},
			expectedIntermediateIds: []string{"a"},
			expectedIdsToLink:       []string{"b"},
		},
		{
			name:                    "fromIds: [a], toIds: [a,b,c], expected: ([], [a], [b,c])",
			fromIds:                 []string{"a"},
			toIds:                   []string{"a", "b", "c"},
			expectedIdsToUnlink:     []string{},
			expectedIntermediateIds: []string{"a"},
			expectedIdsToLink:       []string{"b", "c"},
		},
		{
			name:                    "fromIds: [a,b,c], toIds: [b,d,e,f], expected: ([a,c], [b], [d,e,f])",
			fromIds:                 []string{"a", "b", "c"},
			toIds:                   []string{"b", "d", "e", "f"},
			expectedIdsToUnlink:     []string{"a", "c"},
			expectedIntermediateIds: []string{"b"},
			expectedIdsToLink:       []string{"d", "e", "f"},
		},
		{
			name:                    "fromIds: [a,b,c], toIds: [a], expected: ([b,c], [a], [])",
			fromIds:                 []string{"a", "b", "c"},
			toIds:                   []string{"a"},
			expectedIdsToUnlink:     []string{"b", "c"},
			expectedIntermediateIds: []string{"a"},
			expectedIdsToLink:       []string{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actualIdsToUnlink, actualIntermediateIds, actualIdsToLink := databaselink.LinkUnlink(tc.fromIds, tc.toIds)
			if !slices.Equal(tc.expectedIdsToUnlink, actualIdsToUnlink) {
				t.Errorf("\nexpected idsToUnlink: %v\nactual: %v", tc.expectedIdsToUnlink, actualIdsToUnlink)
			}
			if !slices.Equal(tc.expectedIntermediateIds, actualIntermediateIds) {
				t.Errorf("\nexpected intermediateIds: %v\nactual: %v", tc.expectedIntermediateIds, actualIntermediateIds)
			}
			if !slices.Equal(tc.expectedIdsToLink, actualIdsToLink) {
				t.Errorf("\nexpected idsToLink: %v\nactual: %v", tc.expectedIdsToLink, actualIdsToLink)
			}
		})
	}
}
