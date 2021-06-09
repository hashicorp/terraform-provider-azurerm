package apimanagement_test

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type ApiManagementDataSource struct {
}

func TestAccDataSourceApiManagement_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_api_management", "test")
	r := ApiManagementDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("publisher_email").HasValue("pub1@email.com"),
				check.That(data.ResourceName).Key("publisher_name").HasValue("pub1"),
				check.That(data.ResourceName).Key("sku_name").HasValue("Developer_1"),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
				check.That(data.ResourceName).Key("public_ip_addresses.#").Exists(),
			),
		},
	})
}

func TestAccDataSourceApiManagement_identitySystemAssigned(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_api_management", "test")
	r := ApiManagementDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.identitySystemAssigned(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("publisher_email").HasValue("pub1@email.com"),
				check.That(data.ResourceName).Key("publisher_name").HasValue("pub1"),
				check.That(data.ResourceName).Key("sku_name").HasValue("Developer_1"),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
				check.That(data.ResourceName).Key("public_ip_addresses.#").Exists(),
				check.That(data.ResourceName).Key("identity.#").HasValue("1"),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("SystemAssigned"),
			),
		},
	})
}

func TestAccDataSourceApiManagement_identityUserAssigned(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_api_management", "test")
	r := ApiManagementDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.identityUserAssigned(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("publisher_email").HasValue("pub1@email.com"),
				check.That(data.ResourceName).Key("publisher_name").HasValue("pub1"),
				check.That(data.ResourceName).Key("sku_name").HasValue("Developer_1"),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
				check.That(data.ResourceName).Key("public_ip_addresses.#").Exists(),
				check.That(data.ResourceName).Key("identity.#").HasValue("1"),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("UserAssigned"),
			),
		},
	})
}

func TestAccDataSourceApiManagement_virtualNetwork(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_api_management", "test")
	r := ApiManagementDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.virtualNetwork(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("publisher_email").HasValue("pub1@email.com"),
				check.That(data.ResourceName).Key("publisher_name").HasValue("pub1"),
				check.That(data.ResourceName).Key("sku_name").HasValue("Premium_1"),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
				check.That(data.ResourceName).Key("public_ip_addresses.#").Exists(),
				check.That(data.ResourceName).Key("private_ip_addresses.#").Exists(),
				check.That(data.ResourceName).Key("additional_location.0.public_ip_addresses.#").Exists(),
				check.That(data.ResourceName).Key("additional_location.0.private_ip_addresses.#").Exists(),
			),
		},
	})
}

func (ApiManagementDataSource) basic(data acceptance.TestData) string {
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

func (ApiManagementDataSource) identitySystemAssigned(data acceptance.TestData) string {
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

  identity {
    type = "SystemAssigned"
  }
}

data "azurerm_api_management" "test" {
  name                = azurerm_api_management.test.name
  resource_group_name = azurerm_api_management.test.resource_group_name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (ApiManagementDataSource) identityUserAssigned(data acceptance.TestData) string {
	template := ApiManagementResource{}.identityUserAssigned(data)
	return fmt.Sprintf(`
%s

data "azurerm_api_management" "test" {
  name                = azurerm_api_management.test.name
  resource_group_name = azurerm_api_management.test.resource_group_name
}
`, template)
}

func (ApiManagementDataSource) virtualNetwork(data acceptance.TestData) string {
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
  name                 = "amtestSNET1-%d"
  resource_group_name  = azurerm_resource_group.test1.name
  virtual_network_name = azurerm_virtual_network.test1.name
  address_prefix       = "10.0.1.0/24"
}

resource "azurerm_virtual_network" "test2" {
  name                = "amtestVNET2-%d"
  location            = azurerm_resource_group.test2.location
  resource_group_name = azurerm_resource_group.test2.name
  address_space       = ["10.1.0.0/16"]
}

resource "azurerm_subnet" "test2" {
  name                 = "amtestSNET2-%d"
  resource_group_name  = azurerm_resource_group.test2.name
  virtual_network_name = azurerm_virtual_network.test2.name
  address_prefix       = "10.1.1.0/24"
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
