package signalr_test

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type SignalRServiceDataSource struct{}

func TestAccDataSourceSignalRService_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_signalr_service", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: SignalRServiceDataSource{}.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("id").Exists(),
				check.That(data.ResourceName).Key("hostname").Exists(),
				check.That(data.ResourceName).Key("ip_address").Exists(),
				check.That(data.ResourceName).Key("public_port").Exists(),
				check.That(data.ResourceName).Key("server_port").Exists(),
				check.That(data.ResourceName).Key("primary_access_key").Exists(),
				check.That(data.ResourceName).Key("primary_connection_string").Exists(),
				check.That(data.ResourceName).Key("secondary_access_key").Exists(),
				check.That(data.ResourceName).Key("secondary_connection_string").Exists(),
			),
		},
	})
}

func (r SignalRServiceDataSource) basic(data acceptance.TestData) string {
	template := SignalRServiceResource{}.basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_signalr_service" "test" {
  name                = azurerm_signalr_service.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, template)
}
