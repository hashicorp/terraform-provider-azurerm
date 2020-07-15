package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceArmStorageAccountTableSas_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_storage_account_table_sas", "test")
	utcNow := time.Now().UTC()
	startDate := utcNow.Format(time.RFC3339)
	endDate := utcNow.Add(time.Hour * 24).Format(time.RFC3339)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMStorageAccountTableSas_basic(data, startDate, endDate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "https_only", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "start", startDate),
					resource.TestCheckResourceAttr(data.ResourceName, "expiry", endDate),
					resource.TestCheckResourceAttr(data.ResourceName, "ip_address", "168.1.5.65"),
					resource.TestCheckResourceAttr(data.ResourceName, "permissions.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "permissions.0.read", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "permissions.0.add", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "permissions.0.update", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "permissions.0.delete", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "start_partition_key", "Coho Winery"),
					resource.TestCheckResourceAttr(data.ResourceName, "end_partition_key", "Auburn"),
					resource.TestCheckResourceAttr(data.ResourceName, "start_row_key", "Coho Winery"),
					resource.TestCheckResourceAttr(data.ResourceName, "end_row_key", "Portland"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "sas"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMStorageAccountTableSas_basic(data acceptance.TestData, startDate string, endDate string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "rg" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_storage_account" "storage" {
  name                = "acctestsads%s"
  resource_group_name = azurerm_resource_group.rg.name

  location                 = azurerm_resource_group.rg.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_table" "table" {
  name                  = "sastest"
  storage_account_name  = azurerm_storage_account.storage.name
}

data "azurerm_storage_account_table_sas" "test" {
  connection_string = azurerm_storage_account.storage.primary_connection_string
  table_name        = azurerm_storage_table.table.name
  https_only        = true

  ip_address = "168.1.5.65"

  start  = "%s"
  expiry = "%s"

  permissions {
    read   = true
    add    = true
    update = false
    delete = true
  }

  start_partition_key = "Coho Winery"
  end_partition_key   = "Auburn"
  start_row_key       = "Coho Winery"
  end_row_key         = "Portland"
}

`, data.RandomInteger, data.Locations.Primary, data.RandomString, startDate, endDate)
}

func TestAccDataSourceArmStorageAccountTableSas_permissionsString(t *testing.T) {
	testCases := []struct {
		input    map[string]interface{}
		expected string
	}{
		{map[string]interface{}{"read": true}, "r"},
		{map[string]interface{}{"add": true}, "a"},
		{map[string]interface{}{"update": true}, "u"},
		{map[string]interface{}{"delete": true}, "d"},
		{map[string]interface{}{"add": true, "update": true, "read": true, "delete": true}, "raud"},
	}

	for _, test := range testCases {
		result := storage.BuildTablePermissionsString(test.input)
		if test.expected != result {
			t.Fatalf("Failed to build resource type string: expected: %s, result: %s", test.expected, result)
		}
	}
}
