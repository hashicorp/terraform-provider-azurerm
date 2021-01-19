package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type VirtualWanDataSource struct {
}

func TestAccDataSourceVirtualWan_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_virtual_wan", "test")
	r := VirtualWanDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
			),
		},
	})
}

func (r VirtualWanDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_virtual_wan" "test" {
  name                = "test"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

data "azurerm_virtual_wan" "test" {
  name                = azurerm_virtual_wan.test.name
  resource_group_name = azurerm_virtual_wan.test.resource_group_name
}
`, r.template(data))
}

func (VirtualWanDataSource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}
