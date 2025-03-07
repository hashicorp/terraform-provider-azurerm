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

func TestAccOrchestratedVirtualMachineScaleSet_priority(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_orchestrated_virtual_machine_scale_set", "test")
	r := OrchestratedVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.priorityTemplate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("priority").HasValue("Spot"),
				check.That(data.ResourceName).Key("eviction_policy").HasValue("Deallocate"),
			),
		},
	})
}

func TestAccOrchestratedVirtualMachineScaleSet_zoneUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_orchestrated_virtual_machine_scale_set", "test")
	r := OrchestratedVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.singleZone(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.multipleZone(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccOrchestratedVirtualMachineScaleSet_NonStandardCasing(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_orchestrated_virtual_machine_scale_set", "test")
	r := OrchestratedVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.nonStandardCasing(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:             r.nonStandardCasing(data),
			PlanOnly:           true,
			ExpectNonEmptyPlan: false,
		},
	})
}

func TestAccOrchestratedVirtualMachineScaleSet_AutomaticVMGuestPatchingLinux(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_orchestrated_virtual_machine_scale_set", "test")
	r := OrchestratedVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.linuxVMGuestPatching(data, "AutomaticByPlatform"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccOrchestratedVirtualMachineScaleSet_AutomaticVMGuestPatchingUpdateLinux(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_orchestrated_virtual_machine_scale_set", "test")
	r := OrchestratedVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.linuxVMGuestPatching(data, "ImageDefault"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("os_profile.0.linux_configuration.0.admin_password"),
		{
			Config: r.linuxVMGuestPatching(data, "AutomaticByPlatform"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("os_profile.0.linux_configuration.0.admin_password"),
		{
			Config: r.linuxVMGuestPatching(data, "ImageDefault"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccOrchestratedVirtualMachineScaleSet_PatchAssessmentModeLinux(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_orchestrated_virtual_machine_scale_set", "test")
	r := OrchestratedVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.otherPatchAssessmentModeLinuxDefault(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("os_profile.0.linux_configuration.0.patch_assessment_mode").HasValue("ImageDefault"),
			),
		},
		data.ImportStep("os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccOrchestratedVirtualMachineScaleSet_AutomaticVMGuestPatchingWindows(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_orchestrated_virtual_machine_scale_set", "test")
	r := OrchestratedVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.windowsVMGuestPatching(data, "AutomaticByPlatform"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("os_profile.0.windows_configuration.0.admin_password"),
	})
}

func TestAccOrchestratedVirtualMachineScaleSet_AutomaticVMGuestPatchingUpdateWindows(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_orchestrated_virtual_machine_scale_set", "test")
	r := OrchestratedVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.windowsVMGuestPatching(data, "AutomaticByOS"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("os_profile.0.windows_configuration.0.admin_password"),
		{
			Config: r.windowsVMGuestPatching(data, "AutomaticByPlatform"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("os_profile.0.windows_configuration.0.admin_password"),
		{
			Config: r.windowsVMGuestPatching(data, "AutomaticByOS"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("os_profile.0.windows_configuration.0.admin_password"),
	})
}

func TestAccOrchestratedVirtualMachineScaleSet_AutomaticVMGuestPatchingHotpatchingEnabledWindows(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_orchestrated_virtual_machine_scale_set", "test")
	r := OrchestratedVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.windowsHotpatchingEnabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("os_profile.0.windows_configuration.0.admin_password"),
		{
			Config: r.windowsHotpatchingEnabledUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("os_profile.0.windows_configuration.0.admin_password"),
	})
}

func TestAccOrchestratedVirtualMachineScaleSet_PatchAssessmentModeWindows(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_orchestrated_virtual_machine_scale_set", "test")
	r := OrchestratedVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.otherPatchAssessmentModeWindowsDefault(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("os_profile.0.windows_configuration.0.patch_assessment_mode").HasValue("ImageDefault"),
			),
		},
		data.ImportStep("os_profile.0.windows_configuration.0.admin_password"),
	})
}

func TestAccOrchestratedVirtualMachineScaleSet_otherAutomaticRepairsPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_orchestrated_virtual_machine_scale_set", "test")
	r := OrchestratedVirtualMachineScaleSetResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.otherAutomaticRepairsPolicy(data, "Restart"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("os_profile.0.linux_configuration.0.admin_password"),
		{
			Config: r.otherAutomaticRepairsPolicyDisabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("os_profile.0.linux_configuration.0.admin_password"),
		{
			Config: r.otherAutomaticRepairsPolicy(data, "Reimage"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("os_profile.0.linux_configuration.0.admin_password"),
		{
			Config: r.otherAutomaticRepairsPolicyDisabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("os_profile.0.linux_configuration.0.admin_password"),
		{
			Config: r.otherAutomaticRepairsPolicyEnabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccOrchestratedVirtualMachineScaleSet_otherUltraSsd(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_orchestrated_virtual_machine_scale_set", "test")
	r := OrchestratedVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.otherUltraSsd(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccOrchestratedVirtualMachineScaleSet_otherUpgradePolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_orchestrated_virtual_machine_scale_set", "test")
	r := OrchestratedVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.otherUpgradePolicy(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("os_profile.0.linux_configuration.0.admin_password"),
		{
			Config: r.otherUpgradePolicyUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccOrchestratedVirtualMachineScaleSet_otherRollingUpgrade(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_orchestrated_virtual_machine_scale_set", "test")
	r := OrchestratedVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.otherRollingUpgrade(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccOrchestratedVirtualMachineScaleSet_otherUpgradePolicyErrorValidation(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_orchestrated_virtual_machine_scale_set", "test")
	r := OrchestratedVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.otherRollingWithoutHealthExtension(data),
			ExpectError: regexp.MustCompile("a health extension must be specified when `upgrade_mode` is set to `Rolling`"),
		},
		{
			Config:      r.otherManualWithUpgradePolicy(data),
			ExpectError: regexp.MustCompile("a `rolling_upgrade_policy` block cannot be specified"),
		},
		{
			Config:      r.otherRollingWithoutUpgradePolicy(data),
			ExpectError: regexp.MustCompile("a `rolling_upgrade_policy` block must be specified when `upgrade_mode` is set"),
		},
	})
}

func TestAccOrchestratedVirtualMachineScaleSet_otherCapacityReservationGroupId(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_orchestrated_virtual_machine_scale_set", "test")
	r := OrchestratedVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.otherCapacityReservationGroupId(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccOrchestratedVirtualMachineScaleSet_otherVMAgentDisabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_orchestrated_virtual_machine_scale_set", "test")
	r := OrchestratedVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.otherVMAgentDisabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccOrchestratedVirtualMachineScaleSet_otherVMAgentDisabledWithExtensionDisabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_orchestrated_virtual_machine_scale_set", "test")
	r := OrchestratedVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.otherVMAgentDisabledWithExtensionDisabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccOrchestratedVirtualMachineScaleSet_otherUserData(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_orchestrated_virtual_machine_scale_set", "test")
	r := OrchestratedVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.otherUserData(data, "Hello World"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("os_profile.0.linux_configuration.0.admin_password"),
		{
			Config: r.otherUserData(data, "Goodbye World"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccOrchestratedVirtualMachineScaleSet_otherSinglePlacementGroup(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_orchestrated_virtual_machine_scale_set", "test")
	r := OrchestratedVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			// NOTE: Once single_placement_group has been set to false
			// it can no longer be changed back to true.
			Config: r.otherSinglePlacementGroup(data, "false"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccOrchestratedVirtualMachineScaleSet_otherSinglePlacementGroupInvalid(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_orchestrated_virtual_machine_scale_set", "test")
	r := OrchestratedVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.otherSinglePlacementGroup(data, "true"),
			ExpectError: regexp.MustCompile("Virtual Machine Scale Sets with Flexible Orchestration Mode or Virtual Machines referencing such Virtual Machine Scale Sets cannot use VM Size"),
		},
	})
}

func TestAccOrchestratedVirtualMachineScaleSet_otherSinglePlacementGroupNullMSeriesVmSku(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_orchestrated_virtual_machine_scale_set", "test")
	r := OrchestratedVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.otherSinglePlacementGroupNull(data, "Standard_M32ls"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccOrchestratedVirtualMachineScaleSet_otherSinglePlacementGroupNullStandardVmSku(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_orchestrated_virtual_machine_scale_set", "test")
	r := OrchestratedVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.otherSinglePlacementGroupNull(data, "Standard_D1_v2"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("os_profile.0.linux_configuration.0.admin_password"),
	})
}

func (OrchestratedVirtualMachineScaleSetResource) priorityTemplate(data acceptance.TestData) string {
	r := OrchestratedVirtualMachineScaleSetResource{}
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-OVMSS-%[1]d"
  location = "%[2]s"
}

%[3]s

resource "azurerm_orchestrated_virtual_machine_scale_set" "test" {
  name                = "acctestOVMSS-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  priority        = "Spot"
  eviction_policy = "Deallocate"

  sku_name  = "Standard_D1_v2"
  instances = 1

  platform_fault_domain_count = 2

  os_profile {
    linux_configuration {
      computer_name_prefix = "testvm"
      admin_username       = "myadmin"
      admin_password       = "Passwword1234"

      disable_password_authentication = false
    }
  }

  network_interface {
    name    = "TestNetworkProfile"
    primary = true

    ip_configuration {
      name      = "TestIPConfiguration"
      primary   = true
      subnet_id = azurerm_subnet.test.id

      public_ip_address {
        name                    = "TestPublicIPConfiguration"
        domain_name_label       = "test-domain-label"
        idle_timeout_in_minutes = 4
      }
    }
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }
}
`, data.RandomInteger, data.Locations.Primary, r.natgateway_template(data))
}

func (OrchestratedVirtualMachineScaleSetResource) singleZone(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-OVMSS-%[1]d"
  location = "%[2]s"
}

resource "azurerm_orchestrated_virtual_machine_scale_set" "test" {
  name                = "acctestOVMSS-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  platform_fault_domain_count = 1

  zones = ["1"]

  tags = {
    ENV = "Test"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (OrchestratedVirtualMachineScaleSetResource) multipleZone(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-OVMSS-%[1]d"
  location = "%[2]s"
}

resource "azurerm_orchestrated_virtual_machine_scale_set" "test" {
  name                = "acctestOVMSS-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  platform_fault_domain_count = 1

  zones = ["1", "2"]

  tags = {
    ENV = "Test"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (OrchestratedVirtualMachineScaleSetResource) windowsHotpatchingEnabled(data acceptance.TestData) string {
	r := OrchestratedVirtualMachineScaleSetResource{}
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-OVMSS-%[1]d"
  location = "%[2]s"
}

%[3]s

resource "azurerm_orchestrated_virtual_machine_scale_set" "test" {
  name                = "acctestOVMSS-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku_name = "Standard_F2s_v2"

  # Orchestrated VMSS allocation will timeout at service side due to extension, set instances to 0 to avoid the timeout
  instances = 0

  platform_fault_domain_count = 2

  os_profile {
    windows_configuration {
      computer_name_prefix = "testvm"
      admin_username       = "myadmin"
      admin_password       = "Passwword1234"

      patch_mode          = "AutomaticByPlatform"
      hotpatching_enabled = true
    }
  }

  network_interface {
    name    = "TestNetworkProfile"
    primary = true

    ip_configuration {
      name      = "TestIPConfiguration"
      primary   = true
      subnet_id = azurerm_subnet.test.id

      public_ip_address {
        name                    = "TestPublicIPConfiguration"
        domain_name_label       = "test-domain-label"
        idle_timeout_in_minutes = 4
      }
    }
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2022-datacenter-azure-edition-core"
    version   = "latest"
  }

  extension {
    name                               = "HealthExtension"
    publisher                          = "Microsoft.ManagedServices"
    type                               = "ApplicationHealthWindows"
    type_handler_version               = "1.0"
    auto_upgrade_minor_version_enabled = true

    settings = jsonencode({
      "protocol"    = "http"
      "port"        = "80"
      "requestPath" = "/"
    })
  }
}
`, data.RandomInteger, data.Locations.Primary, r.natgateway_template(data))
}

func (OrchestratedVirtualMachineScaleSetResource) windowsHotpatchingEnabledUpdate(data acceptance.TestData) string {
	r := OrchestratedVirtualMachineScaleSetResource{}
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-OVMSS-%[1]d"
  location = "%[2]s"
}

%[3]s

resource "azurerm_orchestrated_virtual_machine_scale_set" "test" {
  name                = "acctestOVMSS-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku_name = "Standard_F2s_v2"

  # Orchestrated VMSS allocation will timeout at service side due to extension, set instances to 0 to avoid the timeout
  instances = 0

  platform_fault_domain_count = 2

  os_profile {
    windows_configuration {
      computer_name_prefix = "testvm"
      admin_username       = "myadmin"
      admin_password       = "Passwword1234"

      patch_mode          = "AutomaticByPlatform"
      hotpatching_enabled = true
    }
  }

  network_interface {
    name    = "TestNetworkProfile"
    primary = true

    ip_configuration {
      name      = "TestIPConfiguration"
      primary   = true
      subnet_id = azurerm_subnet.test.id

      public_ip_address {
        name                    = "TestPublicIPConfiguration"
        domain_name_label       = "test-domain-label"
        idle_timeout_in_minutes = 4
      }
    }
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2022-datacenter-azure-edition-core"
    version   = "latest"
  }

  extension {
    name                               = "HealthExtension"
    publisher                          = "Microsoft.ManagedServices"
    type                               = "ApplicationHealthWindows"
    type_handler_version               = "1.0"
    auto_upgrade_minor_version_enabled = true

    settings = jsonencode({
      "protocol"    = "http"
      "port"        = "80"
      "requestPath" = "/"
    })
  }

  tags = {
    value = "update"
  }
}
`, data.RandomInteger, data.Locations.Primary, r.natgateway_template(data))
}

func (OrchestratedVirtualMachineScaleSetResource) windowsVMGuestPatching(data acceptance.TestData, patchMode string) string {
	r := OrchestratedVirtualMachineScaleSetResource{}
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-OVMSS-%[1]d"
  location = "%[2]s"
}

%[4]s

resource "azurerm_orchestrated_virtual_machine_scale_set" "test" {
  name                = "acctestOVMSS-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku_name = "Standard_D1_v2"

  # Orchestrated VMSS allocation will timeout at service side due to extension, set instances to 0 to avoid the timeout
  instances = 0

  platform_fault_domain_count = 2

  os_profile {
    windows_configuration {
      computer_name_prefix = "testvm"
      admin_username       = "myadmin"
      admin_password       = "Passwword1234"

      patch_mode = "%[3]s"
    }
  }

  network_interface {
    name    = "TestNetworkProfile"
    primary = true

    ip_configuration {
      name      = "TestIPConfiguration"
      primary   = true
      subnet_id = azurerm_subnet.test.id

      public_ip_address {
        name                    = "TestPublicIPConfiguration"
        domain_name_label       = "test-domain-label"
        idle_timeout_in_minutes = 4
      }
    }
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter"
    version   = "latest"
  }

  extension {
    name                               = "HealthExtension"
    publisher                          = "Microsoft.ManagedServices"
    type                               = "ApplicationHealthWindows"
    type_handler_version               = "1.0"
    auto_upgrade_minor_version_enabled = true

    settings = jsonencode({
      "protocol"    = "http"
      "port"        = "80"
      "requestPath" = "/healthEndpoint"
    })
  }
}
`, data.RandomInteger, data.Locations.Primary, patchMode, r.natgateway_template(data))
}

func (OrchestratedVirtualMachineScaleSetResource) otherPatchAssessmentModeWindowsDefault(data acceptance.TestData) string {
	r := OrchestratedVirtualMachineScaleSetResource{}
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-OVMSS-%[1]d"
  location = "%[2]s"
}

%[3]s

resource "azurerm_orchestrated_virtual_machine_scale_set" "test" {
  name                = "acctestOVMSS-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku_name = "Standard_D1_v2"

  # Orchestrated VMSS allocation will timeout at service side due to extension, set instances to 0 to avoid the timeout
  instances = 0

  platform_fault_domain_count = 2

  os_profile {
    windows_configuration {
      computer_name_prefix = "testvm"
      admin_username       = "myadmin"
      admin_password       = "Passwword1234"
    }
  }

  network_interface {
    name    = "TestNetworkProfile"
    primary = true

    ip_configuration {
      name      = "TestIPConfiguration"
      primary   = true
      subnet_id = azurerm_subnet.test.id

      public_ip_address {
        name                    = "TestPublicIPConfiguration"
        domain_name_label       = "test-domain-label"
        idle_timeout_in_minutes = 4
      }
    }
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter"
    version   = "latest"
  }

  extension {
    name                               = "HealthExtension"
    publisher                          = "Microsoft.ManagedServices"
    type                               = "ApplicationHealthWindows"
    type_handler_version               = "1.0"
    auto_upgrade_minor_version_enabled = true

    settings = jsonencode({
      "protocol"    = "http"
      "port"        = "80"
      "requestPath" = "/healthEndpoint"
    })
  }
}
`, data.RandomInteger, data.Locations.Primary, r.natgateway_template(data))
}

func (OrchestratedVirtualMachineScaleSetResource) linuxVMGuestPatching(data acceptance.TestData, patchMode string) string {
	r := OrchestratedVirtualMachineScaleSetResource{}
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-OVMSS-%[1]d"
  location = "%[2]s"
}

%[4]s

resource "azurerm_orchestrated_virtual_machine_scale_set" "test" {
  name                = "acctestOVMSS-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku_name = "Standard_F2"

  # Orchestrated VMSS allocation will timeout at service side due to extension, set instances to 0 to avoid the timeout
  instances = 0

  platform_fault_domain_count = 2

  os_profile {
    linux_configuration {
      computer_name_prefix = "testvm"
      admin_username       = "myadmin"

      admin_ssh_key {
        username   = "myadmin"
        public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDvXYZAjVUt2aojUV3XIA+PY6gXrgbvktXwf2NoIHGlQFhogpMEyOfqgogCtTBM7MNCS3ELul6SV+mlpH08Ki45ADIQuDXdommCvsMFW096JrsHOJpGfjCsJ1gbbv7brB3Ag+BSGb4qO3pRsEVTtZCeJDwfH5D7vmqP5xXcELKR4UAtKQKUhLvt6mhW90sFLTJeOTiYGbavIKqfCUFSeSMQkUPr8o3uzOfeWyCw7tc7szLuvfwJ5poGHuve73KKAlUnDTPUrhyj7iITZSDl+/i+bpDzPyCyJWDMsC0ON7q2fDr2mEz0L9ACrsI5Nx3lt5fe+IaHSrjivqnL8SqUWSN45o9Qp99sGWFiuTfos8f1jp+AXzC4ArVtKyRg/CnzKRiK0CGSxBJ5s9zAoa7yBBmjCszq89vFa0eMgpEIZFwa6kKJKt9AfRBXgO9YGPV4uaN7topy92/p2pE+vF8IafarbvnTDOQt62mS07tXYqYg1DhecrmBVWKlq9oafBweoeTjoq52SoGsuDc/YAOzIgWVIuvV8yKoh9KbXPWowjLtxDhRIS/d1nMMNdNI8X0TQivgi5+umMgAXhsVAKSNDUauLt4jimYkWAuE+R6KoCqVFdaB9bQDySBjAziruDSe3reToydjzzluvHMjWK8QiDynxs41pi4zZz6gAlca3QPkEQ== hello@world.com"
      }

      patch_mode = "%[3]s"
    }
  }

  network_interface {
    name    = "TestNetworkProfile"
    primary = true

    ip_configuration {
      name      = "TestIPConfiguration"
      primary   = true
      subnet_id = azurerm_subnet.test.id

      public_ip_address {
        name                    = "TestPublicIPConfiguration"
        domain_name_label       = "test-domain-label"
        idle_timeout_in_minutes = 4
      }
    }
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }

  extension {
    name                               = "HealthExtension"
    publisher                          = "Microsoft.ManagedServices"
    type                               = "ApplicationHealthLinux"
    type_handler_version               = "1.0"
    auto_upgrade_minor_version_enabled = true

    settings = jsonencode({
      "protocol"    = "http"
      "port"        = 80
      "requestPath" = "/healthEndpoint"
    })
  }
}
`, data.RandomInteger, data.Locations.Primary, patchMode, r.natgateway_template(data))
}

func (OrchestratedVirtualMachineScaleSetResource) otherPatchAssessmentModeLinuxDefault(data acceptance.TestData) string {
	r := OrchestratedVirtualMachineScaleSetResource{}
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-OVMSS-%[1]d"
  location = "%[2]s"
}

%[3]s

resource "azurerm_orchestrated_virtual_machine_scale_set" "test" {
  name                = "acctestOVMSS-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku_name = "Standard_F2"

  # Orchestrated VMSS allocation will timeout at service side due to extension, set instances to 0 to avoid the timeout
  instances = 0

  platform_fault_domain_count = 2

  os_profile {
    linux_configuration {
      computer_name_prefix = "testvm"
      admin_username       = "myadmin"

      admin_ssh_key {
        username   = "myadmin"
        public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDvXYZAjVUt2aojUV3XIA+PY6gXrgbvktXwf2NoIHGlQFhogpMEyOfqgogCtTBM7MNCS3ELul6SV+mlpH08Ki45ADIQuDXdommCvsMFW096JrsHOJpGfjCsJ1gbbv7brB3Ag+BSGb4qO3pRsEVTtZCeJDwfH5D7vmqP5xXcELKR4UAtKQKUhLvt6mhW90sFLTJeOTiYGbavIKqfCUFSeSMQkUPr8o3uzOfeWyCw7tc7szLuvfwJ5poGHuve73KKAlUnDTPUrhyj7iITZSDl+/i+bpDzPyCyJWDMsC0ON7q2fDr2mEz0L9ACrsI5Nx3lt5fe+IaHSrjivqnL8SqUWSN45o9Qp99sGWFiuTfos8f1jp+AXzC4ArVtKyRg/CnzKRiK0CGSxBJ5s9zAoa7yBBmjCszq89vFa0eMgpEIZFwa6kKJKt9AfRBXgO9YGPV4uaN7topy92/p2pE+vF8IafarbvnTDOQt62mS07tXYqYg1DhecrmBVWKlq9oafBweoeTjoq52SoGsuDc/YAOzIgWVIuvV8yKoh9KbXPWowjLtxDhRIS/d1nMMNdNI8X0TQivgi5+umMgAXhsVAKSNDUauLt4jimYkWAuE+R6KoCqVFdaB9bQDySBjAziruDSe3reToydjzzluvHMjWK8QiDynxs41pi4zZz6gAlca3QPkEQ== hello@world.com"
      }
    }
  }

  network_interface {
    name    = "TestNetworkProfile"
    primary = true

    ip_configuration {
      name      = "TestIPConfiguration"
      primary   = true
      subnet_id = azurerm_subnet.test.id

      public_ip_address {
        name                    = "TestPublicIPConfiguration"
        domain_name_label       = "test-domain-label"
        idle_timeout_in_minutes = 4
      }
    }
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }

  extension {
    name                               = "HealthExtension"
    publisher                          = "Microsoft.ManagedServices"
    type                               = "ApplicationHealthLinux"
    type_handler_version               = "1.0"
    auto_upgrade_minor_version_enabled = true

    settings = jsonencode({
      "protocol"    = "http"
      "port"        = 80
      "requestPath" = "/healthEndpoint"
    })
  }
}
`, data.RandomInteger, data.Locations.Primary, r.natgateway_template(data))
}

func (OrchestratedVirtualMachineScaleSetResource) nonStandardCasing(data acceptance.TestData) string {
	r := OrchestratedVirtualMachineScaleSetResource{}
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-OVMSS-%[1]d"
  location = "%[2]s"
}

%[3]s

resource "azurerm_orchestrated_virtual_machine_scale_set" "test" {
  name                = "acctestOVMSS-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku_name  = "Standard_F2"
  instances = 2

  platform_fault_domain_count = 2

  os_profile {
    linux_configuration {
      computer_name_prefix = "testvm"
      admin_username       = "myadmin"
      admin_password       = "Passwword1234"

      disable_password_authentication = false
    }
  }

  network_interface {
    name    = "TestNetworkProfile-%[1]d"
    primary = true

    ip_configuration {
      name      = "TestIPConfiguration"
      primary   = true
      subnet_id = azurerm_subnet.test.id

      public_ip_address {
        name                    = "TestPublicIPConfiguration"
        domain_name_label       = "test-domain-label"
        idle_timeout_in_minutes = 4
      }
    }
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }
}
`, data.RandomInteger, data.Locations.Primary, r.natgateway_template(data))
}

func (OrchestratedVirtualMachineScaleSetResource) otherAutomaticRepairsPolicyEnabled(data acceptance.TestData) string {
	r := OrchestratedVirtualMachineScaleSetResource{}
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-OVMSS-%[1]d"
  location = "%[2]s"
}

%[3]s

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  name = "acctest%[4]s"
}

resource "azurerm_orchestrated_virtual_machine_scale_set" "test" {
  name                = "acctestOVMSS-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  zones = []

  sku_name = "Standard_D1_v2"

  # Orchestrated VMSS allocation will timeout at service side due to extension, set instances to 0 to avoid the timeout
  instances = 0

  platform_fault_domain_count = 2

  os_profile {
    linux_configuration {
      computer_name_prefix = "testvm-%[1]d"
      admin_username       = "myadmin"
      admin_password       = "Passwword1234"

      disable_password_authentication = false
    }
  }

  network_interface {
    name    = "TestNetworkProfile-%[1]d"
    primary = true

    ip_configuration {
      name      = "TestIPConfiguration"
      primary   = true
      subnet_id = azurerm_subnet.test.id

      public_ip_address {
        name                    = "TestPublicIPConfiguration"
        domain_name_label       = "test-domain-label"
        idle_timeout_in_minutes = 4
      }
    }
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }

  automatic_instance_repair {
    enabled      = true
    grace_period = "PT30M"
  }

  extension {
    name                               = "HealthExtension"
    publisher                          = "Microsoft.ManagedServices"
    type                               = "ApplicationHealthLinux"
    type_handler_version               = "1.0"
    auto_upgrade_minor_version_enabled = true

    settings = jsonencode({
      "protocol"    = "http"
      "port"        = 80
      "requestPath" = "/healthEndpoint"
    })
  }
}
`, data.RandomInteger, data.Locations.Primary, r.natgateway_template(data), data.RandomString)
}

func (OrchestratedVirtualMachineScaleSetResource) otherAutomaticRepairsPolicyDisabled(data acceptance.TestData) string {
	r := OrchestratedVirtualMachineScaleSetResource{}
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-OVMSS-%[1]d"
  location = "%[2]s"
}

%[3]s

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  name = "acctest%[4]s"
}

resource "azurerm_orchestrated_virtual_machine_scale_set" "test" {
  name                = "acctestOVMSS-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  zones = []

  sku_name = "Standard_D1_v2"

  # Orchestrated VMSS allocation will timeout at service side due to extension, set instances to 0 to avoid the timeout
  instances = 0

  platform_fault_domain_count = 2

  os_profile {
    linux_configuration {
      computer_name_prefix = "testvm-%[1]d"
      admin_username       = "myadmin"
      admin_password       = "Passwword1234"

      disable_password_authentication = false
    }
  }

  network_interface {
    name    = "TestNetworkProfile-%[1]d"
    primary = true

    ip_configuration {
      name      = "TestIPConfiguration"
      primary   = true
      subnet_id = azurerm_subnet.test.id

      public_ip_address {
        name                    = "TestPublicIPConfiguration"
        domain_name_label       = "test-domain-label"
        idle_timeout_in_minutes = 4
      }
    }
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }

  automatic_instance_repair {
    enabled = false
  }

  extension {
    name                               = "HealthExtension"
    publisher                          = "Microsoft.ManagedServices"
    type                               = "ApplicationHealthLinux"
    type_handler_version               = "1.0"
    auto_upgrade_minor_version_enabled = true

    settings = jsonencode({
      "protocol"    = "http"
      "port"        = 80
      "requestPath" = "/healthEndpoint"
    })
  }
}
`, data.RandomInteger, data.Locations.Primary, r.natgateway_template(data), data.RandomString)
}

func (OrchestratedVirtualMachineScaleSetResource) otherAutomaticRepairsPolicy(data acceptance.TestData, action string) string {
	r := OrchestratedVirtualMachineScaleSetResource{}
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-OVMSS-%[1]d"
  location = "%[2]s"
}
%[3]s

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  name                = "acctest%[4]s"
}

resource "azurerm_orchestrated_virtual_machine_scale_set" "test" {
  name                = "acctestOVMSS-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  zones               = []
  sku_name            = "Standard_D1_v2"
  # Orchestrated VMSS allocation will timeout at service side due to extension, set instances to 0 to avoid the timeout
  instances                   = 0
  platform_fault_domain_count = 2

  os_profile {
    linux_configuration {
      computer_name_prefix            = "testvm-%[1]d"
      admin_username                  = "myadmin"
      admin_password                  = "Passwword1234"
      disable_password_authentication = false
    }
  }

  network_interface {
    name    = "TestNetworkProfile-%[1]d"
    primary = true
    ip_configuration {
      name      = "TestIPConfiguration"
      primary   = true
      subnet_id = azurerm_subnet.test.id
      public_ip_address {
        name                    = "TestPublicIPConfiguration"
        domain_name_label       = "test-domain-label"
        idle_timeout_in_minutes = 4
      }
    }
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  automatic_instance_repair {
    enabled      = true
    grace_period = "PT30M"
    action       = "%[5]s"
  }

  extension {
    name                               = "HealthExtension"
    publisher                          = "Microsoft.ManagedServices"
    type                               = "ApplicationHealthLinux"
    type_handler_version               = "1.0"
    auto_upgrade_minor_version_enabled = true
    settings = jsonencode({
      "protocol"    = "http"
      "port"        = 80
      "requestPath" = "/healthEndpoint"
    })
  }
}
`, data.RandomInteger, data.Locations.Primary, r.natgateway_template(data), data.RandomString, action)
}

func (OrchestratedVirtualMachineScaleSetResource) otherUltraSsd(data acceptance.TestData) string {
	r := OrchestratedVirtualMachineScaleSetResource{}
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-OVMSS-%[1]d"
  location = "%[2]s"
}

%[3]s

resource "azurerm_orchestrated_virtual_machine_scale_set" "test" {
  name                = "acctestOVMSS-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  zones = ["1"]

  sku_name  = "Standard_D2s_v3"
  instances = 1

  platform_fault_domain_count = 1

  os_profile {
    linux_configuration {
      computer_name_prefix = "testvm-%[1]d"
      admin_username       = "myadmin"
      admin_password       = "Passwword1234"

      disable_password_authentication = false
    }
  }

  network_interface {
    name    = "TestNetworkProfile-%[1]d"
    primary = true

    ip_configuration {
      name      = "TestIPConfiguration"
      primary   = true
      subnet_id = azurerm_subnet.test.id

      public_ip_address {
        name                    = "TestPublicIPConfiguration"
        domain_name_label       = "test-domain-label"
        idle_timeout_in_minutes = 4
      }
    }
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }

  additional_capabilities {
    ultra_ssd_enabled = true
  }
}
`, data.RandomInteger, data.Locations.Primary, r.natgateway_template(data))
}

func (OrchestratedVirtualMachineScaleSetResource) otherUserData(data acceptance.TestData, userData string) string {
	r := OrchestratedVirtualMachineScaleSetResource{}
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-OVMSS-%[1]d"
  location = "%[2]s"
}

%[3]s

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  name = "acctest%[4]s"
}

resource "azurerm_orchestrated_virtual_machine_scale_set" "test" {
  name                = "acctestOVMSS-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  zones = []

  sku_name  = "Standard_D1_v2"
  instances = 1

  platform_fault_domain_count = 2

  os_profile {
    linux_configuration {
      computer_name_prefix = "testvm-%[1]d"
      admin_username       = "myadmin"
      admin_password       = "Passwword1234"

      disable_password_authentication = false
    }
  }

  network_interface {
    name    = "TestNetworkProfile-%[1]d"
    primary = true

    ip_configuration {
      name      = "TestIPConfiguration"
      primary   = true
      subnet_id = azurerm_subnet.test.id

      public_ip_address {
        name                    = "TestPublicIPConfiguration"
        domain_name_label       = "test-domain-label"
        idle_timeout_in_minutes = 4
      }
    }
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }

  user_data_base64 = base64encode("%[5]s")

  tags = {
    "platformsettings.host_environment.service.platform_optedin_for_rootcerts" = "true"
  }
}
`, data.RandomInteger, data.Locations.Primary, r.natgateway_template(data), data.RandomString, userData)
}

func (OrchestratedVirtualMachineScaleSetResource) otherUpgradePolicy(data acceptance.TestData) string {
	r := OrchestratedVirtualMachineScaleSetResource{}
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-OVMSS-%[1]d"
  location = "%[2]s"
}

%[3]s

resource "azurerm_orchestrated_virtual_machine_scale_set" "test" {
  name                        = "acctestOVMSS-%[1]d"
  location                    = azurerm_resource_group.test.location
  resource_group_name         = azurerm_resource_group.test.name
  sku_name                    = "Standard_D1_v2"
  instances                   = 2
  platform_fault_domain_count = 1
  upgrade_mode                = "Automatic"
  zones                       = ["1"]

  rolling_upgrade_policy {
    cross_zone_upgrades_enabled             = true
    max_batch_instance_percent              = 10
    max_unhealthy_instance_percent          = 10
    max_unhealthy_upgraded_instance_percent = 10
    pause_time_between_batches              = "PT0S"
    prioritize_unhealthy_instances_enabled  = true
    maximum_surge_instances_enabled         = false
  }

  os_profile {
    linux_configuration {
      computer_name_prefix            = "testvm-%[1]d"
      admin_username                  = "myadmin"
      admin_password                  = "Passwword1234"
      disable_password_authentication = false
    }
  }

  network_interface {
    name    = "TestNetworkProfile-%[1]d"
    primary = true

    ip_configuration {
      name      = "TestIPConfiguration"
      primary   = true
      subnet_id = azurerm_subnet.test.id

      public_ip_address {
        name                    = "TestPublicIPConfiguration"
        domain_name_label       = "test-domain-label"
        idle_timeout_in_minutes = 4
      }
    }
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }
}
`, data.RandomInteger, data.Locations.Primary, r.natgateway_template(data))
}

func (OrchestratedVirtualMachineScaleSetResource) otherUpgradePolicyUpdate(data acceptance.TestData) string {
	r := OrchestratedVirtualMachineScaleSetResource{}
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-OVMSS-%[1]d"
  location = "%[2]s"
}

%[3]s

resource "azurerm_orchestrated_virtual_machine_scale_set" "test" {
  name                        = "acctestOVMSS-%[1]d"
  location                    = azurerm_resource_group.test.location
  resource_group_name         = azurerm_resource_group.test.name
  sku_name                    = "Standard_D1_v2"
  instances                   = 2
  platform_fault_domain_count = 1
  upgrade_mode                = "Automatic"
  zones                       = ["1"]

  rolling_upgrade_policy {
    cross_zone_upgrades_enabled             = false
    max_batch_instance_percent              = 20
    max_unhealthy_instance_percent          = 20
    max_unhealthy_upgraded_instance_percent = 20
    pause_time_between_batches              = "PT1S"
    prioritize_unhealthy_instances_enabled  = false
    maximum_surge_instances_enabled         = true
  }

  os_profile {
    linux_configuration {
      computer_name_prefix            = "testvm-%[1]d"
      admin_username                  = "myadmin"
      admin_password                  = "Passwword1234"
      disable_password_authentication = false
    }
  }

  network_interface {
    name    = "TestNetworkProfile-%[1]d"
    primary = true

    ip_configuration {
      name      = "TestIPConfiguration"
      primary   = true
      subnet_id = azurerm_subnet.test.id

      public_ip_address {
        name                    = "TestPublicIPConfiguration"
        domain_name_label       = "test-domain-label"
        idle_timeout_in_minutes = 4
      }
    }
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }
}
`, data.RandomInteger, data.Locations.Primary, r.natgateway_template(data))
}

func (OrchestratedVirtualMachineScaleSetResource) otherRollingUpgrade(data acceptance.TestData) string {
	r := OrchestratedVirtualMachineScaleSetResource{}
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-OVMSS-%[1]d"
  location = "%[2]s"
}

%[3]s

resource "azurerm_orchestrated_virtual_machine_scale_set" "test" {
  name                        = "acctestOVMSS-%[1]d"
  location                    = azurerm_resource_group.test.location
  resource_group_name         = azurerm_resource_group.test.name
  sku_name                    = "Standard_D1_v2"
  instances                   = 2
  platform_fault_domain_count = 1
  upgrade_mode                = "Rolling"
  zones                       = ["1"]

  rolling_upgrade_policy {
    cross_zone_upgrades_enabled             = true
    max_batch_instance_percent              = 10
    max_unhealthy_instance_percent          = 10
    max_unhealthy_upgraded_instance_percent = 10
    pause_time_between_batches              = "PT0S"
    prioritize_unhealthy_instances_enabled  = true
    maximum_surge_instances_enabled         = false
  }

  os_profile {
    linux_configuration {
      computer_name_prefix            = "testvm-%[1]d"
      admin_username                  = "myadmin"
      admin_password                  = "Passwword1234"
      disable_password_authentication = false
    }
  }

  network_interface {
    name    = "TestNetworkProfile-%[1]d"
    primary = true

    ip_configuration {
      name      = "TestIPConfiguration"
      primary   = true
      subnet_id = azurerm_subnet.test.id

      public_ip_address {
        name                    = "TestPublicIPConfiguration"
        domain_name_label       = "test-domain-label"
        idle_timeout_in_minutes = 4
      }
    }
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }

  extension {
    name                 = "HealthExtension"
    publisher            = "Microsoft.ManagedServices"
    type                 = "ApplicationHealthLinux"
    type_handler_version = "1.0"
    settings = jsonencode({
      protocol = "https"
      port     = 443
    })
  }
}
`, data.RandomInteger, data.Locations.Primary, r.natgateway_template(data))
}

func (OrchestratedVirtualMachineScaleSetResource) otherRollingWithoutHealthExtension(data acceptance.TestData) string {
	r := OrchestratedVirtualMachineScaleSetResource{}
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-OVMSS-%[1]d"
  location = "%[2]s"
}

%[3]s

resource "azurerm_orchestrated_virtual_machine_scale_set" "test" {
  name                        = "acctestOVMSS-%[1]d"
  location                    = azurerm_resource_group.test.location
  resource_group_name         = azurerm_resource_group.test.name
  sku_name                    = "Standard_D1_v2"
  instances                   = 2
  platform_fault_domain_count = 1
  upgrade_mode                = "Rolling"
  zones                       = ["1"]

  rolling_upgrade_policy {
    cross_zone_upgrades_enabled             = true
    max_batch_instance_percent              = 10
    max_unhealthy_instance_percent          = 10
    max_unhealthy_upgraded_instance_percent = 10
    pause_time_between_batches              = "PT0S"
    prioritize_unhealthy_instances_enabled  = true
    maximum_surge_instances_enabled         = false
  }

  os_profile {
    linux_configuration {
      computer_name_prefix            = "testvm-%[1]d"
      admin_username                  = "myadmin"
      admin_password                  = "Passwword1234"
      disable_password_authentication = false
    }
  }

  network_interface {
    name    = "TestNetworkProfile-%[1]d"
    primary = true

    ip_configuration {
      name      = "TestIPConfiguration"
      primary   = true
      subnet_id = azurerm_subnet.test.id

      public_ip_address {
        name                    = "TestPublicIPConfiguration"
        domain_name_label       = "test-domain-label"
        idle_timeout_in_minutes = 4
      }
    }
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }
}
`, data.RandomInteger, data.Locations.Primary, r.natgateway_template(data))
}

func (OrchestratedVirtualMachineScaleSetResource) otherManualWithUpgradePolicy(data acceptance.TestData) string {
	r := OrchestratedVirtualMachineScaleSetResource{}
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-OVMSS-%[1]d"
  location = "%[2]s"
}

%[3]s

resource "azurerm_orchestrated_virtual_machine_scale_set" "test" {
  name                        = "acctestOVMSS-%[1]d"
  location                    = azurerm_resource_group.test.location
  resource_group_name         = azurerm_resource_group.test.name
  sku_name                    = "Standard_D1_v2"
  instances                   = 2
  platform_fault_domain_count = 1
  upgrade_mode                = "Manual"
  zones                       = ["1"]

  rolling_upgrade_policy {
    cross_zone_upgrades_enabled             = true
    max_batch_instance_percent              = 10
    max_unhealthy_instance_percent          = 10
    max_unhealthy_upgraded_instance_percent = 10
    pause_time_between_batches              = "PT0S"
    prioritize_unhealthy_instances_enabled  = true
    maximum_surge_instances_enabled         = false
  }

  os_profile {
    linux_configuration {
      computer_name_prefix            = "testvm-%[1]d"
      admin_username                  = "myadmin"
      admin_password                  = "Passwword1234"
      disable_password_authentication = false
    }
  }

  network_interface {
    name    = "TestNetworkProfile-%[1]d"
    primary = true

    ip_configuration {
      name      = "TestIPConfiguration"
      primary   = true
      subnet_id = azurerm_subnet.test.id

      public_ip_address {
        name                    = "TestPublicIPConfiguration"
        domain_name_label       = "test-domain-label"
        idle_timeout_in_minutes = 4
      }
    }
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }
}
`, data.RandomInteger, data.Locations.Primary, r.natgateway_template(data))
}

