// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type PublicIPDataSource struct{}

func TestAccDataSourcePublicIP_static(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_public_ip", "test")
	r := PublicIPDataSource{}

	name := fmt.Sprintf("acctestpublicip-%d", data.RandomInteger)
	resourceGroupName := fmt.Sprintf("acctestRG-%d", data.RandomInteger)

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.static(name, resourceGroupName, data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").HasValue(name),
				check.That(data.ResourceName).Key("resource_group_name").HasValue(resourceGroupName),
				check.That(data.ResourceName).Key("domain_name_label").HasValue(fmt.Sprintf("acctest-%d", data.RandomInteger)),
				check.That(data.ResourceName).Key("idle_timeout_in_minutes").HasValue("30"),
				check.That(data.ResourceName).Key("fqdn").Exists(),
				check.That(data.ResourceName).Key("ip_address").Exists(),
				check.That(data.ResourceName).Key("ip_version").HasValue("IPv4"),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.environment").HasValue("test"),
				check.That(data.ResourceName).Key("ip_tags.RoutingPreference").HasValue("Internet"),
				check.That(data.ResourceName).Key("sku_tier").HasValue("Regional"),
			),
		},
	})
}

func TestAccDataSourcePublicIP_staticMinimal(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_public_ip", "test")
	r := PublicIPDataSource{}

	name := fmt.Sprintf("acctestpublicip-%d", data.RandomInteger)
	resourceGroupName := fmt.Sprintf("acctestRG-%d", data.RandomInteger)

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.staticMinimal(data, "IPv4"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").HasValue(name),
				check.That(data.ResourceName).Key("resource_group_name").HasValue(resourceGroupName),
				check.That(data.ResourceName).Key("domain_name_label").HasValue(""),
				check.That(data.ResourceName).Key("fqdn").HasValue(""),
				check.That(data.ResourceName).Key("ip_address").Exists(),
				check.That(data.ResourceName).Key("ip_version").HasValue("IPv4"),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.environment").HasValue("test"),
			),
		},
	})
}

func (PublicIPDataSource) static(name string, resourceGroupName string, data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "%s"
  location = "%s"
}

resource "azurerm_public_ip" "test" {
  name                    = "%s"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  allocation_method       = "Static"
  domain_name_label       = "acctest-%d"
  idle_timeout_in_minutes = 30
  sku                     = "Standard"
  zones                   = ["1", "2", "3"]

  ip_tags = {
    RoutingPreference = "Internet"
  }

  tags = {
    environment = "test"
  }
}

data "azurerm_public_ip" "test" {
  name                = azurerm_public_ip.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, resourceGroupName, data.Locations.Primary, name, data.RandomInteger)
}

func (PublicIPDataSource) staticMinimal(data acceptance.TestData, ipVersion string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpublicip-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"

  ip_version = "%s"

  tags = {
    environment = "test"
  }
}

data "azurerm_public_ip" "test" {
  name                = azurerm_public_ip.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, ipVersion)
}

func TestAccDataSourcePublicIP_prefixAndDomainNameScope(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_public_ip", "test")
	r := PublicIPDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.prefixAndDomainNameScope(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("public_ip_prefix_id").IsSet(),
				check.That(data.ResourceName).Key("domain_name_label_scope").HasValue("TenantReuse"),
			),
		},
	})
}

func (PublicIPDataSource) prefixAndDomainNameScope(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_public_ip_prefix" "test" {
  name                = "acctestpublicipprefix-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_public_ip" "test" {
  name                    = "acctestpublicip-%[1]d"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  allocation_method       = "Static"
  sku                     = "Standard"
  public_ip_prefix_id     = azurerm_public_ip_prefix.test.id
  domain_name_label       = "acctest-%[1]d"
  domain_name_label_scope = "TenantReuse"
}

data "azurerm_public_ip" "test" {
  name                = azurerm_public_ip.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary)
}

func TestAccDataSourcePublicIP_edgeZone(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_public_ip", "test")
	r := PublicIPDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.edgeZone(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("edge_zone").IsSet(),
			),
		},
	})
}

func (PublicIPDataSource) edgeZone(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_public_ip" "test" {
  name                = azurerm_public_ip.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, PublicIpResource{}.edgeZone(data))
}
