// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import "testing"

func TestPostgresqlFlexibleServerVirtualEndpointID(t *testing.T) {
	cases := []struct {
		Input string
		Valid bool
	}{
		{
			Input: "",
			Valid: false,
		},
		{
			Input: "/",
			Valid: false,
		},
		{
			Input: "/subscriptions/",
			Valid: false,
		},
		{
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/",
			Valid: false,
		},
		{
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/",
			Valid: false,
		},
		{
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.DBforPostgreSQL/",
			Valid: false,
		},
		{
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.DBforPostgreSQL/flexibleServers/",
			Valid: false,
		},
		{
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.DBforPostgreSQL/flexibleServers/server1",
			Valid: false,
		},
		{
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.DBforPostgreSQL/flexibleServers/server1/virtualEndpoints/endpoint1",
			Valid: true,
		},
		{
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.DBforPostgreSQL/flexibleServers/server1/virtualendpoints/endpoint1",
			Valid: true,
		},
		{
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/RESGROUP1/PROVIDERS/MICROSOFT.DBFORPOSTGRESQL/FLEXIBLESERVERS/SERVER1/VIRTUALENDPOINTS/ENDPOINT1",
			Valid: true,
		},
	}
	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := PostgresqlFlexibleServerVirtualEndpointIDInsensitively(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}
