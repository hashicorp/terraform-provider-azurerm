// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type VirtualMachineDataSource struct{}

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
				check.That(data.ResourceName).Key("private_ip_address").HasValue("10.0.2.4"),
				check.That(data.ResourceName).Key("power_state").HasValue("running"),
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
				check.That(data.ResourceName).Key("private_ip_address").HasValue("10.0.2.4"),
				check.That(data.ResourceName).Key("power_state").HasValue("running"),
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
