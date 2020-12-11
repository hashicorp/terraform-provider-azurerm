package storage_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMStorageAccount_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_storage_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMStorageAccount_basic(data),
			},
			{
				Config: testAccDataSourceAzureRMStorageAccount_basicWithDataSource(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "account_tier", "Standard"),
					resource.TestCheckResourceAttr(data.ResourceName, "account_replication_type", "LRS"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.environment", "production"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMStorageAccount_withWriteLock(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_storage_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMStorageAccount_basicWriteLock(data),
			},
			{
				Config: testAccDataSourceAzureRMStorageAccount_basicWriteLockWithDataSource(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "account_tier", "Standard"),
					resource.TestCheckResourceAttr(data.ResourceName, "account_replication_type", "LRS"),
					resource.TestCheckResourceAttr(data.ResourceName, "primary_connection_string", ""),
					resource.TestCheckResourceAttr(data.ResourceName, "secondary_connection_string", ""),
					resource.TestCheckResourceAttr(data.ResourceName, "primary_blob_connection_string", ""),
					resource.TestCheckResourceAttr(data.ResourceName, "secondary_blob_connection_string", ""),
					resource.TestCheckResourceAttr(data.ResourceName, "primary_access_key", ""),
					resource.TestCheckResourceAttr(data.ResourceName, "secondary_access_key", ""),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMStorageAccount_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                = "acctestsads%s"
  resource_group_name = azurerm_resource_group.test.name

  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "production"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccDataSourceAzureRMStorageAccount_basicWriteLock(data acceptance.TestData) string {
	template := testAccDataSourceAzureRMStorageAccount_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_management_lock" "test" {
  name       = "acctestlock-%d"
  scope      = azurerm_storage_account.test.id
  lock_level = "ReadOnly"
}
`, template, data.RandomInteger)
}

func testAccDataSourceAzureRMStorageAccount_basicWithDataSource(data acceptance.TestData) string {
	config := testAccDataSourceAzureRMStorageAccount_basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_storage_account" "test" {
  name                = azurerm_storage_account.test.name
  resource_group_name = azurerm_storage_account.test.resource_group_name
}
`, config)
}

func testAccDataSourceAzureRMStorageAccount_basicWriteLockWithDataSource(data acceptance.TestData) string {
	config := testAccDataSourceAzureRMStorageAccount_basicWriteLock(data)
	return fmt.Sprintf(`
%s

data "azurerm_storage_account" "test" {
  name                = azurerm_storage_account.test.name
  resource_group_name = azurerm_storage_account.test.resource_group_name
}
`, config)
}
