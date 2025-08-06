// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

import "testing"

func TestManagedHSMDataPlaneRoleDefinitionID(t *testing.T) {
	cases := []struct {
		input    string
		expected *ManagedHSMDataPlaneRoleDefinitionId
	}{
		{
			input:    "",
			expected: nil,
		},
		{
			// missing the path
			input:    "https://my-hsm.managedhsm.azure.net/",
			expected: nil,
		},
		{
			// scope but no middle component or role definition name
			input:    "https://my-hsm.managedhsm.azure.net///test",
			expected: nil,
		},
		{
			// missing role definition name
			input:    "https://my-hsm.managedhsm.azure.net////providers/Microsoft.Authorization/roleDefinitions/",
			expected: nil,
		},
		{
			// Key Vault URIs are not valid
			input:    "https://my-hsm.vault.azure.net///providers/Microsoft.Authorization/roleDefinitions/test",
			expected: nil,
		},
		{
			// scope is missing
			input:    "https://my-hsm.managedhsm.azure.net//providers/Microsoft.Authorization/roleDefinitions/test",
			expected: nil,
		},
		{
			input: "https://my-hsm.managedhsm.azure.net///providers/Microsoft.Authorization/roleDefinitions/test",
			expected: &ManagedHSMDataPlaneRoleDefinitionId{
				ManagedHSMName:     "my-hsm",
				DomainSuffix:       "managedhsm.azure.net",
				Scope:              "/",
				RoleDefinitionName: "test",
			},
		},
		{
			input: "https://my-hsm.managedhsm.azure.net//keys/providers/Microsoft.Authorization/roleDefinitions/1492",
			expected: &ManagedHSMDataPlaneRoleDefinitionId{
				ManagedHSMName:     "my-hsm",
				DomainSuffix:       "managedhsm.azure.net",
				Scope:              "/keys",
				RoleDefinitionName: "1492",
			},
		},
		{
			input: "https://my-hsm.managedhsm.azure.net//keys/abc123/providers/Microsoft.Authorization/roleDefinitions/1492",
			expected: &ManagedHSMDataPlaneRoleDefinitionId{
				ManagedHSMName:     "my-hsm",
				DomainSuffix:       "managedhsm.azure.net",
				Scope:              "/keys/abc123",
				RoleDefinitionName: "1492",
			},
		},
		{
			input: "https://my-hsm.managedhsm.azure.cn///providers/Microsoft.Authorization/roleDefinitions/test",
			expected: &ManagedHSMDataPlaneRoleDefinitionId{
				ManagedHSMName:     "my-hsm",
				DomainSuffix:       "managedhsm.azure.cn",
				Scope:              "/",
				RoleDefinitionName: "test",
			},
		},
		{
			input: "https://my-hsm.managedhsm.usgovcloudapi.net//keys/providers/Microsoft.Authorization/roleDefinitions/1492",
			expected: &ManagedHSMDataPlaneRoleDefinitionId{
				ManagedHSMName:     "my-hsm",
				DomainSuffix:       "managedhsm.usgovcloudapi.net",
				Scope:              "/keys",
				RoleDefinitionName: "1492",
			},
		},
		{
			// extra suffix at the end
			input:    "https://my-hsm.managedhsm.azure.net//keys//providers/Microsoft.Authorization/roleDefinitions/1492/suffix",
			expected: nil,
		},
		{
			// valid format but missing scope
			input:    "https://my-hsm.managedhsm.azure.net/providers/Microsoft.Authorization/roleDefinitions/000-000",
			expected: nil,
		},
	}

	for _, test := range cases {
		t.Logf("Testing %q..", test.input)
		actual, err := ManagedHSMDataPlaneRoleDefinitionID(test.input, nil)
		if err != nil {
			if test.expected == nil {
				continue
			}
			t.Fatalf("unexpected error: %+v", err)
		}
		if test.expected == nil {
			if actual == nil {
				continue
			}
			t.Fatalf("expected nothing but got %+v", actual)
		}
		if actual == nil {
			t.Fatalf("expected %+v but got nil", test.expected)
		}
		if test.expected.ManagedHSMName != actual.ManagedHSMName {
			t.Fatalf("expected ManagedHSMName to be %q but got %q", test.expected.ManagedHSMName, actual.ManagedHSMName)
		}
		if test.expected.DomainSuffix != actual.DomainSuffix {
			t.Fatalf("expected DomainSuffix to be %q but got %q", test.expected.DomainSuffix, actual.DomainSuffix)
		}
		if test.expected.Scope != actual.Scope {
			t.Fatalf("expected Scope to be %q but got %q", test.expected.Scope, actual.Scope)
		}
		if test.expected.RoleDefinitionName != actual.RoleDefinitionName {
			t.Fatalf("expected RoleDefinitionName to be %q but got %q", test.expected.RoleDefinitionName, actual.RoleDefinitionName)
		}
	}
}
