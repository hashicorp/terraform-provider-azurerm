package network_test

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type AppGatewayDataSource struct {
}

func TestAccDataSourceAppGateway_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_application_gateway", "test")
	r := AppGatewayDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("location").Exists(),
			),
		},
	})
}

func TestAccDataSourceAppGateway_userAssignedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_application_gateway", "test")
	r := AppGatewayDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.userAssignedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("location").Exists(),
				check.That(data.ResourceName).Key("identity.0.identity_ids.#").HasValue("1"),
			),
		},
	})
}

func (AppGatewayDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_application_gateway" "test" {
  resource_group_name = azurerm_application_gateway.test.resource_group_name
  name                = azurerm_application_gateway.test.name
}
`, ApplicationGatewayResource{}.basic(data))
}

func (AppGatewayDataSource) userAssignedIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_application_gateway" "test" {
  resource_group_name = azurerm_application_gateway.test.resource_group_name
  name                = azurerm_application_gateway.test.name
}
`, ApplicationGatewayResource{}.UserDefinedIdentity(data))
}
