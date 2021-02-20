package network_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type AppGatewayDataSource struct {
}

func TestAccDataSourceAppGateway_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_application_gateway", "test")
	r := AppGatewayDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("location").Exists(),
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
