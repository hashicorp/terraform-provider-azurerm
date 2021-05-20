package appconfiguration_test

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type AppConfigurationDataSource struct {
}

func TestAccAppConfigurationDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_app_configuration", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: AppConfigurationResource{}.standard(data),
		},
		{
			Config: AppConfigurationDataSource{}.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(AppConfigurationResource{}),
				check.That(data.ResourceName).Key("endpoint").Exists(),
				check.That(data.ResourceName).Key("location").Exists(),
				check.That(data.ResourceName).Key("sku").Exists(),
				check.That(data.ResourceName).Key("primary_read_key.0.connection_string").Exists(),
				check.That(data.ResourceName).Key("primary_read_key.0.id").Exists(),
				check.That(data.ResourceName).Key("primary_read_key.0.secret").Exists(),
				check.That(data.ResourceName).Key("primary_write_key.0.connection_string").Exists(),
				check.That(data.ResourceName).Key("primary_write_key.0.id").Exists(),
				check.That(data.ResourceName).Key("primary_write_key.0.secret").Exists(),
				check.That(data.ResourceName).Key("secondary_read_key.0.connection_string").Exists(),
				check.That(data.ResourceName).Key("secondary_read_key.0.id").Exists(),
				check.That(data.ResourceName).Key("secondary_read_key.0.secret").Exists(),
				check.That(data.ResourceName).Key("secondary_write_key.0.connection_string").Exists(),
				check.That(data.ResourceName).Key("secondary_write_key.0.id").Exists(),
				check.That(data.ResourceName).Key("secondary_write_key.0.secret").Exists(),
			),
		},
	})
}

func (AppConfigurationDataSource) basic(data acceptance.TestData) string {
	template := AppConfigurationResource{}.standard(data)
	return fmt.Sprintf(`
%s

data "azurerm_app_configuration" "test" {
  name                = azurerm_app_configuration.test.name
  resource_group_name = azurerm_app_configuration.test.resource_group_name
}
`, template)
}
