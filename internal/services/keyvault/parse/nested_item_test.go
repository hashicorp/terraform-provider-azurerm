// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

import "testing"

func TestNewNestedItemID(t *testing.T) {
	childType := NestedItemTypeKey
	childName := "test"
	childVersion := "testVersionString"
	cases := []struct {
		Scenario        string
		keyVaultBaseUrl string
		Expected        string
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
			keyVaultBaseUrl: "https://test.vault.azure.net",
			Expected:        "https://test.vault.azure.net/keys/test/testVersionString",
			ExpectError:     false,
		},
		{
			Scenario:        "valid, with port",
			keyVaultBaseUrl: "https://test.vault.azure.net:443",
			Expected:        "https://test.vault.azure.net/keys/test/testVersionString",
			ExpectError:     false,
		},
		{
			Scenario:        "mhsm valid, with port",
			keyVaultBaseUrl: "https://test.managedhsm.azure.net:443",
			Expected:        "https://test.managedhsm.azure.net/keys/test/testVersionString",
			ExpectError:     false,
		},
	}
	for _, tc := range cases {
		t.Logf("[DEBUG] Testing %q", tc.keyVaultBaseUrl)

		id, err := NewNestedItemID(tc.keyVaultBaseUrl, childType, childName, childVersion)
		if err != nil {
			if !tc.ExpectError {
				t.Fatalf("Got error for New Resource ID '%s': %+v", tc.keyVaultBaseUrl, err)
				return
			}
			t.Logf("[DEBUG]   --> [Received Expected Error]: %+v", err)
			continue
		}
		if id.ID() != tc.Expected {
			t.Fatalf("Expected id for %q to be %q, got %q", tc.keyVaultBaseUrl, tc.Expected, id)
		}
		t.Logf("[DEBUG]   --> [Valid Value]: %+v", tc.keyVaultBaseUrl)
	}
}

func TestParseNestedItemID(t *testing.T) {
	cases := []struct {
		Input       string
		Expected    NestedItemId
		ExpectError bool
	}{
		{
			Input:       "",
			ExpectError: true,
		},
		{
			Input:       "https",
			ExpectError: true,
		},
		{
			Input:       "https://",
			ExpectError: true,
		},
		{
			Input:       "https://my-keyvault.vault.azure.net",
			ExpectError: true,
		},
		{
			Input:       "https://my-keyvault.vault.azure.net/",
			ExpectError: true,
		},
		{
			Input:       "https://my-keyvault.vault.azure.net/secrets",
			ExpectError: true,
		},
		{
			Input:       "https://my-keyvault.vault.azure.net/secrets/bird",
			ExpectError: true,
		},
		{
			Input:       "https://my-keyvault.vault.azure.net/invalidNestedItemObjectType/bird/fdf067c93bbb4b22bff4d8b7a9a56217",
			ExpectError: true,
		},
		{
			Input:       "https://my-keyvault.vault.azure.net/secrets/bird/fdf067c93bbb4b22bff4d8b7a9a56217",
			ExpectError: false,
			Expected: NestedItemId{
				Name:            "bird",
				KeyVaultBaseUrl: "https://my-keyvault.vault.azure.net/",
				Version:         "fdf067c93bbb4b22bff4d8b7a9a56217",
			},
		},
		{
			Input:       "https://my-keyvault.vault.azure.net/certificates/hello/world",
			ExpectError: false,
			Expected: NestedItemId{
				Name:            "hello",
				KeyVaultBaseUrl: "https://my-keyvault.vault.azure.net/",
				Version:         "world",
			},
		},
		{
			Input:       "https://my-keyvault.vault.azure.net/keys/castle/1492",
			ExpectError: false,
			Expected: NestedItemId{
				Name:            "castle",
				KeyVaultBaseUrl: "https://my-keyvault.vault.azure.net/",
				Version:         "1492",
			},
		},
		{
			Input:       "https://my-keyvault.managedhsm.azure.net/keys/castle/1492",
			ExpectError: false,
			Expected: NestedItemId{
				Name:            "castle",
				KeyVaultBaseUrl: "https://my-keyvault.managedhsm.azure.net/",
				Version:         "1492",
				NestedItemType:  NestedItemTypeKey,
			},
		},
		{
			Input:       "https://my-keyvault.vault.azure.net/secrets/bird/fdf067c93bbb4b22bff4d8b7a9a56217/XXX",
			ExpectError: true,
		},
	}

	for _, tc := range cases {
		t.Logf("[DEBUG] Testing %q", tc.Input)

		secretId, err := ParseNestedItemID(tc.Input)
		if err != nil {
			if tc.ExpectError {
				t.Logf("[DEBUG]   --> [Received Expected Error]: %+v", err)
				continue
			}

			t.Fatalf("Got error for ID '%s': %+v", tc.Input, err)
		}

		if secretId == nil {
			t.Fatalf("Expected a SecretID to be parsed for ID '%s', got nil.", tc.Input)
		}

		if tc.Expected.KeyVaultBaseUrl != secretId.KeyVaultBaseUrl {
			t.Fatalf("Expected 'KeyVaultBaseUrl' to be '%s', got '%s' for ID '%s'", tc.Expected.KeyVaultBaseUrl, secretId.KeyVaultBaseUrl, tc.Input)
		}

		if tc.Expected.Name != secretId.Name {
			t.Fatalf("Expected 'Version' to be '%s', got '%s' for ID '%s'", tc.Expected.Name, secretId.Name, tc.Input)
		}

		if tc.Expected.Version != secretId.Version {
			t.Fatalf("Expected 'Version' to be '%s', got '%s' for ID '%s'", tc.Expected.Version, secretId.Version, tc.Input)
		}

		if tc.Input != secretId.ID() {
			t.Fatalf("Expected 'ID()' to be '%s', got '%s'", tc.Input, secretId.ID())
		}
		t.Logf("[DEBUG]   --> [Valid Value]: %+v", tc.Input)
	}
}

