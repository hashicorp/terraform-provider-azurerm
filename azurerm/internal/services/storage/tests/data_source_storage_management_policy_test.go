package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMStorageManagementPolicy_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_storage_account", "test")
	config := testAccDataSourceAzureRMStorageManagementPolicy_basic(data)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "rule.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.name", "rule1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.filters.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.filters.0.prefix_match.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.filters.0.prefix_match.3439697764", "container1/prefix1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.filters.0.blob_types.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.filters.0.blob_types.1068358194", "blockBlob"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.actions.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.actions.0.base_blob.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.actions.0.base_blob.0.tier_to_cool_after_days_since_modification_greater_than", "10"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.actions.0.base_blob.0.tier_to_archive_after_days_since_modification_greater_than", "50"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.actions.0.base_blob.0.delete_after_days_since_modification_greater_than", "100"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.actions.0.snapshot.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "rule.0.actions.0.snapshot.0.delete_after_days_since_creation_greater_than", "30"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMStorageManagementPolicy_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = "${azurerm_resource_group.test.name}"

  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
  account_kind             = "BlobStorage"
}

resource "azurerm_storage_management_policy" "test" {
  storage_account_id = "${azurerm_storage_account.test.id}"

  rule {
    name    = "rule1"
    enabled = true
    filters {
      prefix_match = ["container1/prefix1"]
      blob_types   = ["blockBlob"]
    }
    actions {
      base_blob {
        tier_to_cool_after_days_since_modification_greater_than    = 10
        tier_to_archive_after_days_since_modification_greater_than = 50
        delete_after_days_since_modification_greater_than          = 100
      }
      snapshot {
        delete_after_days_since_creation_greater_than = 30
      }
    }
  }
}

data "azurerm_storage_management_policy" "test" {
  storage_account_id = "${azurerm_storage_management_policy.test.storage_account_id}"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
