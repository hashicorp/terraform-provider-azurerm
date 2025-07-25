// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

func TestAccWindowsVirtualMachineScaleSet_resiliency_vmPoliciesOnly(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine_scale_set", "test")
	data.Locations.Primary = "eastus2" // Resiliency policies are only supported in specific regions

	r := WindowsVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.resiliencyVMPolicies(data, true, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("resilient_vm_creation_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("resilient_vm_deletion_enabled").HasValue("true"),
			),
		},
		data.ImportStep("admin_password"),
	})
}

func TestAccWindowsVirtualMachineScaleSet_resiliency_vmCreationOnly(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine_scale_set", "test")
	data.Locations.Primary = "eastus2" // Resiliency policies are only supported in specific regions

	r := WindowsVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.resiliencyVMPolicies(data, true, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("resilient_vm_creation_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("resilient_vm_deletion_enabled").HasValue("false"),
			),
		},
		data.ImportStep("admin_password"),
	})
}

func TestAccWindowsVirtualMachineScaleSet_resiliency_vmDeletionOnly(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine_scale_set", "test")
	data.Locations.Primary = "eastus2" // Resiliency policies are only supported in specific regions

	r := WindowsVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.resiliencyVMPolicies(data, false, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("resilient_vm_creation_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("resilient_vm_deletion_enabled").HasValue("true"),
			),
		},
		data.ImportStep("admin_password"),
	})
}

func TestAccWindowsVirtualMachineScaleSet_resiliency_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine_scale_set", "test")
	data.Locations.Primary = "eastus2" // Resiliency policies are only supported in specific regions

	r := WindowsVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.resiliencyVMPolicies(data, true, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("resilient_vm_creation_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("resilient_vm_deletion_enabled").HasValue("false"),
			),
		},
		data.ImportStep("admin_password"),
		{
			Config: r.resiliencyVMPolicies(data, true, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("resilient_vm_creation_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("resilient_vm_deletion_enabled").HasValue("true"),
			),
		},
		data.ImportStep("admin_password"),
		{
			Config:      r.resiliencyVMPolicies(data, false, false),
			ExpectError: regexp.MustCompile("Azure does not support disabling resiliency policies\\. Once the `resilient_vm.*_enabled` field is set to `true`, it cannot be reverted to `false`"),
		},
	})
}

func TestAccWindowsVirtualMachineScaleSet_resiliency_vmCreationEnabledOnly(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine_scale_set", "test")
	data.Locations.Primary = "eastus2" // Resiliency policies are only supported in specific regions

	r := WindowsVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.resiliencyVMCreationOnly(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("resilient_vm_creation_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("resilient_vm_deletion_enabled").HasValue("false"),
			),
		},
		data.ImportStep("admin_password"),
	})
}

func TestAccWindowsVirtualMachineScaleSet_resiliency_vmDeletionEnabledOnly(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine_scale_set", "test")
	data.Locations.Primary = "eastus2" // Resiliency policies are only supported in specific regions

	r := WindowsVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.resiliencyVMDeletionOnly(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("resilient_vm_creation_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("resilient_vm_deletion_enabled").HasValue("true"),
			),
		},
		data.ImportStep("admin_password"),
	})
}

func TestAccWindowsVirtualMachineScaleSet_resiliency_fieldsNotSetInState(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine_scale_set", "test")
	data.Locations.Primary = "eastus2" // Resiliency policies are only supported in specific regions

	r := WindowsVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.resiliencyFieldsNotConfigured(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				// Verify that resilient fields are NOT present in state when not configured
				check.That(data.ResourceName).Key("resilient_vm_creation_enabled").DoesNotExist(),
				check.That(data.ResourceName).Key("resilient_vm_deletion_enabled").DoesNotExist(),
			),
		},
		data.ImportStep("admin_password"),
	})
}

func (r WindowsVirtualMachineScaleSetResource) resiliencyVMCreationOnly(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_windows_virtual_machine_scale_set" "test" {
  name                 = "acctestvmss-%d"
  resource_group_name  = azurerm_resource_group.test.name
  location             = azurerm_resource_group.test.location
  sku                  = "Standard_B1ls"
  instances            = 1
  admin_username       = "adminuser"
  admin_password       = "P@55w0rd1234!"
  computer_name_prefix = "vm-"

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2022-Datacenter"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  network_interface {
    name    = "example"
    primary = true

    ip_configuration {
      name      = "internal"
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }

  resilient_vm_creation_enabled = true
}
`, r.template(data), data.RandomInteger)
}

func (r WindowsVirtualMachineScaleSetResource) resiliencyVMDeletionOnly(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_windows_virtual_machine_scale_set" "test" {
  name                 = "acctestvmss-%d"
  resource_group_name  = azurerm_resource_group.test.name
  location             = azurerm_resource_group.test.location
  sku                  = "Standard_B1ls"
  instances            = 1
  admin_username       = "adminuser"
  admin_password       = "P@55w0rd1234!"
  computer_name_prefix = "vm-"

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2022-Datacenter"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  network_interface {
    name    = "example"
    primary = true

    ip_configuration {
      name      = "internal"
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }

  resilient_vm_deletion_enabled = true
}
`, r.template(data), data.RandomInteger)
}

func (r WindowsVirtualMachineScaleSetResource) resiliencyVMPolicies(data acceptance.TestData, vmCreationEnabled, vmDeletionEnabled bool) string {
	return fmt.Sprintf(`
%s

resource "azurerm_windows_virtual_machine_scale_set" "test" {
  name                 = "acctestvmss-%d"
  resource_group_name  = azurerm_resource_group.test.name
  location             = azurerm_resource_group.test.location
  sku                  = "Standard_B1ls"
  instances            = 1
  admin_username       = "adminuser"
  admin_password       = "P@55w0rd1234!"
  computer_name_prefix = "vm-"

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2022-Datacenter"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  network_interface {
    name    = "example"
    primary = true

    ip_configuration {
      name      = "internal"
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }

  resilient_vm_creation_enabled = %t
  resilient_vm_deletion_enabled = %t
}
`, r.template(data), data.RandomInteger, vmCreationEnabled, vmDeletionEnabled)
}

func (r WindowsVirtualMachineScaleSetResource) resiliencyFieldsNotConfigured(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_windows_virtual_machine_scale_set" "test" {
  name                 = "acctestvmss-%d"
  resource_group_name  = azurerm_resource_group.test.name
  location             = azurerm_resource_group.test.location
  sku                  = "Standard_B1ls"
  instances            = 1
  admin_username       = "adminuser"
  admin_password       = "P@55w0rd1234!"
  computer_name_prefix = "vm-"

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2022-Datacenter"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  network_interface {
    name    = "example"
    primary = true

    ip_configuration {
      name      = "internal"
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }

  # NOTE: resilient_vm_creation_enabled and resilient_vm_deletion_enabled are intentionally
  # NOT configured here to test backward compatibility - they should not appear in state
}
`, r.template(data), data.RandomInteger)
}
