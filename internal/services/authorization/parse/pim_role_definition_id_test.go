// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package parse

import "testing"

func TestPimRoleDefinitionID(t *testing.T) {
	cases := []struct {
		name          string
		input         string
		expectedScope string
		expectedUUID  string
		wantErr       bool
	}{
		{
			name:          "subscription-scoped",
			input:         "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/roleDefinitions/acdd72a7-3385-48ef-bd42-f606fba81ae7",
			expectedScope: "/subscriptions/00000000-0000-0000-0000-000000000000",
			expectedUUID:  "acdd72a7-3385-48ef-bd42-f606fba81ae7",
		},
		{
			name:          "resource group-scoped",
			input:         "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myRG/providers/Microsoft.Authorization/roleDefinitions/acdd72a7-3385-48ef-bd42-f606fba81ae7",
			expectedScope: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myRG",
			expectedUUID:  "acdd72a7-3385-48ef-bd42-f606fba81ae7",
		},
		{
			name:          "resource-scoped (virtual network)",
			input:         "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myRG/providers/Microsoft.Network/virtualNetworks/myVNet/providers/Microsoft.Authorization/roleDefinitions/acdd72a7-3385-48ef-bd42-f606fba81ae7",
			expectedScope: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myRG/providers/Microsoft.Network/virtualNetworks/myVNet",
			expectedUUID:  "acdd72a7-3385-48ef-bd42-f606fba81ae7",
		},
		{
			name:          "resource-scoped (storage account)",
			input:         "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myRG/providers/Microsoft.Storage/storageAccounts/myAccount/providers/Microsoft.Authorization/roleDefinitions/acdd72a7-3385-48ef-bd42-f606fba81ae7",
			expectedScope: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myRG/providers/Microsoft.Storage/storageAccounts/myAccount",
			expectedUUID:  "acdd72a7-3385-48ef-bd42-f606fba81ae7",
		},
		{
			name:          "unscoped (provider-scoped)",
			input:         "/providers/Microsoft.Authorization/roleDefinitions/acdd72a7-3385-48ef-bd42-f606fba81ae7",
			expectedScope: "",
			expectedUUID:  "acdd72a7-3385-48ef-bd42-f606fba81ae7",
		},
		{
			name:          "mixed case",
			input:         "/Providers/Microsoft.Authorization/RoleDefinitions/ACDD72A7-3385-48EF-BD42-F606FBA81AE7",
			expectedScope: "",
			expectedUUID:  "acdd72a7-3385-48ef-bd42-f606fba81ae7",
		},
		{
			name:    "invalid format",
			input:   "not-a-role-definition-id",
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			id, err := PimRoleDefinitionID(tc.input)
			if tc.wantErr {
				if err == nil {
					t.Errorf("expected error but got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if id.Scope != tc.expectedScope {
				t.Errorf("scope: got %q, want %q", id.Scope, tc.expectedScope)
			}
			if id.RoleDefinitionUuid != tc.expectedUUID {
				t.Errorf("uuid: got %q, want %q", id.RoleDefinitionUuid, tc.expectedUUID)
			}
		})
	}
}

func TestPimRoleDefinitionIdsMatch(t *testing.T) {
	cases := []struct {
		name     string
		a        string
		b        string
		expected bool
	}{
		{
			name:     "identical scoped IDs",
			a:        "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/roleDefinitions/acdd72a7-3385-48ef-bd42-f606fba81ae7",
			b:        "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/roleDefinitions/acdd72a7-3385-48ef-bd42-f606fba81ae7",
			expected: true,
		},
		{
			name:     "unscoped vs scoped same UUID",
			a:        "/providers/Microsoft.Authorization/roleDefinitions/acdd72a7-3385-48ef-bd42-f606fba81ae7",
			b:        "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/roleDefinitions/acdd72a7-3385-48ef-bd42-f606fba81ae7",
			expected: true,
		},
		{
			name:     "scoped vs unscoped same UUID",
			a:        "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/roleDefinitions/acdd72a7-3385-48ef-bd42-f606fba81ae7",
			b:        "/providers/Microsoft.Authorization/roleDefinitions/acdd72a7-3385-48ef-bd42-f606fba81ae7",
			expected: true,
		},
		{
			name:     "case insensitive",
			a:        "/Providers/Microsoft.Authorization/RoleDefinitions/ACDD72A7-3385-48EF-BD42-F606FBA81AE7",
			b:        "/subscriptions/00000000-0000-0000-0000-000000000000/providers/microsoft.authorization/roledefinitions/acdd72a7-3385-48ef-bd42-f606fba81ae7",
			expected: true,
		},
		{
			name:     "different UUIDs",
			a:        "/providers/Microsoft.Authorization/roleDefinitions/acdd72a7-3385-48ef-bd42-f606fba81ae7",
			b:        "/providers/Microsoft.Authorization/roleDefinitions/ba92f5b4-2d11-453d-a403-e96b0029c9fe",
			expected: false,
		},
		{
			name:     "different subscriptions same UUID",
			a:        "/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Authorization/roleDefinitions/acdd72a7-3385-48ef-bd42-f606fba81ae7",
			b:        "/subscriptions/22222222-2222-2222-2222-222222222222/providers/Microsoft.Authorization/roleDefinitions/acdd72a7-3385-48ef-bd42-f606fba81ae7",
			expected: true,
		},
		{
			name:     "resource group-scoped vs unscoped same UUID",
			a:        "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myRG/providers/Microsoft.Authorization/roleDefinitions/acdd72a7-3385-48ef-bd42-f606fba81ae7",
			b:        "/providers/Microsoft.Authorization/roleDefinitions/acdd72a7-3385-48ef-bd42-f606fba81ae7",
			expected: true,
		},
		{
			name:     "resource-scoped vs subscription-scoped same UUID",
			a:        "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myRG/providers/Microsoft.Network/virtualNetworks/myVNet/providers/Microsoft.Authorization/roleDefinitions/acdd72a7-3385-48ef-bd42-f606fba81ae7",
			b:        "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/roleDefinitions/acdd72a7-3385-48ef-bd42-f606fba81ae7",
			expected: true,
		},
		{
			name:     "management group-scoped vs unscoped same UUID",
			a:        "/providers/Microsoft.Management/managementGroups/myMG/providers/Microsoft.Authorization/roleDefinitions/acdd72a7-3385-48ef-bd42-f606fba81ae7",
			b:        "/providers/Microsoft.Authorization/roleDefinitions/acdd72a7-3385-48ef-bd42-f606fba81ae7",
			expected: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if result := PimRoleDefinitionIdsMatch(tc.a, tc.b); result != tc.expected {
				t.Errorf("PimRoleDefinitionIdsMatch(%q, %q) = %v, want %v", tc.a, tc.b, result, tc.expected)
			}
		})
	}
}

func TestPimRoleDefinitionId_ID(t *testing.T) {
	cases := []struct {
		name     string
		id       PimRoleDefinitionId
		expected string
	}{
		{
			name:     "scoped",
			id:       NewPimRoleDefinitionID("/subscriptions/00000000-0000-0000-0000-000000000000", "acdd72a7-3385-48ef-bd42-f606fba81ae7"),
			expected: "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/roleDefinitions/acdd72a7-3385-48ef-bd42-f606fba81ae7",
		},
		{
			name:     "unscoped",
			id:       NewPimRoleDefinitionID("", "acdd72a7-3385-48ef-bd42-f606fba81ae7"),
			expected: "/providers/Microsoft.Authorization/roleDefinitions/acdd72a7-3385-48ef-bd42-f606fba81ae7",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if result := tc.id.ID(); result != tc.expected {
				t.Errorf("got %q, want %q", result, tc.expected)
			}
		})
	}
}
