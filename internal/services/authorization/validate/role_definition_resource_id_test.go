// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import "testing"

func TestRoleDefinitionResourceID(t *testing.T) {
	testCases := []struct {
		name     string
		input    interface{}
		errCount int
	}{
		{
			name:     "tenant scoped",
			input:    "/providers/Microsoft.Authorization/roleDefinitions/acdd72a7-3385-48ef-bd42-f606fba81ae7",
			errCount: 0,
		},
		{
			name:     "resource group scoped",
			input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.Authorization/roleDefinitions/acdd72a7-3385-48ef-bd42-f606fba81ae7",
			errCount: 0,
		},
		{
			name:     "invalid format",
			input:    "not-a-role-definition-id",
			errCount: 1,
		},
		{
			name:     "non string",
			input:    123,
			errCount: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, errors := RoleDefinitionResourceID(tc.input, "role_definition_id")
			if len(errors) != tc.errCount {
				t.Fatalf("expected %d errors but got %d for %q", tc.errCount, len(errors), tc.name)
			}
		})
	}
}
