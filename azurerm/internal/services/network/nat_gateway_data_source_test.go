package network_test

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type NatGatewayDataSource struct {
}

func TestAccDataSourceatGateway_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_nat_gateway", "test")
	r := NatGatewayDataSource{}
	// Using alt location because the resource currently in private preview and is only available in eastus2.

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("location").Exists(),
			),
		},
	})
}

func TestAccDataSourceatGateway_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_nat_gateway", "test")
	r := NatGatewayDataSource{}
	// Using alt location because the resource currently in private preview and is only available in eastus2.

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("public_ip_address_ids.#").HasValue("1"),
				check.That(data.ResourceName).Key("public_ip_prefix_ids.#").HasValue("1"),
				check.That(data.ResourceName).Key("sku_name").HasValue("Standard"),
				check.That(data.ResourceName).Key("idle_timeout_in_minutes").HasValue("10"),
			),
		},
	})
}

func (NatGatewayDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_nat_gateway" "test" {
  resource_group_name = azurerm_nat_gateway.test.resource_group_name
  name                = azurerm_nat_gateway.test.name
}
`, NatGatewayResource{}.basic(data))
}

func (NatGatewayDataSource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_nat_gateway" "test" {
  resource_group_name = azurerm_nat_gateway.test.resource_group_name
  name                = azurerm_nat_gateway.test.name
}
`, NatGatewayResource{}.complete(data))
}
