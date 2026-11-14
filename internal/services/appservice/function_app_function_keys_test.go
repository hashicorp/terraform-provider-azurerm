// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package appservice

import (
	"reflect"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
)

func TestFlattenFunctionAppFunctionKeys(t *testing.T) {
	cases := []struct {
		name     string
		input    *webapps.StringDictionary
		expected map[string]string
	}{
		{
			name:     "nil response",
			input:    nil,
			expected: map[string]string{},
		},
		{
			name:     "nil properties",
			input:    &webapps.StringDictionary{},
			expected: map[string]string{},
		},
		{
			name: "default and custom keys",
			input: &webapps.StringDictionary{
				Properties: &map[string]string{
					"default": "abc",
					"ci":      "def",
				},
			},
			expected: map[string]string{
				"default": "abc",
				"ci":      "def",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual := flattenFunctionAppFunctionKeys(tc.input)
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Fatalf("expected %+v, got %+v", tc.expected, actual)
			}
		})
	}
}
