// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = StorageTableResourceManagerId{}

func TestStorageTableResourceManagerIDFormatter(t *testing.T) {
	actual := NewStorageTableResourceManagerID("12345678-1234-9876-4563-123456789012", "resGroup1", "storageAccount1", "tableService1", "table1").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Storage/storageAccounts/storageAccount1/tableServices/tableService1/tables/table1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestStorageTableResourceManagerID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *StorageTableResourceManagerId
	}{

		{
			// empty
			Input: "",
			Error: true,
		},

		{
			// missing SubscriptionId
			Input: "/",
			Error: true,
		},

		{
			// missing value for SubscriptionId
			Input: "/subscriptions/",
			Error: true,
		},

		{
			// missing ResourceGroup
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/",
			Error: true,
		},

		{
			// missing value for ResourceGroup
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/",
			Error: true,
		},

		{
			// missing StorageAccountName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Storage/",
			Error: true,
		},

		{
			// missing value for StorageAccountName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Storage/storageAccounts/",
			Error: true,
		},

		{
			// missing TableServiceName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Storage/storageAccounts/storageAccount1/",
			Error: true,
		},

		{
			// missing value for TableServiceName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Storage/storageAccounts/storageAccount1/tableServices/",
			Error: true,
		},

		{
			// missing TableName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Storage/storageAccounts/storageAccount1/tableServices/tableService1/",
			Error: true,
		},

		{
			// missing value for TableName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Storage/storageAccounts/storageAccount1/tableServices/tableService1/tables/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Storage/storageAccounts/storageAccount1/tableServices/tableService1/tables/table1",
			Expected: &StorageTableResourceManagerId{
				SubscriptionId:     "12345678-1234-9876-4563-123456789012",
				ResourceGroup:      "resGroup1",
				StorageAccountName: "storageAccount1",
				TableServiceName:   "tableService1",
				TableName:          "table1",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/RESGROUP1/PROVIDERS/MICROSOFT.STORAGE/STORAGEACCOUNTS/STORAGEACCOUNT1/TABLESERVICES/TABLESERVICE1/TABLES/TABLE1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := StorageTableResourceManagerID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expect a value but got an error: %s", err)
		}
		if v.Error {
			t.Fatal("Expect an error but didn't get one")
		}

		if actual.SubscriptionId != v.Expected.SubscriptionId {
			t.Fatalf("Expected %q but got %q for SubscriptionId", v.Expected.SubscriptionId, actual.SubscriptionId)
		}
		if actual.ResourceGroup != v.Expected.ResourceGroup {
			t.Fatalf("Expected %q but got %q for ResourceGroup", v.Expected.ResourceGroup, actual.ResourceGroup)
		}
		if actual.StorageAccountName != v.Expected.StorageAccountName {
			t.Fatalf("Expected %q but got %q for StorageAccountName", v.Expected.StorageAccountName, actual.StorageAccountName)
		}
		if actual.TableServiceName != v.Expected.TableServiceName {
			t.Fatalf("Expected %q but got %q for TableServiceName", v.Expected.TableServiceName, actual.TableServiceName)
		}
		if actual.TableName != v.Expected.TableName {
			t.Fatalf("Expected %q but got %q for TableName", v.Expected.TableName, actual.TableName)
		}
	}
}
