package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type PublicIPPrefixDataSource struct {
}

func TestAccDataSourcePublicIPPrefix_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_public_ip_prefix", "test")
	r := PublicIPPrefixDataSource{}
	name := fmt.Sprintf("acctestpublicipprefix-%d", data.RandomInteger)
	resourceGroupName := fmt.Sprintf("acctestRG-%d", data.RandomInteger)

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.basic(name, resourceGroupName, data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").HasValue(name),
				check.That(data.ResourceName).Key("resource_group_name").HasValue(resourceGroupName),
				check.That(data.ResourceName).Key("location").HasValue(data.Locations.Primary),
				check.That(data.ResourceName).Key("sku").HasValue("Standard"),
				check.That(data.ResourceName).Key("prefix_length").HasValue("31"),
				check.That(data.ResourceName).Key("ip_prefix").Exists(),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.env").HasValue("test"),
			),
		},
	})
}

func (PublicIPPrefixDataSource) basic(name string, resourceGroupName string, data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "%s"
  location = "%s"

  tags = {
    env = "test"
  }
}

resource "azurerm_public_ip_prefix" "test" {
  name                = "%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
  prefix_length       = 31

  tags = {
    env = "test"
  }
}

data "azurerm_public_ip_prefix" "test" {
  name                = azurerm_public_ip_prefix.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, resourceGroupName, data.Locations.Primary, name)
}
