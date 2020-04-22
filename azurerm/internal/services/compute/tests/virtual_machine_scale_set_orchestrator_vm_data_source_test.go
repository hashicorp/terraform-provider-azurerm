package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMWindowsVirtualMachineScaleSetOrchestratorVM_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_virtual_machine_scale_set_orchestrator_vm", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMWindowsVirtualMachineScaleSetOrchestratorVMDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMWindowsVirtualMachineScaleSetOrchestratorVM_basic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.ENV", "Test"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMWindowsVirtualMachineScaleSetOrchestratorVM_basic(data acceptance.TestData) string {
	template := testAccAzureRMWindowsVirtualMachineScaleSetOrchestratorVM_basicLinux(data)
	return fmt.Sprintf(`
%s

data "azurerm_virtual_machine_scale_set_orchestrator_vm" "test" {
  name                = azurerm_virtual_machine_scale_set_orchestrator_vm.test.name
  resource_group_name = azurerm_virtual_machine_scale_set_orchestrator_vm.test.resource_group_name
}
`, template)
}
