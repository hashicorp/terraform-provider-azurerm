// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"strings"
	"testing"
)

func TestParseManagedHSMNestedKeyVersionlessId(t *testing.T) {
	cases := []struct {
		Input        string
		Insenstively bool
		Expected     ManagedHSMNestedKeyWithVersionId
		ExpectError  bool
	}{
		{
			Input:       "",
			ExpectError: true,
		},
		{
			Input:       "https://my-hsm.managedhsm.azure.net/keys",
			ExpectError: true,
		},
		{
			Input:       "https://my-hsm.managedhsm.azure.net/invalidNestedItemObjectType/hello/world",
			ExpectError: true,
		},
		{
			Input:       "https://my-hsm.managedhsm.azure.net/Keys/hello/world",
			ExpectError: true,
		},
		{
			Input:        "https://my-hsm.managedhsm.azure.net/Keys/hello",
			Insenstively: true,
			Expected: ManagedHSMNestedKeyWithVersionId{
				KeyName: "hello",
				BaseURI: "https://my-hsm.managedhsm.azure.net",
			},
		},
		{
			Input:       "https://my-hsm.managedhsm.azure.net/keys/bird/version",
			ExpectError: true,
		},
		{
			Input: "https://my-hsm.managedhsm.azure.net/keys/bird",
			Expected: ManagedHSMNestedKeyWithVersionId{
				KeyName: "bird",
				BaseURI: "https://my-hsm.managedhsm.azure.net",
			},
		},
		{
			Input:        "https://my-hsm.managedhsm.azure.net/Keys/castle",
			Insenstively: true,
			Expected: ManagedHSMNestedKeyWithVersionId{
				KeyName: "castle",
				BaseURI: "https://my-hsm.managedhsm.azure.net",
			},
		},
	}

	for _, tc := range cases {
		t.Logf("[DEBUG] Testing %q", tc.Input)

		var keyId *ManagedHSMNestedkeyVersionless
		var err error
		if tc.Insenstively {
			keyId, err = ParseManagedHSMNestedkeyVersionlessInsensitively(tc.Input)
		} else {
			keyId, err = ParseManagedHSMNestedkeyVersionless(tc.Input)
		}
		if err != nil {
			if tc.ExpectError {
				t.Logf("[DEBUG]   --> [Received Expected Error]: %+v", err)
				continue
			}

			t.Fatalf("Got error for ID '%s': %+v", tc.Input, err)
		}

		if keyId == nil {
			t.Fatalf("Expected a SecretID to be parsed for ID '%s', got nil.", tc.Input)
		}

		if tc.Expected.BaseURI != keyId.BaseURI {
			t.Fatalf("Expected 'BaseURI' to be '%s', got '%s' for ID '%s'", tc.Expected.BaseURI, keyId.BaseURI, tc.Input)
		}

		if tc.Expected.KeyName != keyId.KeyName {
			t.Fatalf("Expected 'Version' to be '%s', got '%s' for ID '%s'", tc.Expected.KeyName, keyId.KeyName, tc.Input)
		}

		if !(tc.Input == keyId.ID() || (tc.Insenstively && strings.EqualFold(tc.Input, keyId.ID()))) {
			t.Fatalf("Expected 'ID()' to be '%s', got '%s'", tc.Input, keyId.ID())
		}
		t.Logf("[DEBUG]   --> [Valid Value]: %+v", tc.Input)
	}
}

func TestValidateManagedHSMNestedVersionlessId(t *testing.T) {
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
			Input:       "https://my-hsm.managedhsm.azure.net/Keys/bird",
			ExpectError: true,
		},
		{
			Input:       "https://my-hsm.managedhsm.azure.net/invalidNestedItemObjectType/bird/fdf067c93bbb4b22bff4d8b7a9a56217",
			ExpectError: true,
		},
		{
			Input:       "https://my-hsm.managedhsm.azure.net/keys/bird/fdf067c93bbb4b22bff4d8b7a9a56217",
			ExpectError: true,
		},
		{
			Input:       "https://my-hsm.managedhsm.azure.net/keys/bird/fdf067c93bbb4b22bff4d8b7a9a56217/XXX",
			ExpectError: true,
		},
		{
			Input:       "https://my-hsm.managedhsm.azure.net/keys/bird",
			ExpectError: false,
		},
	}

	for idx, tc := range cases {
		t.Logf("[DEBUG] Testing %q", tc.Input)

		_, errors := ValidateManagedHSMNestedKeyVersionlessID(tc.Input, "test")
		if (len(errors) > 0) != tc.ExpectError {
			t.Fatalf("case %d: expect error %t, but got %t", idx, tc.ExpectError, len(errors) > 0)
		}
	}
}
