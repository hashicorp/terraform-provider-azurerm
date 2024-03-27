// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"fmt"
	"testing"
)

func TestNewMHSMNestedItemID(t *testing.T) {
	mhsmType := RoleDefinitionType
	cases := []struct {
		Scenario        string
		keyVaultBaseUrl string
		Expected        string
		Scope           string
		Name            string
		ExpectError     bool
	}{
		{
			Scenario:        "empty values",
			keyVaultBaseUrl: "",
			Expected:        "",
			ExpectError:     true,
		},
		{
			Scenario:        "valid, no port",
			keyVaultBaseUrl: "https://test.managedhsm.azure.net",
			Scope:           "/",
			Name:            "test",
			Expected:        fmt.Sprintf("https://test.managedhsm.azure.net///%s/test", mhsmType),
			ExpectError:     false,
		},
		{
			Scenario:        "valid, with port",
			keyVaultBaseUrl: "https://test.managedhsm.azure.net:443",
			Scope:           "/keys",
			Name:            "test",
			Expected:        fmt.Sprintf("https://test.managedhsm.azure.net//keys/%s/test", mhsmType),
			ExpectError:     false,
		},
	}
	for idx, tc := range cases {
		id, err := NewRoleNestedItemID(tc.keyVaultBaseUrl, tc.Scope, mhsmType, tc.Name)
		if err != nil {
			if !tc.ExpectError {
				t.Fatalf("Got error for New Resource ID '%s': %+v", tc.keyVaultBaseUrl, err)
				return
			}
			continue
		}
		if id.ID() != tc.Expected {
			t.Fatalf("Expected %d id for %q to be %q, got %q", idx, tc.keyVaultBaseUrl, tc.Expected, id)
		}
	}
}

func TestParseMHSMNestedItemID(t *testing.T) {
	typ := RoleDefinitionType
	cases := []struct {
		Input       string
		Expected    RoleNestedItemId
		ExpectError bool
	}{
		{
			Input:       "",
			ExpectError: true,
		},
		{
			Input:       fmt.Sprintf("https://my-keyvault.managedhsm.azure.net///%s/test", typ),
			ExpectError: true,
			Expected: RoleNestedItemId{
				Name:         "test",
				VaultBaseUrl: "https://my-keyvault.managedhsm.azure.net/",
				Scope:        "/",
			},
		},
		{
			Input:       fmt.Sprintf("https://my-keyvault.managedhsm.azure.net///%s/bird", typ),
			ExpectError: true,
			Expected: RoleNestedItemId{
				Name:         "bird",
				VaultBaseUrl: "https://my-keyvault.managedhsm.azure.net/",
				Scope:        "/",
			},
		},
		{
			Input:       fmt.Sprintf("https://my-keyvault.managedhsm.azure.net///%s/bird", typ),
			ExpectError: false,
			Expected: RoleNestedItemId{
				Name:         "bird",
				VaultBaseUrl: "https://my-keyvault.managedhsm.azure.net/",
				Scope:        "/",
			},
		},
		{
			Input:       fmt.Sprintf("https://my-keyvault.managedhsm.azure.net//keys/%s/world", typ),
			ExpectError: false,
			Expected: RoleNestedItemId{
				Name:         "world",
				VaultBaseUrl: "https://my-keyvault.managedhsm.azure.net/",
				Scope:        "/keys",
			},
		},
		{
			Input:       fmt.Sprintf("https://my-keyvault.managedhsm.azure.net//keys/%s/fdf067c93bbb4b22bff4d8b7a9a56217", typ),
			ExpectError: true,
			Expected: RoleNestedItemId{
				Name:         "fdf067c93bbb4b22bff4d8b7a9a56217",
				VaultBaseUrl: "https://my-keyvault.managedhsm.azure.net/",
				Scope:        "/keys",
			},
		},
		{
			Input:       "https://kvhsm23030816100222.managedhsm.azure.net///RoleDefinition/862d4d5e-bf01-11ed-a49d-00155d61ee9e",
			ExpectError: true,
			Expected: RoleNestedItemId{
				Name:         "862d4d5e-bf01-11ed-a49d-00155d61ee9e",
				VaultBaseUrl: "https://kvhsm23030816100222.managedhsm.azure.net/",
				Scope:        "/",
			},
		},
	}

	for idx, tc := range cases {
		secretId, err := RoleNestedItemID(tc.Input)
		if err != nil {
			if tc.ExpectError {
				continue
			}

			t.Fatalf("Got error for ID '%s': %+v", tc.Input, err)
		}

		if secretId == nil {
			t.Fatalf("Expected a SecretID to be parsed for ID '%s', got nil.", tc.Input)
		}

		if tc.Expected.VaultBaseUrl != secretId.VaultBaseUrl {
			t.Fatalf("Expected %d 'KeyVaultBaseUrl' to be '%s', got '%s' for ID '%s'", idx, tc.Expected.VaultBaseUrl, secretId.VaultBaseUrl, tc.Input)
		}

		if tc.Expected.Name != secretId.Name {
			t.Fatalf("Expected 'Name' to be '%s', got '%s' for ID '%s'", tc.Expected.Name, secretId.Name, tc.Input)
		}

		if tc.Expected.Scope != secretId.Scope {
			t.Fatalf("Expected 'Scope' to be '%s', got '%s' for ID '%s'", tc.Expected.Scope, secretId.Scope, tc.Input)
		}

		if tc.Input != secretId.ID() {
			t.Fatalf("Expected 'ID()' to be '%s', got '%s'", tc.Input, secretId.ID())
		}
	}
}
