// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

/*

import (
	"fmt"
	"testing"
)

func TestParseManagedHSMRoleDefinitionID(t *testing.T) {
	typ := "RoleAssignment"
	cases := []struct {
		Input       string
		Expected    ManagedHSMRoleDefinitionId
		ExpectError bool
	}{
		{
			Input:       "",
			ExpectError: true,
		},
		{
			Input:       fmt.Sprintf("https://my-keyvault.managedhsm.azure.net///%s/test", typ),
			ExpectError: true,
			Expected: ManagedHSMRoleDefinitionId{
				Name:         "test",
				VaultBaseUrl: "https://my-keyvault.managedhsm.azure.net/",
				Scope:        "/",
			},
		},
		{
			Input:       fmt.Sprintf("https://my-keyvault.managedhsm.azure.net///%s/bird", typ),
			ExpectError: true,
			Expected: ManagedHSMRoleDefinitionId{
				Name:         "bird",
				VaultBaseUrl: "https://my-keyvault.managedhsm.azure.net/",
				Scope:        "/",
			},
		},
		{
			Input:       fmt.Sprintf("https://my-keyvault.managedhsm.azure.net///%s/bird", typ),
			ExpectError: false,
			Expected: ManagedHSMRoleDefinitionId{
				Name:         "bird",
				VaultBaseUrl: "https://my-keyvault.managedhsm.azure.net/",
				Scope:        "/",
			},
		},
		{
			Input:       fmt.Sprintf("https://my-keyvault.managedhsm.azure.net//keys/%s/world", typ),
			ExpectError: false,
			Expected: ManagedHSMRoleDefinitionId{
				Name:         "world",
				VaultBaseUrl: "https://my-keyvault.managedhsm.azure.net/",
				Scope:        "/keys",
			},
		},
		{
			Input:       fmt.Sprintf("https://my-keyvault.managedhsm.azure.net//keys/%s/fdf067c93bbb4b22bff4d8b7a9a56217", typ),
			ExpectError: true,
			Expected: ManagedHSMRoleDefinitionId{
				Name:         "fdf067c93bbb4b22bff4d8b7a9a56217",
				VaultBaseUrl: "https://my-keyvault.managedhsm.azure.net/",
				Scope:        "/keys",
			},
		},
		{
			Input:       "https://kvhsm23030816100222.managedhsm.azure.net///RoleAssignment/862d4d5e-bf01-11ed-a49d-00155d61ee9e",
			ExpectError: true,
			Expected: ManagedHSMRoleDefinitionId{
				Name:         "862d4d5e-bf01-11ed-a49d-00155d61ee9e",
				VaultBaseUrl: "https://kvhsm23030816100222.managedhsm.azure.net/",
				Scope:        "/",
			},
		},
	}

	for idx, tc := range cases {
		roleId, err := ManagedHSMRoleAssignmentID(tc.Input)
		if err != nil {
			if tc.ExpectError {
				continue
			}

			t.Fatalf("Got error for ID '%s': %+v", tc.Input, err)
		}

		if roleId == nil {
			t.Fatalf("Expected a SecretID to be parsed for ID '%s', got nil.", tc.Input)
		}

		if tc.Expected.VaultBaseUrl != roleId.VaultBaseUrl {
			t.Fatalf("Expected %d 'KeyVaultBaseUrl' to be '%s', got '%s' for ID '%s'", idx, tc.Expected.VaultBaseUrl, roleId.VaultBaseUrl, tc.Input)
		}

		if tc.Expected.Name != roleId.Name {
			t.Fatalf("Expected 'Name' to be '%s', got '%s' for ID '%s'", tc.Expected.Name, roleId.Name, tc.Input)
		}

		if tc.Expected.Scope != roleId.Scope {
			t.Fatalf("Expected 'Scope' to be '%s', got '%s' for ID '%s'", tc.Expected.Scope, roleId.Scope, tc.Input)
		}

		if tc.Input != roleId.ID() {
			t.Fatalf("Expected 'ID()' to be '%s', got '%s'", tc.Input, roleId.ID())
		}
	}
}
*/
