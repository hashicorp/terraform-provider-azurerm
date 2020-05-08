package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func testAccAzureRMNetworkDDoSProtectionPlanDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_network_ddos_protection_plan", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkDDoSProtectionPlanDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkDDoSProtectionPlanDataSource_basicConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkDDoSProtectionPlanExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "virtual_network_ids.#"),
				),
			},
		},
	})
}

func testAccAzureRMNetworkDDoSProtectionPlanDataSource_basicConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_network_ddos_protection_plan" "test" {
  name                = azurerm_network_ddos_protection_plan.test.name
  resource_group_name = azurerm_network_ddos_protection_plan.test.resource_group_name
}
`, testAccAzureRMNetworkDDoSProtectionPlan_basicConfig(data))
}
