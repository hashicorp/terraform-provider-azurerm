package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type NetworkDDoSProtectionPlanDataSource struct {
}

func testAccNetworkDDoSProtectionPlanDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_network_ddos_protection_plan", "test")
	r := NetworkDDoSProtectionPlanDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.basicConfig(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("virtual_network_ids.#").Exists(),
			),
		},
	})
}

func (NetworkDDoSProtectionPlanDataSource) basicConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_network_ddos_protection_plan" "test" {
  name                = azurerm_network_ddos_protection_plan.test.name
  resource_group_name = azurerm_network_ddos_protection_plan.test.resource_group_name
}
`, NetworkDDoSProtectionPlanResource{}.basicConfig(data))
}
