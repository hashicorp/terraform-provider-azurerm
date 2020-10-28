package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMApiManagement_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_api_management", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceApiManagement_basic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "publisher_email", "pub1@email.com"),
					resource.TestCheckResourceAttr(data.ResourceName, "publisher_name", "pub1"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku_name", "Developer_1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "0"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "public_ip_addresses.#"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMApiManagement_virtualNetwork(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_api_management", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceApiManagement_virtualNetwork(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "publisher_email", "pub1@email.com"),
					resource.TestCheckResourceAttr(data.ResourceName, "publisher_name", "pub1"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku_name", "Premium_1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "0"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "public_ip_addresses.#"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "private_ip_addresses.#"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "additional_location.0.public_ip_addresses.#"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "additional_location.0.private_ip_addresses.#"),
				),
			},
		},
	})
}

func testAccDataSourceApiManagement_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "amtestRG-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"
  sku_name            = "Developer_1"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

data "azurerm_api_management" "test" {
  name                = azurerm_api_management.test.name
  resource_group_name = azurerm_api_management.test.resource_group_name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccDataSourceApiManagement_virtualNetwork(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test1" {
  name     = "amestRG1-%d"
  location = "%s"
}

resource "azurerm_resource_group" "test2" {
  name     = "amestRG2-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test1" {
  name                = "amtestVNET1-%d"
  location            = azurerm_resource_group.test1.location
  resource_group_name = azurerm_resource_group.test1.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "test1" {
  name               = "amtestSNET1-%d"
  virtual_network_id = azurerm_virtual_network.test1.id
  address_prefix     = "10.0.1.0/24"
}

resource "azurerm_virtual_network" "test2" {
  name                = "amtestVNET2-%d"
  location            = azurerm_resource_group.test2.location
  resource_group_name = azurerm_resource_group.test2.name
  address_space       = ["10.1.0.0/16"]
}

resource "azurerm_subnet" "test2" {
  name               = "amtestSNET2-%d"
  virtual_network_id = azurerm_virtual_network.test2.id
  address_prefix     = "10.1.1.0/24"
}

resource "azurerm_api_management" "test" {
  name                = "amtestAM-%d"
  location            = azurerm_resource_group.test1.location
  resource_group_name = azurerm_resource_group.test1.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"

  sku_name = "Premium_1"

  additional_location {
    location = azurerm_resource_group.test2.location
    virtual_network_configuration {
      subnet_id = azurerm_subnet.test2.id
    }
  }

  virtual_network_type = "Internal"
  virtual_network_configuration {
    subnet_id = azurerm_subnet.test1.id
  }
}

data "azurerm_api_management" "test" {
  name                = azurerm_api_management.test.name
  resource_group_name = azurerm_api_management.test.resource_group_name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Secondary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
