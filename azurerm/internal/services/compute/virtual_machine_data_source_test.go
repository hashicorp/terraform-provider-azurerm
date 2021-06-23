package compute_test

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type VirtualMachineDataSource struct {
}

func TestAccDataSourceAzureRMVirtualMachine_basicLinux(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_virtual_machine", "test")
	r := VirtualMachineDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basicLinux(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("identity.#").HasValue("1"),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("SystemAssigned"),
				check.That(data.ResourceName).Key("identity.0.principal_id").Exists(),
				check.That(data.ResourceName).Key("identity.0.tenant_id").Exists(),
			),
		},
	})
}

func TestAccDataSourceAzureRMVirtualMachine_basicWindows(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_virtual_machine", "test")
	r := VirtualMachineDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basicWindows(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("identity.#").HasValue("1"),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("SystemAssigned"),
				check.That(data.ResourceName).Key("identity.0.principal_id").Exists(),
				check.That(data.ResourceName).Key("identity.0.tenant_id").Exists(),
			),
		},
	})
}

func (VirtualMachineDataSource) basicLinux(data acceptance.TestData) string {
	template := LinuxVirtualMachineResource{}.identitySystemAssigned(data)
	return fmt.Sprintf(`
%s

data "azurerm_virtual_machine" "test" {
  name                = azurerm_linux_virtual_machine.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, template)
}

func (VirtualMachineDataSource) basicWindows(data acceptance.TestData) string {
	template := WindowsVirtualMachineResource{}.identitySystemAssigned(data)
	return fmt.Sprintf(`
%s

data "azurerm_virtual_machine" "test" {
  name                = azurerm_windows_virtual_machine.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, template)
}
