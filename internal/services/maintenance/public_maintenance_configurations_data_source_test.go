package maintenance_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type PublicMaintenanceConfigurationsDataSource struct{}

func TestAccDataSourcePublicMaintenanceConfigurations_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_public_maintenance_configurations", "test")
	r := PublicMaintenanceConfigurationsDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("public_maintenance_configurations.0.maintenance_scope").HasValue("SQLManagedInstance"),
			),
		},
	})
}

func (PublicMaintenanceConfigurationsDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
data "azurerm_public_maintenance_configurations" "test" {
  location_filter = "%s"
  scope_filter = "SQLManagedInstance"
  recur_window_filter = "weekMondayToThursday"
}
`, data.Locations.Primary)
}
