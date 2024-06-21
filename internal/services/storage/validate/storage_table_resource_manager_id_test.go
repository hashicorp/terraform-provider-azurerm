// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import "testing"

func TestStorageTableResourceManagerID(t *testing.T) {
	cases := []struct {
		Input string
		Valid bool
	}{

		{
			// empty
			Input: "",
			Valid: false,
		},

		{
			// missing SubscriptionId
			Input: "/",
			Valid: false,
		},

		{
			// missing value for SubscriptionId
			Input: "/subscriptions/",
			Valid: false,
		},

		{
			// missing ResourceGroup
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/",
			Valid: false,
		},

		{
			// missing value for ResourceGroup
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/",
			Valid: false,
		},

		{
			// missing StorageAccountName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Storage/",
			Valid: false,
		},

		{
			// missing value for StorageAccountName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Storage/storageAccounts/",
			Valid: false,
		},

		{
			// missing TableServiceName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Storage/storageAccounts/storageAccount1/",
			Valid: false,
		},

		{
			// missing value for TableServiceName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Storage/storageAccounts/storageAccount1/tableServices/",
			Valid: false,
		},

		{
			// missing TableName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Storage/storageAccounts/storageAccount1/tableServices/tableService1/",
			Valid: false,
		},

		{
			// missing value for TableName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Storage/storageAccounts/storageAccount1/tableServices/tableService1/tables/",
			Valid: false,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Storage/storageAccounts/storageAccount1/tableServices/tableService1/tables/table1",
			Valid: true,
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/RESGROUP1/PROVIDERS/MICROSOFT.STORAGE/STORAGEACCOUNTS/STORAGEACCOUNT1/TABLESERVICES/TABLESERVICE1/TABLES/TABLE1",
			Valid: false,
		},
	}
	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := StorageTableResourceManagerID(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}
