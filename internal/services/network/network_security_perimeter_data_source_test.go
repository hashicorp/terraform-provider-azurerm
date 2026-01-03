package network_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type NetworkSecurityPerimeterDataSource struct{}

func TestAccNetworkSecurityPerimeterDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_network_security_perimeter", "test")
	r := NetworkSecurityPerimeterDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("location").HasValue(location.Normalize(data.Locations.Primary)),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.env").HasValue("test"),
			),
		},
	})
}

func (NetworkSecurityPerimeterDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_network_security_perimeter" "test" {
  name     = "acctestNsp-%d"
  resource_group_name = azurerm_resource_group.test.name
  location = "%s"

  tags = {
    env = "test"
  }
}

data "azurerm_network_security_perimeter" "test" {
  name = azurerm_network_security_perimeter.test.name
  resource_group_name = azurerm_network_security_perimeter.test.resource_group_name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Primary)
}