func (OrchestratedVirtualMachineScaleSetResource) otherRollingWithoutUpgradePolicy(data acceptance.TestData) string {
	r := OrchestratedVirtualMachineScaleSetResource{}
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-OVMSS-%[1]d"
  location = "%[2]s"
}

%[3]s

resource "azurerm_orchestrated_virtual_machine_scale_set" "test" {
  name                        = "acctestOVMSS-%[1]d"
  location                    = azurerm_resource_group.test.location
  resource_group_name         = azurerm_resource_group.test.name
  sku_name                    = "Standard_D1_v2"
  instances                   = 2
  platform_fault_domain_count = 1
  upgrade_mode                = "Rolling"
  zones                       = ["1"]

  os_profile {
    linux_configuration {
      computer_name_prefix            = "testvm-%[1]d"
      admin_username                  = "myadmin"
      admin_password                  = "Passwword1234"
      disable_password_authentication = false
    }
  }

  network_interface {
    name    = "TestNetworkProfile-%[1]d"
    primary = true

    ip_configuration {
      name      = "TestIPConfiguration"
      primary   = true
      subnet_id = azurerm_subnet.test.id

      public_ip_address {
        name                    = "TestPublicIPConfiguration"
        domain_name_label       = "test-domain-label"
        idle_timeout_in_minutes = 4
      }
    }
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }
}
`, data.RandomInteger, data.Locations.Primary, r.natgateway_template(data))
}

func (OrchestratedVirtualMachineScaleSetResource) otherCapacityReservationGroupId(data acceptance.TestData) string {
	r := OrchestratedVirtualMachineScaleSetResource{}
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-OVMSS-%[1]d"
  location = "%[2]s"
}

resource "azurerm_capacity_reservation_group" "test" {
  name                = "acctest-ccrg-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_capacity_reservation" "test" {
  name                          = "acctest-ccr-%[1]d"
  capacity_reservation_group_id = azurerm_capacity_reservation_group.test.id
  sku {
    name     = "Standard_F2"
    capacity = 1
  }
}

%[3]s

resource "azurerm_orchestrated_virtual_machine_scale_set" "test" {
  name                = "acctestOVMSS-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku_name  = "Standard_F2"
  instances = 1

  platform_fault_domain_count   = 2
  capacity_reservation_group_id = azurerm_capacity_reservation_group.test.id

  os_profile {
    linux_configuration {
      computer_name_prefix = "testvm"
      admin_username       = "myadmin"

      admin_ssh_key {
        username   = "myadmin"
        public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDvXYZAjVUt2aojUV3XIA+PY6gXrgbvktXwf2NoIHGlQFhogpMEyOfqgogCtTBM7MNCS3ELul6SV+mlpH08Ki45ADIQuDXdommCvsMFW096JrsHOJpGfjCsJ1gbbv7brB3Ag+BSGb4qO3pRsEVTtZCeJDwfH5D7vmqP5xXcELKR4UAtKQKUhLvt6mhW90sFLTJeOTiYGbavIKqfCUFSeSMQkUPr8o3uzOfeWyCw7tc7szLuvfwJ5poGHuve73KKAlUnDTPUrhyj7iITZSDl+/i+bpDzPyCyJWDMsC0ON7q2fDr2mEz0L9ACrsI5Nx3lt5fe+IaHSrjivqnL8SqUWSN45o9Qp99sGWFiuTfos8f1jp+AXzC4ArVtKyRg/CnzKRiK0CGSxBJ5s9zAoa7yBBmjCszq89vFa0eMgpEIZFwa6kKJKt9AfRBXgO9YGPV4uaN7topy92/p2pE+vF8IafarbvnTDOQt62mS07tXYqYg1DhecrmBVWKlq9oafBweoeTjoq52SoGsuDc/YAOzIgWVIuvV8yKoh9KbXPWowjLtxDhRIS/d1nMMNdNI8X0TQivgi5+umMgAXhsVAKSNDUauLt4jimYkWAuE+R6KoCqVFdaB9bQDySBjAziruDSe3reToydjzzluvHMjWK8QiDynxs41pi4zZz6gAlca3QPkEQ== hello@world.com"
      }
    }
  }

  network_interface {
    name    = "TestNetworkProfile"
    primary = true

    ip_configuration {
      name      = "TestIPConfiguration"
      primary   = true
      subnet_id = azurerm_subnet.test.id

      public_ip_address {
        name                    = "TestPublicIPConfiguration"
        domain_name_label       = "test-domain-label"
        idle_timeout_in_minutes = 4
      }
    }
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }

  depends_on = [
    azurerm_capacity_reservation.test,
  ]
}
`, data.RandomInteger, data.Locations.Primary, r.natgateway_template(data))
}

