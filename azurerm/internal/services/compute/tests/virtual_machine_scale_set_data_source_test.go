package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMVirtualMachineScaleSet_basicLinux(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_virtual_machine_scale_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMVirtualMachineScaleSet_basicLinux(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "identity.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "identity.0.type", "SystemAssigned"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "identity.0.principal_id"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMVirtualMachineScaleSet_basicWindows(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_virtual_machine_scale_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMVirtualMachineScaleSet_basicWindows(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "identity.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "identity.0.type", "SystemAssigned"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "identity.0.principal_id"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMVirtualMachineScaleSet_orchestrated(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_virtual_machine_scale_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMVirtualMachineScaleSet_orchestrated(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "id"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMVirtualMachineScaleSet_basicLinux(data acceptance.TestData) string {
	template := testAccAzureRMLinuxVirtualMachineScaleSet_identitySystemAssigned(data)
	return fmt.Sprintf(`
%s

data "azurerm_virtual_machine_scale_set" "test" {
  name                = azurerm_linux_virtual_machine_scale_set.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, template)
}

func testAccDataSourceAzureRMVirtualMachineScaleSet_basicWindows(data acceptance.TestData) string {
	template := testAccAzureRMWindowsVirtualMachineScaleSet_identitySystemAssigned(data)
	return fmt.Sprintf(`
%s

data "azurerm_virtual_machine_scale_set" "test" {
  name                = azurerm_windows_virtual_machine_scale_set.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, template)
}

func testAccDataSourceAzureRMVirtualMachineScaleSet_orchestrated(data acceptance.TestData) string {
	template := testAccAzureRMWindowsVirtualMachine_orchestratedZonal(data)
	return fmt.Sprintf(`
%s

data "azurerm_virtual_machine_scale_set" "test" {
  name                = azurerm_orchestrated_virtual_machine_scale_set.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, template)
}
