// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

import "testing"

func TestParseManagedHSMNestedKeyWithVersionId(t *testing.T) {
	cases := []struct {
		Input       string
		Expected    ManagedHSMNestedKeyWithVersionId
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
			Input:       "https://my-hsm.managedhsm.azure.net",
			ExpectError: true,
		},
		{
			Input:       "https://my-hsm.managedhsm.azure.net/",
			ExpectError: true,
		},
		{
			Input:       "https://my-hsm.managedhsm.azure.net/keys",
			ExpectError: true,
		},
		{
			Input:       "https://my-hsm.managedhsm.azure.net/keys/bird",
			ExpectError: true,
		},
		{
			Input:       "https://my-hsm.managedhsm.azure.net/invalidNestedItemObjectType/bird/fdf067c93bbb4b22bff4d8b7a9a56217",
			ExpectError: true,
		},
		{
			Input:       "https://my-hsm.managedhsm.azure.net/keys/bird/fdf067c93bbb4b22bff4d8b7a9a56217",
			ExpectError: false,
			Expected: ManagedHSMNestedKeyWithVersionId{
				KeyName: "bird",
				BaseURI: "https://my-hsm.managedhsm.azure.net",
				Version: "fdf067c93bbb4b22bff4d8b7a9a56217",
			},
		},
		{
			Input:       "https://my-hsm.managedhsm.azure.net/keys/hello/world",
			ExpectError: false,
			Expected: ManagedHSMNestedKeyWithVersionId{
				KeyName: "hello",
				BaseURI: "https://my-hsm.managedhsm.azure.net",
				Version: "world",
			},
		},
		{
			Input:       "https://my-hsm.managedhsm.azure.net/keys/castle/1492",
			ExpectError: false,
			Expected: ManagedHSMNestedKeyWithVersionId{
				KeyName: "castle",
				BaseURI: "https://my-hsm.managedhsm.azure.net",
				Version: "1492",
			},
		},
		{
			Input:       "https://my-hsm.managedhsm.azure.net/keys/bird/fdf067c93bbb4b22bff4d8b7a9a56217/XXX",
			ExpectError: true,
		},
	}

	for _, tc := range cases {
		t.Logf("[DEBUG] Testing %q", tc.Input)

		keyID, err := ParseManagedHSMNestedKeyWithVersionID(tc.Input)
		if err != nil {
			if tc.ExpectError {
				t.Logf("[DEBUG]   --> [Received Expected Error]: %+v", err)
				continue
			}

			t.Fatalf("Got error for ID '%s': %+v", tc.Input, err)
		}

		if keyID == nil {
			t.Fatalf("Expected a SecretID to be parsed for ID '%s', got nil.", tc.Input)
		}

		if tc.Expected.BaseURI != keyID.BaseURI {
			t.Fatalf("Expected 'BaseURI' to be '%s', got '%s' for ID '%s'", tc.Expected.BaseURI, keyID.BaseURI, tc.Input)
		}

		if tc.Expected.KeyName != keyID.KeyName {
			t.Fatalf("Expected 'Version' to be '%s', got '%s' for ID '%s'", tc.Expected.KeyName, keyID.KeyName, tc.Input)
		}

		if tc.Expected.Version != keyID.Version {
			t.Fatalf("Expected 'Version' to be '%s', got '%s' for ID '%s'", tc.Expected.Version, keyID.Version, tc.Input)
		}

		if tc.Input != keyID.ID() {
			t.Fatalf("Expected 'ID()' to be '%s', got '%s'", tc.Input, keyID.ID())
		}
		t.Logf("[DEBUG]   --> [Valid Value]: %+v", tc.Input)
	}
}

func TestValidateManagedHSMNestedWithVersionId(t *testing.T) {
	cases := []struct {
		Input       string
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
			Input:       "https://my-hsm.managedhsm.azure.net",
			ExpectError: true,
		},
		{
			Input:       "https://my-hsm.managedhsm.azure.net/",
			ExpectError: true,
		},
		{
			Input:       "https://my-hsm.managedhsm.azure.net/keys",
			ExpectError: true,
		},
		{
			Input:       "https://my-hsm.managedhsm.azure.net/keys/bird",
			ExpectError: true,
		},
		{
			Input:       "https://my-hsm.managedhsm.azure.net/Keys/bird",
			ExpectError: true,
		},
		{
			Input:       "https://my-hsm.managedhsm.azure.net/invalidNestedItemObjectType/bird/fdf067c93bbb4b22bff4d8b7a9a56217",
			ExpectError: true,
		},
		{
			Input:       "https://my-hsm.managedhsm.azure.net/keys/bird/fdf067c93bbb4b22bff4d8b7a9a56217",
			ExpectError: false,
		},
		{
			Input:       "https://my-hsm.managedhsm.azure.net/keys/hello/world",
			ExpectError: false,
		},
		{
			Input:       "https://my-hsm.managedhsm.azure.net/keys/castle/1492",
			ExpectError: false,
		},
		{
			Input:       "https://my-hsm.managedhsm.azure.net/keys/bird/fdf067c93bbb4b22bff4d8b7a9a56217/XXX",
			ExpectError: true,
		},
	}

	for _, tc := range cases {
		t.Logf("[DEBUG] Testing %q", tc.Input)

		_, errors := ValidateManagedHSMNestedKeyWithVersionID(tc.Input, "test")
		if (len(errors) > 0) != tc.ExpectError {
			t.Fatalf("expect error %T, but got %T", tc.ExpectError, len(errors) > 0)
		}
	}
}
