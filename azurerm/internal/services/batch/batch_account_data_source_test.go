package batch_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type BatchAccountDataSource struct {
}

func TestAccBatchAccountDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_batch_account", "test")
	r := BatchAccountDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("location").HasValue(azure.NormalizeLocation(data.Locations.Primary)),
				check.That(data.ResourceName).Key("pool_allocation_mode").HasValue("BatchService"),
			),
		},
	})
}

func TestAccBatchAccountDataSource_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_batch_account", "test")
	r := BatchAccountDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("location").HasValue(azure.NormalizeLocation(data.Locations.Primary)),
				check.That(data.ResourceName).Key("pool_allocation_mode").HasValue("BatchService"),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.env").HasValue("test"),
			),
		},
	})
}

func TestAccBatchAccountDataSource_userSubscription(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_batch_account", "test")
	r := BatchAccountDataSource{}

	tenantID := os.Getenv("ARM_TENANT_ID")
	subscriptionID := os.Getenv("ARM_SUBSCRIPTION_ID")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.userSubscription(data, tenantID, subscriptionID),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("location").HasValue(azure.NormalizeLocation(data.Locations.Primary)),
				check.That(data.ResourceName).Key("pool_allocation_mode").HasValue("UserSubscription"),
				check.That(data.ResourceName).Key("key_vault_reference.#").HasValue("1"),
			),
		},
	})
}

func (BatchAccountDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "testaccRG-batch-%d"
  location = "%s"
}

resource "azurerm_batch_account" "test" {
  name                 = "testaccbatch%s"
  resource_group_name  = azurerm_resource_group.test.name
  location             = azurerm_resource_group.test.location
  pool_allocation_mode = "BatchService"
}

data "azurerm_batch_account" "test" {
  name                = azurerm_batch_account.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (BatchAccountDataSource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "testaccRG-batch-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "testaccsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_batch_account" "test" {
  name                 = "testaccbatch%s"
  resource_group_name  = azurerm_resource_group.test.name
  location             = azurerm_resource_group.test.location
  pool_allocation_mode = "BatchService"
  storage_account_id   = azurerm_storage_account.test.id

  tags = {
    env = "test"
  }
}

data "azurerm_batch_account" "test" {
  name                = azurerm_batch_account.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString)
}

func (BatchAccountDataSource) userSubscription(data acceptance.TestData, tenantID string, subscriptionID string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

provider "azuread" {}

data "azuread_service_principal" "test" {
  display_name = "Microsoft Azure Batch"
}

resource "azurerm_resource_group" "test" {
  name     = "testaccRG-batch-%d"
  location = "%s"
}

resource "azurerm_key_vault" "test" {
  name                            = "batchkv%s"
  location                        = "${azurerm_resource_group.test.location}"
  resource_group_name             = "${azurerm_resource_group.test.name}"
  enabled_for_disk_encryption     = true
  enabled_for_deployment          = true
  enabled_for_template_deployment = true
  tenant_id                       = "%s"

  sku_name = "standard"

  access_policy {
    tenant_id = "%s"
    object_id = "${data.azuread_service_principal.test.object_id}"

    secret_permissions = [
      "get",
      "list",
      "set",
      "delete",
      "recover"
    ]

  }
}

resource "azurerm_role_assignment" "contribrole" {
  scope                = "/subscriptions/%s"
  role_definition_name = "Contributor"
  principal_id         = "${data.azuread_service_principal.test.object_id}"
}

resource "azurerm_batch_account" "test" {
  name                = "testaccbatch%s"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  pool_allocation_mode = "UserSubscription"

  key_vault_reference {
    id  = "${azurerm_key_vault.test.id}"
    url = "${azurerm_key_vault.test.vault_uri}"
  }
}

data "azurerm_batch_account" "test" {
  name                = "${azurerm_batch_account.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, tenantID, tenantID, subscriptionID, data.RandomString)
}
