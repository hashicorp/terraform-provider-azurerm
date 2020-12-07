package storage_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

func TestAzureRMStorageShareMigrateStateV0ToV1(t *testing.T) {
	clouds := []azure.Environment{
		azure.ChinaCloud,
		azure.GermanCloud,
		azure.PublicCloud,
		azure.USGovernmentCloud,
	}

	for _, cloud := range clouds {
		t.Logf("[DEBUG] Testing with Cloud %q", cloud.Name)

		input := map[string]interface{}{
			"id":                   "share1",
			"name":                 "share1",
			"resource_group_name":  "group1",
			"storage_account_name": "account1",
			"quota":                5120,
		}
		meta := &clients.Client{
			Account: &clients.ResourceManagerAccount{
				Environment: cloud,
			},
		}
		expected := map[string]interface{}{
			"id":                   "share1/group1/account1",
			"name":                 "share1",
			"resource_group_name":  "group1",
			"storage_account_name": "account1",
			"quota":                5120,
		}

		actual, err := ResourceStorageShareStateUpgradeV0ToV1(input, meta)
		if err != nil {
			t.Fatalf("Expected no error but got: %s", err)
		}

		if !reflect.DeepEqual(expected, actual) {
			t.Fatalf("Expected %+v. Got %+v. But expected them to be the same", expected, actual)
		}

		t.Logf("[DEBUG] Ok!")
	}
}

func TestAzureRMStorageShareMigrateStateV1ToV2(t *testing.T) {
	clouds := []azure.Environment{
		azure.ChinaCloud,
		azure.GermanCloud,
		azure.PublicCloud,
		azure.USGovernmentCloud,
	}

	for _, cloud := range clouds {
		t.Logf("[DEBUG] Testing with Cloud %q", cloud.Name)

		input := map[string]interface{}{
			"id":                   "share1/group1/account1",
			"name":                 "share1",
			"resource_group_name":  "group1",
			"storage_account_name": "account1",
			"quota":                5120,
		}
		meta := &clients.Client{
			Account: &clients.ResourceManagerAccount{
				Environment: cloud,
			},
		}
		expected := map[string]interface{}{
			"id":                   fmt.Sprintf("https://account1.file.%s/share1", cloud.StorageEndpointSuffix),
			"name":                 "share1",
			"resource_group_name":  "group1",
			"storage_account_name": "account1",
			"quota":                5120,
		}

		actual, err := ResourceStorageShareStateUpgradeV1ToV2(input, meta)
		if err != nil {
			t.Fatalf("Expected no error but got: %s", err)
		}

		if !reflect.DeepEqual(expected, actual) {
			t.Fatalf("Expected %+v. Got %+v. But expected them to be the same", expected, actual)
		}

		t.Logf("[DEBUG] Ok!")
	}
}
