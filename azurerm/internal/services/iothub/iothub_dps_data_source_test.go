package iothub_test

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type IotHubDPSDataSource struct {
}

func TestAccDataSourceIotHubDPS_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_iothub_dps", "test")
	r := IotHubDPSDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("allocation_policy").Exists(),
				check.That(data.ResourceName).Key("device_provisioning_host_name").Exists(),
				check.That(data.ResourceName).Key("id_scope").Exists(),
				check.That(data.ResourceName).Key("service_operations_host_name").Exists(),
			),
		},
	})
}

func (IotHubDPSDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_iothub_dps" "test" {
  name                = azurerm_iothub_dps.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, IotHubDPSResource{}.basic(data))
}
