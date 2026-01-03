package network_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type NetworkSecurityPerimeterProfileDataSource struct{}

func TestAccNetworkSecurityPerimeterProfileDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_network_security_perimeter_profile", "test")
	r := NetworkSecurityPerimeterProfileDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").HasValue(fmt.Sprintf("acctestProfile-%d", data.RandomInteger)),
				check.That(data.ResourceName).Key("perimeter_id").Exists(),
				check.That(data.ResourceName).Key("id").Exists(),
			),
		},
	})
}

func (NetworkSecurityPerimeterProfileDataSource) basic(data acceptance.TestData) string {
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

resource "azurerm_network_security_perimeter_profile" "test" {
	name = "acctestProfile-%d"
	perimeter_id = azurerm_network_security_perimeter.test.id
}

data "azurerm_network_security_perimeter_profile" "test" {
	name = azurerm_network_security_perimeter_profile.test.name 
	perimeter_id = azurerm_network_security_perimeter.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