func (OrchestratedVirtualMachineScaleSetResource) otherVMAgentDisabled(data acceptance.TestData) string {
	r := OrchestratedVirtualMachineScaleSetResource{}
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-OVMSS-%[1]d"
  location = "%[2]s"
}

%[3]s

resource "azurerm_orchestrated_virtual_machine_scale_set" "test" {
  name                = "acctestOVMSS-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  platform_fault_domain_count = 2

  sku_name  = "Standard_F2"
  instances = 1

  extension_operations_enabled = false

  os_profile {
    linux_configuration {
      computer_name_prefix = "testvm"
      admin_username       = "myadmin"
      provision_vm_agent   = false

      admin_ssh_key {
        username   = "myadmin"
        public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDvXYZAjVUt2aojUV3XIA+PY6gXrgbvktXwf2NoIHGlQFhogpMEyOfqgogCtTBM7MNCS3ELul6SV+mlpH08Ki45ADIQuDXdommCvsMFW096JrsHOJpGfjCsJ1gbbv7brB3Ag+BSGb4qO3pRsEVTtZCeJDwfH5D7vmqP5xXcELKR4UAtKQKUhLvt6mhW90sFLTJeOTiYGbavIKqfCUFSeSMQkUPr8o3uzOfeWyCw7tc7szLuvfwJ5poGHuve73KKAlUnDTPUrhyj7iITZSDl+/i+bpDzPyCyJWDMsC0ON7q2fDr2mEz0L9ACrsI5Nx3lt5fe+IaHSrjivqnL8SqUWSN45o9Qp99sGWFiuTfos8f1jp+AXzC4ArVtKyRg/CnzKRiK0CGSxBJ5s9zAoa7yBBmjCszq89vFa0eMgpEIZFwa6kKJKt9AfRBXgO9YGPV4uaN7topy92/p2pE+vF8IafarbvnTDOQt62mS07tXYqYg1DhecrmBVWKlq9oafBweoeTjoq52SoGsuDc/YAOzIgWVIuvV8yKoh9KbXPWowjLtxDhRIS/d1nMMNdNI8X0TQivgi5+umMgAXhsVAKSNDUauLt4jimYkWAuE+R6KoCqVFdaB9bQDySBjAziruDSe3reToydjzzluvHMjWK8QiDynxs41pi4zZz6gAlca3QPkEQ== hello@world.com"
      }
    }
  }

  network_interface {
    name    = "TestNetworkProfile"
    primary = true

    ip_configuration {
      name      = "TestIPConfiguration"
      primary   = true
      subnet_id = azurerm_subnet.test.id

      public_ip_address {
        name                    = "TestPublicIPConfiguration"
        domain_name_label       = "test-domain-label"
        idle_timeout_in_minutes = 4
      }
    }
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }
}
`, data.RandomInteger, data.Locations.Primary, r.natgateway_template(data))
}

func (OrchestratedVirtualMachineScaleSetResource) otherVMAgentDisabledWithExtensionDisabled(data acceptance.TestData) string {
	r := OrchestratedVirtualMachineScaleSetResource{}
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-OVMSS-%[1]d"
  location = "%[2]s"
}

%[3]s

resource "azurerm_orchestrated_virtual_machine_scale_set" "test" {
  name                = "acctestOVMSS-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  platform_fault_domain_count = 2

  extension_operations_enabled = false
  sku_name                     = "Standard_F2"
  instances                    = 1

  os_profile {
    linux_configuration {
      computer_name_prefix = "testvm"
      admin_username       = "myadmin"
      provision_vm_agent   = false

      admin_ssh_key {
        username   = "myadmin"
        public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDvXYZAjVUt2aojUV3XIA+PY6gXrgbvktXwf2NoIHGlQFhogpMEyOfqgogCtTBM7MNCS3ELul6SV+mlpH08Ki45ADIQuDXdommCvsMFW096JrsHOJpGfjCsJ1gbbv7brB3Ag+BSGb4qO3pRsEVTtZCeJDwfH5D7vmqP5xXcELKR4UAtKQKUhLvt6mhW90sFLTJeOTiYGbavIKqfCUFSeSMQkUPr8o3uzOfeWyCw7tc7szLuvfwJ5poGHuve73KKAlUnDTPUrhyj7iITZSDl+/i+bpDzPyCyJWDMsC0ON7q2fDr2mEz0L9ACrsI5Nx3lt5fe+IaHSrjivqnL8SqUWSN45o9Qp99sGWFiuTfos8f1jp+AXzC4ArVtKyRg/CnzKRiK0CGSxBJ5s9zAoa7yBBmjCszq89vFa0eMgpEIZFwa6kKJKt9AfRBXgO9YGPV4uaN7topy92/p2pE+vF8IafarbvnTDOQt62mS07tXYqYg1DhecrmBVWKlq9oafBweoeTjoq52SoGsuDc/YAOzIgWVIuvV8yKoh9KbXPWowjLtxDhRIS/d1nMMNdNI8X0TQivgi5+umMgAXhsVAKSNDUauLt4jimYkWAuE+R6KoCqVFdaB9bQDySBjAziruDSe3reToydjzzluvHMjWK8QiDynxs41pi4zZz6gAlca3QPkEQ== hello@world.com"
      }
    }
  }

  network_interface {
    name    = "TestNetworkProfile"
    primary = true

    ip_configuration {
      name      = "TestIPConfiguration"
      primary   = true
      subnet_id = azurerm_subnet.test.id

      public_ip_address {
        name                    = "TestPublicIPConfiguration"
        domain_name_label       = "test-domain-label"
        idle_timeout_in_minutes = 4
      }
    }
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }
}
`, data.RandomInteger, data.Locations.Primary, r.natgateway_template(data))
}

