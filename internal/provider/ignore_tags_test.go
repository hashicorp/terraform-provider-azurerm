// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"sort"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ignoreTagsBlock(keys, prefixes []string) []interface{} {
	toIface := func(in []string) []interface{} {
		out := make([]interface{}, len(in))
		for i, v := range in {
			out[i] = v
		}
		return out
	}
	return []interface{}{
		map[string]interface{}{
			"keys":         schema.NewSet(schema.HashString, toIface(keys)),
			"key_prefixes": schema.NewSet(schema.HashString, toIface(prefixes)),
		},
	}
}

func TestExpandIgnoreTags(t *testing.T) {
	testCases := []struct {
		name             string
		input            []interface{}
		expectNil        bool
		expectedKeys     []string
		expectedPrefixes []string
	}{
		{
			name:      "absent block returns nil",
			input:     []interface{}{},
			expectNil: true,
		},
		{
			name:      "empty block returns nil",
			input:     ignoreTagsBlock(nil, nil),
			expectNil: true,
		},
		{
			name:         "keys only",
			input:        ignoreTagsBlock([]string{"createdBy", "owner"}, nil),
			expectedKeys: []string{"createdBy", "owner"},
		},
		{
			name:             "prefixes only",
			input:            ignoreTagsBlock(nil, []string{"azure-policy-"}),
			expectedPrefixes: []string{"azure-policy-"},
		},
		{
			name:             "keys and prefixes",
			input:            ignoreTagsBlock([]string{"createdBy"}, []string{"internal:", "azure-policy-"}),
			expectedKeys:     []string{"createdBy"},
			expectedPrefixes: []string{"azure-policy-", "internal:"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := expandIgnoreTags(tc.input)
			if tc.expectNil {
				if result != nil {
					t.Fatalf("expected nil, got %+v", result)
				}
				return
			}
			if result == nil {
				t.Fatalf("expected non-nil IgnoreConfig")
			}

			gotKeys := append([]string{}, result.Keys...)
			gotPrefixes := append([]string{}, result.KeyPrefixes...)
			sort.Strings(gotKeys)
			sort.Strings(gotPrefixes)
			sort.Strings(tc.expectedKeys)
			sort.Strings(tc.expectedPrefixes)

			if !equalStringSlices(gotKeys, tc.expectedKeys) {
				t.Fatalf("keys: expected %v, got %v", tc.expectedKeys, gotKeys)
			}
			if !equalStringSlices(gotPrefixes, tc.expectedPrefixes) {
				t.Fatalf("key_prefixes: expected %v, got %v", tc.expectedPrefixes, gotPrefixes)
			}
		})
	}
}

func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
