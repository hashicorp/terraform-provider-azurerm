package migration

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

func TestTableStateV0ToV1(t *testing.T) {
	clouds := []azure.Environment{
		azure.ChinaCloud,
		azure.GermanCloud,
		azure.PublicCloud,
		azure.USGovernmentCloud,
	}

	for _, cloud := range clouds {
		t.Logf("[DEBUG] Testing with Cloud %q", cloud.Name)

		input := map[string]interface{}{
			"id":                   "table1",
			"name":                 "table1",
			"storage_account_name": "account1",
		}
		meta := &clients.Client{
			Account: &clients.ResourceManagerAccount{
				Environment: cloud,
			},
		}
		suffix := meta.Account.Environment.StorageEndpointSuffix

		expected := map[string]interface{}{
			"id":                   fmt.Sprintf("https://account1.table.%s/table1", suffix),
			"name":                 "table1",
			"storage_account_name": "account1",
		}

		actual, err := TableV0ToV1{}.UpgradeFunc()(context.TODO(), input, meta)
		if err != nil {
			t.Fatalf("Expected no error but got: %s", err)
		}

		if !reflect.DeepEqual(expected, actual) {
			t.Fatalf("Expected %+v. Got %+v. But expected them to be the same", expected, actual)
		}

		t.Logf("[DEBUG] Ok!")
	}
}

func TestTableStateV1ToV2(t *testing.T) {
	clouds := []azure.Environment{
		azure.ChinaCloud,
		azure.GermanCloud,
		azure.PublicCloud,
		azure.USGovernmentCloud,
	}

	for _, cloud := range clouds {
		t.Logf("[DEBUG] Testing with Cloud %q", cloud.Name)

		meta := &clients.Client{
			Account: &clients.ResourceManagerAccount{
				Environment: cloud,
			},
		}
		suffix := meta.Account.Environment.StorageEndpointSuffix

		input := map[string]interface{}{
			"id":                   fmt.Sprintf("https://account1.table.%s/table1", suffix),
			"name":                 "table1",
			"storage_account_name": "account1",
		}
		expected := map[string]interface{}{
			"id":                   fmt.Sprintf("https://account1.table.%s/Tables('table1')", suffix),
			"name":                 "table1",
			"storage_account_name": "account1",
		}

		actual, err := TableV1ToV2{}.UpgradeFunc()(context.TODO(), input, meta)
		if err != nil {
			t.Fatalf("Expected no error but got: %s", err)
		}

		if !reflect.DeepEqual(expected, actual) {
			t.Fatalf("Expected %+v. Got %+v. But expected them to be the same", expected, actual)
		}

		t.Logf("[DEBUG] Ok!")
	}
}
