package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceArmStorageContainer_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_storage_container", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMStorageContainer_basic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "container_access_type", "private"),
					resource.TestCheckResourceAttr(data.ResourceName, "has_immutability_policy", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "metadata.%", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "metadata.k1", "v1"),
					resource.TestCheckResourceAttr(data.ResourceName, "metadata.k2", "v2"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMStorageContainer_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "containerdstest-%s"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                = "acctestsadsc%s"
  resource_group_name = "${azurerm_resource_group.test.name}"

  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "containerdstest-%s"
  storage_account_name  = "${azurerm_storage_account.test.name}"
  container_access_type = "private"
  metadata = {
    k1 = "v1"
    k2 = "v2"
  }
}

data "azurerm_storage_container" "test" {
  name                 = azurerm_storage_container.test.name
  storage_account_name = azurerm_storage_container.test.storage_account_name
}

`, data.RandomString, data.Locations.Primary, data.RandomString, data.RandomString)
}
