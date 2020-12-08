package resource_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMResources_ByName(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_resources", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMResources_template(data),
			},
			{
				Config: testAccDataSourceAzureRMResources_ByName(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "resources.#", "1"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMResources_ByResourceGroup(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_resources", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMResources_template(data),
			},
			{
				Config: testAccDataSourceAzureRMResources_ByResourceGroup(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "resources.#", "1"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMResources_ByResourceType(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_resources", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMResources_template(data),
			},
			{
				Config: testAccDataSourceAzureRMResources_ByResourceType(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "resources.#", "1"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMResources_FilteredByTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_resources", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMResources_template(data),
			},
			{
				Config: testAccDataSourceAzureRMResources_FilteredByTags(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "resources.#", "1"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMResources_ByName(data acceptance.TestData) string {
	r := testAccDataSourceAzureRMResources_template(data)
	return fmt.Sprintf(`
%s

data "azurerm_resources" "test" {
  name = azurerm_storage_account.test.name
}
`, r)
}

func testAccDataSourceAzureRMResources_ByResourceGroup(data acceptance.TestData) string {
	r := testAccDataSourceAzureRMResources_template(data)
	return fmt.Sprintf(`
%s

data "azurerm_resources" "test" {
  resource_group_name = azurerm_storage_account.test.resource_group_name
}
`, r)
}

func testAccDataSourceAzureRMResources_ByResourceType(data acceptance.TestData) string {
	r := testAccDataSourceAzureRMResources_template(data)
	return fmt.Sprintf(`
%s

data "azurerm_resources" "test" {
  resource_group_name = azurerm_storage_account.test.resource_group_name
  type                = "Microsoft.Storage/storageAccounts"
}
`, r)
}

func testAccDataSourceAzureRMResources_FilteredByTags(data acceptance.TestData) string {
	r := testAccDataSourceAzureRMResources_template(data)
	return fmt.Sprintf(`
%s

data "azurerm_resources" "test" {
  name                = azurerm_storage_account.test.name
  resource_group_name = azurerm_storage_account.test.resource_group_name

  required_tags = {
    environment = "production"
  }
}
`, r)
}

func testAccDataSourceAzureRMResources_template(data acceptance.TestData) string {
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
