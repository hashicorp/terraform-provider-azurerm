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

func TestAccLinuxVirtualMachineScaleSet_resiliency_vmPoliciesOnly(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine_scale_set", "test")
	r := LinuxVirtualMachineScaleSetResource{}

	// TODO - Remove this override when resiliency policies are rolled out to all regions
	// Currently only supported in specific regions like EastUS2
	data.Locations.Primary = "eastus2"

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

func TestAccLinuxVirtualMachineScaleSet_resiliency_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine_scale_set", "test")
	r := LinuxVirtualMachineScaleSetResource{}

	// TODO - Remove this override when resiliency policies are rolled out to all regions
	// Currently only supported in specific regions like EastUS2
	data.Locations.Primary = "eastus2"

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

func TestAccLinuxVirtualMachineScaleSet_resiliency_vmCreationEnabledOnly(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine_scale_set", "test")
	r := LinuxVirtualMachineScaleSetResource{}

	// TODO - Remove this override when resiliency policies are rolled out to all regions
	// Currently only supported in specific regions like EastUS2
	data.Locations.Primary = "eastus2"

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

func TestAccLinuxVirtualMachineScaleSet_resiliency_vmDeletionEnabledOnly(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine_scale_set", "test")
	r := LinuxVirtualMachineScaleSetResource{}

	// TODO - Remove this override when resiliency policies are rolled out to all regions
	// Currently only supported in specific regions like EastUS2
	data.Locations.Primary = "eastus2"

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

func TestAccLinuxVirtualMachineScaleSet_resiliency_fieldsNotSetInState(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine_scale_set", "test")
	r := LinuxVirtualMachineScaleSetResource{}

	// TODO - Remove this override when resiliency policies are rolled out to all regions
	// Currently only supported in specific regions like EastUS2
	data.Locations.Primary = "eastus2"

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

func (r LinuxVirtualMachineScaleSetResource) resiliencyVMCreationOnly(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  name                            = "acctestvmss-%d"
  resource_group_name             = azurerm_resource_group.test.name
  location                        = azurerm_resource_group.test.location
  sku                             = "Standard_F2"
  instances                       = 1
  admin_username                  = "adminuser"
  disable_password_authentication = true

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  admin_ssh_key {
    username   = "adminuser"
    public_key = local.first_public_key
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

func (r LinuxVirtualMachineScaleSetResource) resiliencyVMDeletionOnly(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  name                            = "acctestvmss-%d"
  resource_group_name             = azurerm_resource_group.test.name
  location                        = azurerm_resource_group.test.location
  sku                             = "Standard_F2"
  instances                       = 1
  admin_username                  = "adminuser"
  disable_password_authentication = true

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  admin_ssh_key {
    username   = "adminuser"
    public_key = local.first_public_key
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

func (r LinuxVirtualMachineScaleSetResource) resiliencyVMPolicies(data acceptance.TestData, vmCreationEnabled, vmDeletionEnabled bool) string {
	return fmt.Sprintf(`
%s

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  name                            = "acctestvmss-%d"
  resource_group_name             = azurerm_resource_group.test.name
  location                        = azurerm_resource_group.test.location
  sku                             = "Standard_F2"
  instances                       = 1
  admin_username                  = "adminuser"
  disable_password_authentication = true

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  admin_ssh_key {
    username   = "adminuser"
    public_key = local.first_public_key
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

func (r LinuxVirtualMachineScaleSetResource) resiliencyFieldsNotConfigured(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  name                            = "acctestvmss-%d"
  resource_group_name             = azurerm_resource_group.test.name
  location                        = azurerm_resource_group.test.location
  sku                             = "Standard_F2"
  instances                       = 1
  admin_username                  = "adminuser"
  disable_password_authentication = true

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  admin_ssh_key {
    username   = "adminuser"
    public_key = local.first_public_key
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

  # Note: resilient_vm_creation_enabled and resilient_vm_deletion_enabled are intentionally NOT configured
  # This tests backward compatibility - these fields should not appear in state
}
`, r.template(data), data.RandomInteger)
}
