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
	data.Locations.Primary = "eastus2" // Resiliency policies are only supported in specific regions

	r := LinuxVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.resiliencyVMPolicies(data, true, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("admin_password"),
	})
}

func TestAccLinuxVirtualMachineScaleSet_resiliency_vmCreationOnly(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine_scale_set", "test")
	data.Locations.Primary = "eastus2" // Resiliency policies are only supported in specific regions
	r := LinuxVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.resiliencyVMPolicies(data, true, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("admin_password"),
	})
}

func TestAccLinuxVirtualMachineScaleSet_resiliency_vmDeletionOnly(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine_scale_set", "test")
	data.Locations.Primary = "eastus2" // Resiliency policies are only supported in specific regions
	r := LinuxVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.resiliencyVMPolicies(data, false, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("admin_password"),
	})
}

func TestAccLinuxVirtualMachineScaleSet_resiliency_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine_scale_set", "test")
	data.Locations.Primary = "eastus2" // Resiliency policies are only supported in specific regions

	r := LinuxVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.resiliencyVMPolicies(data, false, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("admin_password"),
		{
			Config: r.resiliencyVMPolicies(data, true, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("admin_password"),
		{
			Config:      r.resiliencyVMPolicies(data, false, false),
			ExpectError: regexp.MustCompile("Azure does not support disabling resiliency policies\\. Once the `resilient_vm.*_enabled` field is set to `true`, it cannot be reverted to `false`"),
		},
	})
}

func TestAccLinuxVirtualMachineScaleSet_resiliency_fieldsNotSetInState(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine_scale_set", "test")
	data.Locations.Primary = "eastus2" // Resiliency policies are only supported in specific regions

	r := LinuxVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.resiliencyFieldsNotConfigured(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("admin_password"),
	})
}

func TestAccLinuxVirtualMachineScaleSet_resiliency_explicitFalse(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine_scale_set", "test")
	data.Locations.Primary = "eastus2" // Resiliency policies are only supported in specific regions

	r := LinuxVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			// Test that explicit false values can be set and appear correctly in state
			Config: r.resiliencyVMPolicies(data, false, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("resilient_vm_creation_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("resilient_vm_deletion_enabled").HasValue("false"),
			),
		},
		// Import test with field exclusions: Azure API doesn't return resiliency fields when they're false,
		// treating them as "not configured". This is expected behavior, so we exclude these fields from
		// import state verification to avoid false failures.
		data.ImportStep("admin_password", "resilient_vm_creation_enabled", "resilient_vm_deletion_enabled"),
		{
			// Plan stability test: Verify that after Azure API interactions, applying the same config
			// with explicit false values doesn't produce unexpected diffs. This ensures the provider
			// correctly handles Azure's treatment of false values as "not configured".
			Config: r.resiliencyVMPolicies(data, false, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("resilient_vm_creation_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("resilient_vm_deletion_enabled").HasValue("false"),
			),
		},
	})
}

func TestAccLinuxVirtualMachineScaleSet_resiliency_unsupportedRegion(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine_scale_set", "test")
	data.Locations.Primary = "chilecentral" //Unsupported region

	r := LinuxVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.resiliencyVMPolicies(data, true, true),
			ExpectError: regexp.MustCompile("the resiliency policies.*are not supported in the.*chilecentral.*region"),
		},
	})
}

func (r LinuxVirtualMachineScaleSetResource) resiliencyVMPolicies(data acceptance.TestData, vmCreationEnabled, vmDeletionEnabled bool) string {
	return fmt.Sprintf(`
%s

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  name                            = "acctestvmss-%d"
  resource_group_name             = azurerm_resource_group.test.name
  location                        = azurerm_resource_group.test.location
  sku                             = "Standard_B1ls"
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
  sku                             = "Standard_B1ls"
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
