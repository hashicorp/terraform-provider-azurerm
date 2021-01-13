package compute

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type VirtualMachineScaleSetDataSource struct {
}

func TestAccDataSourceVirtualMachineScaleSet_basicLinux(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_virtual_machine_scale_set", "test")
	r := VirtualMachineScaleSetDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.basicLinux(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("identity.#").HasValue("1"),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("SystemAssigned"),
				check.That(data.ResourceName).Key("identity.0.principal_id").Exists(),
			),
		},
	})
}

func TestAccDataSourceVirtualMachineScaleSet_basicWindows(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_virtual_machine_scale_set", "test")
	r := VirtualMachineScaleSetDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.basicWindows(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("identity.#").HasValue("1"),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("SystemAssigned"),
				check.That(data.ResourceName).Key("identity.0.principal_id").Exists(),
			),
		},
	})
}

func TestAccDataSourceVirtualMachineScaleSet_orchestrated(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_virtual_machine_scale_set", "test")
	r := VirtualMachineScaleSetDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.orchestrated(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("id").Exists(),
			),
		},
	})
}

func (VirtualMachineScaleSetDataSource) basicLinux(data acceptance.TestData) string {
	template := LinuxVirtualMachineScaleSetResource{}.identitySystemAssigned(data)
	return fmt.Sprintf(`
%s

data "azurerm_virtual_machine_scale_set" "test" {
  name                = azurerm_linux_virtual_machine_scale_set.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, template)
}

func (VirtualMachineScaleSetDataSource) basicWindows(data acceptance.TestData) string {
	template := WindowsVirtualMachineScaleSetResource{}.identitySystemAssigned(data)
	return fmt.Sprintf(`
%s

data "azurerm_virtual_machine_scale_set" "test" {
  name                = azurerm_windows_virtual_machine_scale_set.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, template)
}

func (VirtualMachineScaleSetDataSource) orchestrated(data acceptance.TestData) string {
	template := WindowsVirtualMachineResource{}.orchestratedZonal(data)
	return fmt.Sprintf(`
%s

data "azurerm_virtual_machine_scale_set" "test" {
  name                = azurerm_orchestrated_virtual_machine_scale_set.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, template)
}