func (OrchestratedVirtualMachineScaleSetResource) otherSinglePlacementGroup(data acceptance.TestData, singlePlacementGroup string) string {
	r := OrchestratedVirtualMachineScaleSetResource{}
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-OVMSS-%[1]d"
  location = "%[2]s"
}

%[3]s

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  name = "acctest%[4]s"
}

resource "azurerm_orchestrated_virtual_machine_scale_set" "test" {
  name                = "acctestOVMSS-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  zones = []

  sku_name  = "Standard_D1_v2"
  instances = 1

  platform_fault_domain_count = 2

  os_profile {
    linux_configuration {
      computer_name_prefix = "testvm-%[1]d"
      admin_username       = "myadmin"
      admin_password       = "Passwword1234"

      disable_password_authentication = false
    }
  }

  network_interface {
    name    = "TestNetworkProfile-%[1]d"
    primary = true

    ip_configuration {
      name      = "TestIPConfiguration"
      primary   = true
      subnet_id = azurerm_subnet.test.id

      public_ip_address {
        name                    = "TestPublicIPConfiguration"
        domain_name_label       = "test-domain-label"
        idle_timeout_in_minutes = 4
      }
    }
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }

  single_placement_group = %[5]s

  tags = {
    "platformsettings.host_environment.service.platform_optedin_for_rootcerts" = "true"
  }
}
`, data.RandomInteger, data.Locations.Primary, r.natgateway_template(data), data.RandomString, singlePlacementGroup)
}

func (OrchestratedVirtualMachineScaleSetResource) otherSinglePlacementGroupNull(data acceptance.TestData, vmSku string) string {
	r := OrchestratedVirtualMachineScaleSetResource{}
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-OVMSS-%[1]d"
  location = "%[2]s"
}

%[3]s

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  name = "acctest%[4]s"
}

resource "azurerm_orchestrated_virtual_machine_scale_set" "test" {
  name                = "acctestOVMSS-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  zones = []

  sku_name  = "%[5]s"
  instances = 1

  platform_fault_domain_count = 2

  os_profile {
    linux_configuration {
      computer_name_prefix = "testvm-%[1]d"
      admin_username       = "myadmin"
      admin_password       = "Passwword1234"

      disable_password_authentication = false
    }
  }

  network_interface {
    name    = "TestNetworkProfile-%[1]d"
    primary = true

    ip_configuration {
      name      = "TestIPConfiguration"
      primary   = true
      subnet_id = azurerm_subnet.test.id

      public_ip_address {
        name                    = "TestPublicIPConfiguration"
        domain_name_label       = "test-domain-label"
        idle_timeout_in_minutes = 4
      }
    }
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }

  tags = {
    "platformsettings.host_environment.service.platform_optedin_for_rootcerts" = "true"
  }
}
`, data.RandomInteger, data.Locations.Primary, r.natgateway_template(data), data.RandomString, vmSku)
}
