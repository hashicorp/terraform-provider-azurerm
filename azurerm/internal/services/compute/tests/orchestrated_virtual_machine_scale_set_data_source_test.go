package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMOrchestratedVirtualMachineScaleSet_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_orchestrated_virtual_machine_scale_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMOrchestratedVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMOrchestratedVirtualMachineScaleSet_basic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "location"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "unique_id"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.ENV", "Test"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMOrchestratedVirtualMachineScaleSet_basic(data acceptance.TestData) string {
	template := testAccAzureRMOrchestratedVirtualMachineScaleSet_basicLinux(data)
	return fmt.Sprintf(`
%s

data "azurerm_orchestrated_virtual_machine_scale_set" "test" {
  name                = azurerm_orchestrated_virtual_machine_scale_set.test.name
  resource_group_name = azurerm_orchestrated_virtual_machine_scale_set.test.resource_group_name
}
`, template)
}
