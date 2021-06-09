package storage_test

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type StorageManagementPolicyDataSource struct{}

func TestAccDataSourceStorageManagementPolicy_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_storage_management_policy", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: StorageManagementPolicyDataSource{}.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("rule.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.0.name").HasValue("rule1"),
				check.That(data.ResourceName).Key("rule.0.enabled").HasValue("true"),
				check.That(data.ResourceName).Key("rule.0.filters.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.0.filters.0.prefix_match.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.0.filters.0.blob_types.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.0.actions.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.0.actions.0.base_blob.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.0.actions.0.base_blob.0.tier_to_cool_after_days_since_modification_greater_than").HasValue("10"),
				check.That(data.ResourceName).Key("rule.0.actions.0.base_blob.0.tier_to_archive_after_days_since_modification_greater_than").HasValue("50"),
				check.That(data.ResourceName).Key("rule.0.actions.0.base_blob.0.delete_after_days_since_modification_greater_than").HasValue("100"),
				check.That(data.ResourceName).Key("rule.0.actions.0.snapshot.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.0.actions.0.snapshot.0.delete_after_days_since_creation_greater_than").HasValue("30"),
			),
		},
	})
}

func TestAccDataSourceStorageManagementPolicy_blobTypes(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_storage_management_policy", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: StorageManagementPolicyDataSource{}.blobTypes(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("rule.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.0.name").HasValue("rule1"),
				check.That(data.ResourceName).Key("rule.0.enabled").HasValue("true"),
				check.That(data.ResourceName).Key("rule.0.filters.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.0.filters.0.prefix_match.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.0.filters.0.blob_types.#").HasValue("2"),
				check.That(data.ResourceName).Key("rule.0.actions.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.0.actions.0.base_blob.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.0.actions.0.base_blob.0.delete_after_days_since_modification_greater_than").HasValue("100"),
				check.That(data.ResourceName).Key("rule.0.actions.0.snapshot.#").HasValue("1"),
				check.That(data.ResourceName).Key("rule.0.actions.0.snapshot.0.delete_after_days_since_creation_greater_than").HasValue("30"),
			),
		},
	})
}

func (d StorageManagementPolicyDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = azurerm_resource_group.test.name

  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
  account_kind             = "BlobStorage"
}

resource "azurerm_storage_management_policy" "test" {
  storage_account_id = azurerm_storage_account.test.id

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
  storage_account_id = azurerm_storage_management_policy.test.storage_account_id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (d StorageManagementPolicyDataSource) blobTypes(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = azurerm_resource_group.test.name

  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
  account_kind             = "BlobStorage"
}

resource "azurerm_storage_management_policy" "test" {
  storage_account_id = azurerm_storage_account.test.id

  rule {
    name    = "rule1"
    enabled = true
    filters {
      prefix_match = ["container1/prefix1"]
      blob_types   = ["blockBlob", "appendBlob"]
    }
    actions {
      base_blob {
        delete_after_days_since_modification_greater_than = 100
      }
      snapshot {
        delete_after_days_since_creation_greater_than = 30
      }
    }
  }
}

data "azurerm_storage_management_policy" "test" {
  storage_account_id = azurerm_storage_management_policy.test.storage_account_id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
