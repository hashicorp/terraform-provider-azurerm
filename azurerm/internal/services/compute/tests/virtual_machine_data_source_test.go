package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMVirtualMachine_basicLinux(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMVirtualMachine_basicLinux(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "identity.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "identity.0.type", "SystemAssigned"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "identity.0.principal_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "identity.0.tenant_id"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMVirtualMachine_basicWindows(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMVirtualMachine_basicWindows(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "identity.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "identity.0.type", "SystemAssigned"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "identity.0.principal_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "identity.0.tenant_id"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMVirtualMachine_basicLinux(data acceptance.TestData) string {
	template := testLinuxVirtualMachine_identitySystemAssigned(data)
	return fmt.Sprintf(`
%s

data "azurerm_virtual_machine" "test" {
  name                = azurerm_linux_virtual_machine.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, template)
}

func testAccDataSourceAzureRMVirtualMachine_basicWindows(data acceptance.TestData) string {
	template := testWindowsVirtualMachine_identitySystemAssigned(data)
	return fmt.Sprintf(`
%s

data "azurerm_virtual_machine" "test" {
  name                = azurerm_windows_virtual_machine.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, template)
}
