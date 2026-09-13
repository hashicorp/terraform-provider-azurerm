// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type ImageDataSource struct{}

func TestAccDataSourceImage_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_image", "test")
	r := ImageDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("os_disk.#").HasValue("1"),
				check.That(data.ResourceName).Key("os_disk.0.caching").HasValue("None"),
				check.That(data.ResourceName).Key("os_disk.0.os_type").HasValue("Linux"),
				check.That(data.ResourceName).Key("os_disk.0.os_state").HasValue("Generalized"),
				check.That(data.ResourceName).Key("os_disk.0.size_gb").HasValue("30"),
				check.That(data.ResourceName).Key("data_disk.#").HasValue("0"),
				check.That(data.ResourceName).Key("tags.%").HasValue("2"),
				check.That(data.ResourceName).Key("tags.environment").HasValue("Dev"),
				check.That(data.ResourceName).Key("tags.cost-center").HasValue("Ops"),
			),
		},
	})
}

func TestAccDataSourceImage_localFilter(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_image", "test1")
	r := ImageDataSource{}

	descDataSourceName := "data.azurerm_image.test2"
	data.DataSourceTest(t, []acceptance.TestStep{
		{
			// We have to create the images first explicitly, then retrieve the data source, because in this case we do not have explicit dependency on the image resources
			Config: r.localFilter_setup(data),
		},
		{
			Config: r.localFilter(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("name").MatchesOtherKey(check.That(descDataSourceName).Key("name")),
				check.That(data.ResourceName).Key("resource_group_name").MatchesOtherKey(check.That(descDataSourceName).Key("resource_group_name")),
			),
		},
	})
}

func (r ImageDataSource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpip%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
  domain_name_label   = "acctestpip%[1]d"
}

resource "azurerm_network_interface" "testsource" {
  name                = "acctestnic-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "testconfigurationsource"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = azurerm_public_ip.test.id
  }
}

resource "azurerm_linux_virtual_machine" "testsource" {
  name                            = "acctestvm-%[1]d"
  location                        = azurerm_resource_group.test.location
  resource_group_name             = azurerm_resource_group.test.name
  network_interface_ids           = [azurerm_network_interface.testsource.id]
  size                            = "Standard_D1_v2"
  computer_name                   = "acctest-%[1]d"
  admin_username                  = "tfuser"
  admin_password                  = "P@ssW0RD7890"
  disable_password_authentication = false

  os_disk {
    name                 = "myosdisk1-%[1]d"
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }

  tags = {
    environment = "Dev"
    cost-center = "Ops"
  }
}

data "azurerm_managed_disk" "testsource" {
  name                = azurerm_linux_virtual_machine.testsource.os_disk.0.name
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r ImageDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_image" "test" {
  name                = "acctest-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  os_disk {
    os_type         = "Linux"
    os_state        = "Generalized"
    managed_disk_id = data.azurerm_managed_disk.testsource.id
    size_gb         = 30
    caching         = "None"
    storage_type    = "Standard_LRS"
  }

  tags = {
    environment = "Dev"
    cost-center = "Ops"
  }
}

data "azurerm_image" "test" {
  name                = azurerm_image.test.name
  resource_group_name = azurerm_resource_group.test.name
}

output "location" {
  value = data.azurerm_image.test.location
}
`, r.template(data), data.RandomInteger)
}

func (r ImageDataSource) localFilter_setup(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_image" "abc" {
  name                = "abc-acctest-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  os_disk {
    os_type         = "Linux"
    os_state        = "Generalized"
    managed_disk_id = data.azurerm_managed_disk.testsource.id
    size_gb         = 30
    caching         = "None"
    storage_type    = "Standard_LRS"
  }

  tags = {
    environment = "Dev"
    cost-center = "Ops"
  }
}

resource "azurerm_image" "def" {
  name                = "def-acctest-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  os_disk {
    os_type         = "Linux"
    os_state        = "Generalized"
    managed_disk_id = data.azurerm_managed_disk.testsource.id
    size_gb         = 30
    caching         = "None"
    storage_type    = "Standard_LRS"
  }

  tags = {
    environment = "Dev"
    cost-center = "Ops"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r ImageDataSource) localFilter(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_image" "test1" {
  name_regex          = "^def-acctest-\\d+"
  resource_group_name = azurerm_resource_group.test.name
}

data "azurerm_image" "test2" {
  name_regex          = "^[a-z]+-acctest-\\d+"
  sort_descending     = true
  resource_group_name = azurerm_resource_group.test.name
}
`, r.localFilter_setup(data))
}
