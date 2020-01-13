package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMBatchAccount_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_batch_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMBatchAccount_basic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "name", fmt.Sprintf("testaccbatch%s", data.RandomString)),
					resource.TestCheckResourceAttr(data.ResourceName, "location", azure.NormalizeLocation(data.Locations.Primary)),
					resource.TestCheckResourceAttr(data.ResourceName, "pool_allocation_mode", "BatchService"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMBatchAccount_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_batch_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMBatchAccount_complete(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "name", fmt.Sprintf("testaccbatch%s", data.RandomString)),
					resource.TestCheckResourceAttr(data.ResourceName, "location", azure.NormalizeLocation(data.Locations.Primary)),
					resource.TestCheckResourceAttr(data.ResourceName, "pool_allocation_mode", "BatchService"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.env", "test"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMBatchAccount_userSubscription(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_batch_account", "test")

	tenantID := os.Getenv("ARM_TENANT_ID")
	subscriptionID := os.Getenv("ARM_SUBSCRIPTION_ID")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureBatchAccount_userSubscription(data, tenantID, subscriptionID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "name", fmt.Sprintf("testaccbatch%s", data.RandomString)),
					resource.TestCheckResourceAttr(data.ResourceName, "location", azure.NormalizeLocation(data.Locations.Primary)),
					resource.TestCheckResourceAttr(data.ResourceName, "pool_allocation_mode", "UserSubscription"),
					resource.TestCheckResourceAttr(data.ResourceName, "key_vault_reference.#", "1"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMBatchAccount_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "testaccRG-%d-batch"
  location = "%s"
}

resource "azurerm_batch_account" "test" {
  name                 = "testaccbatch%s"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  location             = "${azurerm_resource_group.test.location}"
  pool_allocation_mode = "BatchService"
}

data "azurerm_batch_account" "test" {
  name                = "${azurerm_batch_account.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccDataSourceAzureRMBatchAccount_complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "testaccRG-%d-batch"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "testaccsa%s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_batch_account" "test" {
  name                 = "testaccbatch%s"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  location             = "${azurerm_resource_group.test.location}"
  pool_allocation_mode = "BatchService"
  storage_account_id   = "${azurerm_storage_account.test.id}"

  tags = {
    env = "test"
  }
}

data "azurerm_batch_account" "test" {
  name                = "${azurerm_batch_account.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString)
}

func testAccDataSourceAzureBatchAccount_userSubscription(data acceptance.TestData, tenantID string, subscriptionID string) string {
	return fmt.Sprintf(`
data "azurerm_azuread_service_principal" "test" {
  display_name = "Microsoft Azure Batch"
}

resource "azurerm_resource_group" "test" {
  name     = "testaccRG-%d-batchaccount"
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

  sku {
    name = "standard"
  }

  access_policy {
    tenant_id = "%s"
    object_id = "${data.azurerm_azuread_service_principal.test.object_id}"

    secret_permissions = [
      "get",
      "list",
      "set",
      "delete"
    ]

  }
}

resource "azurerm_role_assignment" "contribrole" {
  scope                = "/subscriptions/%s"
  role_definition_name = "Contributor"
  principal_id         = "${data.azurerm_azuread_service_principal.test.object_id}"
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