func TestParseOptionallyVersionedNestedItemID(t *testing.T) {
	cases := []struct {
		Input       string
		Expected    NestedItemId
		ExpectError bool
	}{
		{
			Input:       "",
			ExpectError: true,
		},
		{
			Input:       "https://my-keyvault.vault.azure.net/secrets",
			ExpectError: true,
		},
		{
			Input:       "https://my-keyvault.vault.azure.net/invalidNestedItemObjectType/hello/world",
			ExpectError: true,
		},
		{
			Input:       "https://my-keyvault.vault.azure.net/secrets/bird",
			ExpectError: false,
			Expected: NestedItemId{
				Name:            "bird",
				KeyVaultBaseUrl: "https://my-keyvault.vault.azure.net/",
				Version:         "",
			},
		},
		{
			Input:       "https://my-keyvault.vault.azure.net/secrets/bird/fdf067c93bbb4b22bff4d8b7a9a56217",
			ExpectError: false,
			Expected: NestedItemId{
				Name:            "bird",
				KeyVaultBaseUrl: "https://my-keyvault.vault.azure.net/",
				Version:         "fdf067c93bbb4b22bff4d8b7a9a56217",
			},
		},
		{
			Input:       "https://my-keyvault.vault.azure.net/certificates/hello/world",
			ExpectError: false,
			Expected: NestedItemId{
				Name:            "hello",
				KeyVaultBaseUrl: "https://my-keyvault.vault.azure.net/",
				Version:         "world",
			},
		},
		{
			Input:       "https://my-keyvault.vault.azure.net/keys/castle/1492",
			ExpectError: false,
			Expected: NestedItemId{
				Name:            "castle",
				KeyVaultBaseUrl: "https://my-keyvault.vault.azure.net/",
				Version:         "1492",
			},
		},
		{
			Input:       "https://my-keyvault.managedhsm.azure.net/keys/castle/1492",
			ExpectError: false,
			Expected: NestedItemId{
				Name:            "castle",
				KeyVaultBaseUrl: "https://my-keyvault.managedhsm.azure.net/",
				Version:         "1492",
				NestedItemType:  NestedItemTypeKey,
			},
		},
		{
			Input:       "https://my-keyvault.vault.azure.net/secrets/bird/fdf067c93bbb4b22bff4d8b7a9a56217/XXX",
			ExpectError: true,
		},
	}

	for _, tc := range cases {
		t.Logf("[DEBUG] Testing %q", tc.Input)

		secretId, err := ParseOptionallyVersionedNestedItemID(tc.Input)
		if err != nil {
			if tc.ExpectError {
				t.Logf("[DEBUG]   --> [Received Expected Error]: %+v", err)
				continue
			}

			t.Fatalf("Got error for ID '%s': %+v", tc.Input, err)
		}

		if secretId == nil {
			t.Fatalf("Expected a SecretID to be parsed for ID '%s', got nil.", tc.Input)
		}

		if tc.Expected.KeyVaultBaseUrl != secretId.KeyVaultBaseUrl {
			t.Fatalf("Expected 'KeyVaultBaseUrl' to be '%s', got '%s' for ID '%s'", tc.Expected.KeyVaultBaseUrl, secretId.KeyVaultBaseUrl, tc.Input)
		}

		if tc.Expected.Name != secretId.Name {
			t.Fatalf("Expected 'Version' to be '%s', got '%s' for ID '%s'", tc.Expected.Name, secretId.Name, tc.Input)
		}

		if tc.Expected.Version != secretId.Version {
			t.Fatalf("Expected 'Version' to be '%s', got '%s' for ID '%s'", tc.Expected.Version, secretId.Version, tc.Input)
		}

		if tc.Input != secretId.ID() {
			t.Fatalf("Expected 'ID()' to be '%s', got '%s'", tc.Input, secretId.ID())
		}
		t.Logf("[DEBUG]   --> [Valid Value]: %+v", tc.Input)
	}
}
