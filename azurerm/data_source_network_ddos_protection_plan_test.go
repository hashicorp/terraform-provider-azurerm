package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func testAccAzureRMNetworkDDoSProtectionPlanDataSource_basic(t *testing.T) {
	dsn := "azurerm_network_ddos_protection_plan.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkDDoSProtectionPlanDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkDDoSProtectionPlanDataSource_basicConfig(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkDDoSProtectionPlanExists(dsn),
					resource.TestCheckResourceAttrSet(dsn, "virtual_network_ids.#"),
				),
			},
		},
	})
}

func testAccAzureRMNetworkDDoSProtectionPlanDataSource_basicConfig(rInt int, location string) string {
	return fmt.Sprintf(`
%s

data "azurerm_network_ddos_protection_plan" "test" {
  name                = "${azurerm_network_ddos_protection_plan.test.name}"
  resource_group_name = "${azurerm_network_ddos_protection_plan.test.resource_group_name}"
}
`, testAccAzureRMNetworkDDoSProtectionPlan_basicConfig(rInt, location))
}
