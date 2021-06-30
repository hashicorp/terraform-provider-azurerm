package maintenance_test

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type MaintenanceConfigurationDataSource struct {
}

func TestAccMaintenanceConfigurationDataSource_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_maintenance_configuration", "test")
	r := MaintenanceConfigurationDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("scope").HasValue("Host"),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.env").HasValue("TesT"),
			),
		},
	})
}

func (MaintenanceConfigurationDataSource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_maintenance_configuration" "test" {
  name                = azurerm_maintenance_configuration.test.name
  resource_group_name = azurerm_maintenance_configuration.test.resource_group_name
}
`, MaintenanceConfigurationResource{}.complete(data))
}
