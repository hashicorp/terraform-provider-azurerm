// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

import "testing"

func TestNewNestedItemID(t *testing.T) {
	childType := NestedItemTypeOrganization
	childId := "test"
	cases := []struct {
		Scenario          string
		IotCentralBaseUrl string
		Expected          string
		ExpectError       bool
	}{
		{
			Scenario:          "empty values",
			IotCentralBaseUrl: "",
			Expected:          "",
			ExpectError:       true,
		},
		{
			Scenario:          "valid, no port",
			IotCentralBaseUrl: "https://subdomain.baseDomain",
			Expected:          "https://subdomain.baseDomain/api/organizations/test",
			ExpectError:       false,
		},
		{
			Scenario:          "valid, with port",
			IotCentralBaseUrl: "https://subdomain.baseDomain:443",
			Expected:          "https://subdomain.baseDomain/api/organizations/test",
			ExpectError:       false,
		},
	}
	for _, tc := range cases {
		t.Logf("[DEBUG] Testing %q", tc.IotCentralBaseUrl)

		id, err := NewNestedItemID(tc.IotCentralBaseUrl, childType, childId)
		if err != nil {
			if !tc.ExpectError {
				t.Fatalf("Got error for New Resource ID '%s': %+v", tc.IotCentralBaseUrl, err)
				return
			}
			t.Logf("[DEBUG]   --> [Received Expected Error]: %+v", err)
			continue
		}
		if id.ID() != tc.Expected {
			t.Fatalf("Expected id for %q to be %q, got %q", tc.IotCentralBaseUrl, tc.Expected, id)
		}
		t.Logf("[DEBUG]   --> [Valid Value]: %+v", tc.IotCentralBaseUrl)
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
			Input:       "https://subdomain.baseDomain",
			ExpectError: true,
		},
		{
			Input:       "https://subdomain.baseDomain/",
			ExpectError: true,
		},
		{
			Input:       "https://subdomain.baseDomain/api",
			ExpectError: true,
		},
		{
			Input:       "https://subdomain.baseDomain/api/organizations",
			ExpectError: true,
		},
		{
			Input:       "https://subdomain.baseDomain/api/domains/fdf067c93bbb4b22bff4d8b7a9a56217",
			ExpectError: true,
		},
		{
			Input:       "https://subdomain.baseDomain/api/organizations/fdf067c93bbb4b22bff4d8b7a9a56217",
			ExpectError: false,
			Expected: NestedItemId{
				Id:                "fdf067c93bbb4b22bff4d8b7a9a56217",
				IotcentralBaseUrl: "https://subdomain.baseDomain/",
				SubDomain:         "subdomain",
			},
		},
		{
			Input:       "https://subdomain.baseDomain/api/users/fdf067c93bbb4b22bff4d8b7a9a56217",
			ExpectError: false,
			Expected: NestedItemId{
				Id:                "fdf067c93bbb4b22bff4d8b7a9a56217",
				IotcentralBaseUrl: "https://subdomain.baseDomain/",
				SubDomain:         "subdomain",
			},
		},
	}

	for _, tc := range cases {
		t.Logf("[DEBUG] Testing %q", tc.Input)

		orgId, err := ParseNestedItemID(tc.Input)
		if err != nil {
			if tc.ExpectError {
				t.Logf("[DEBUG]   --> [Received Expected Error]: %+v", err)
				continue
			}

			t.Fatalf("Got error for ID '%s': %+v", tc.Input, err)
		}

		if orgId == nil {
			t.Fatalf("Expected a SecretID to be parsed for ID '%s', got nil.", tc.Input)
		}

		if tc.Expected.IotcentralBaseUrl != orgId.IotcentralBaseUrl {
			t.Fatalf("Expected 'IotcentralBaseUrl' to be '%s', got '%s' for ID '%s'", tc.Expected.IotcentralBaseUrl, orgId.IotcentralBaseUrl, tc.Input)
		}

		if tc.Expected.Id != orgId.Id {
			t.Fatalf("Expected 'Id' to be '%s', got '%s' for ID '%s'", tc.Expected.Id, orgId.Id, tc.Input)
		}

		if tc.Input != orgId.ID() {
			t.Fatalf("Expected 'ID()' to be '%s', got '%s'", tc.Input, orgId.ID())
		}
		t.Logf("[DEBUG]   --> [Valid Value]: %+v", tc.Input)
	}
}
