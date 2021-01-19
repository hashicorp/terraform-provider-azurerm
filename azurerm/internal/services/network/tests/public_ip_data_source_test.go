package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type PublicIPDataSource struct {
}

func TestAccDataSourcePublicIP_static(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_public_ip", "test")
	r := PublicIPDataSource{}

	name := fmt.Sprintf("acctestpublicip-%d", data.RandomInteger)
	resourceGroupName := fmt.Sprintf("acctestRG-%d", data.RandomInteger)

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.static(name, resourceGroupName, data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").HasValue(name),
				check.That(data.ResourceName).Key("resource_group_name").HasValue(resourceGroupName),
				check.That(data.ResourceName).Key("domain_name_label").HasValue(fmt.Sprintf("acctest-%d", data.RandomInteger)),
				check.That(data.ResourceName).Key("idle_timeout_in_minutes").HasValue("30"),
				check.That(data.ResourceName).Key("fqdn").Exists(),
				check.That(data.ResourceName).Key("ip_address").Exists(),
				check.That(data.ResourceName).Key("ip_version").HasValue("IPv4"),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.environment").HasValue("test"),
			),
		},
	})
}

func TestAccDataSourcePublicIP_dynamic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_public_ip", "test")
	r := PublicIPDataSource{}

	name := fmt.Sprintf("acctestpublicip-%d", data.RandomInteger)
	resourceGroupName := fmt.Sprintf("acctestRG-%d", data.RandomInteger)

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.dynamic(data, "Ipv4"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").HasValue(name),
				check.That(data.ResourceName).Key("resource_group_name").HasValue(resourceGroupName),
				check.That(data.ResourceName).Key("domain_name_label").HasValue(""),
				check.That(data.ResourceName).Key("fqdn").HasValue(""),
				check.That(data.ResourceName).Key("ip_address").HasValue(""),
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

func (PublicIPDataSource) dynamic(data acceptance.TestData, ipVersion string) string {
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
  allocation_method   = "Dynamic"

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
