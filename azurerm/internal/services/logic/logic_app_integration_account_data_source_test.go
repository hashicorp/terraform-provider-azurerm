package logic_test

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type LogicAppIntegrationAccountDataSource struct {
}

func TestAccLogicAppIntegrationAccountDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_logic_app_integration_account", "test")
	r := LogicAppIntegrationAccountDataSource{}
	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("location").Exists(),
				check.That(data.ResourceName).Key("sku_name").Exists(),
				check.That(data.ResourceName).Key("tags.%").Exists(),
				check.That(data.ResourceName).Key("tags.ENV").Exists(),
			),
		},
	})
}

func (LogicAppIntegrationAccountDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_logic_app_integration_account" "test" {
  name                = azurerm_logic_app_integration_account.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, LogicAppIntegrationAccountResource{}.complete(data))
}
