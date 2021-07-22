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
				check.That(data.ResourceName).Key("scope").HasValue("SQLDB"),
				check.That(data.ResourceName).Key("visibility").HasValue("Custom"),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.enV").HasValue("TesT"),
				check.That(data.ResourceName).Key("window.0.start_date_time").HasValue("5555-12-31 00:00"),
				check.That(data.ResourceName).Key("window.0.expiration_date_time").HasValue("6666-12-31 00:00"),
				check.That(data.ResourceName).Key("window.0.duration").HasValue("06:00"),
				check.That(data.ResourceName).Key("window.0.time_zone").HasValue("Pacific Standard Time"),
				check.That(data.ResourceName).Key("window.0.recur_every").HasValue("2Days"),
				check.That(data.ResourceName).Key("properties.%").HasValue("1"),
				check.That(data.ResourceName).Key("properties.description").HasValue("acceptance test"),
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
