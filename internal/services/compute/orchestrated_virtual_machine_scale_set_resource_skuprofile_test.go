// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

func TestAccOrchestratedVirtualMachineScaleSet_skuProfile_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_orchestrated_virtual_machine_scale_set", "test")
	r := OrchestratedVirtualMachineScaleSetResource{}

	// TODO - Remove this override when Preview is rolled out to westeurope - currently only supported in EastUS, WestUS, EastUS2, and WestUS2
	data.Locations.Primary = "eastus2"

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.skuProfileBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("os_profile.0.windows_configuration.0.admin_password"),
	})
}

func TestAccOrchestratedVirtualMachineScaleSet_skuProfile_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_orchestrated_virtual_machine_scale_set", "test")
	r := OrchestratedVirtualMachineScaleSetResource{}

	// TODO - Remove this override when Preview is rolled out to westeurope - currently only supported in EastUS, WestUS, EastUS2, and WestUS2
	data.Locations.Primary = "eastus2"

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.skuProfileBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("os_profile.0.windows_configuration.0.admin_password"),
		{
			Config: r.skuProfileUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("os_profile.0.windows_configuration.0.admin_password"),
		{
			Config: r.skuProfileBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("os_profile.0.windows_configuration.0.admin_password"),
	})
}

func TestAccOrchestratedVirtualMachineScaleSet_skuProfile_withRank(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_orchestrated_virtual_machine_scale_set", "test")
	r := OrchestratedVirtualMachineScaleSetResource{}

	// TODO - Remove this override when Preview is rolled out to westeurope - currently only supported in EastUS, WestUS, EastUS2, and WestUS2
	data.Locations.Primary = "eastus2"

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.skuProfileWithRank(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("os_profile.0.windows_configuration.0.admin_password"),
	})
}

func TestAccOrchestratedVirtualMachineScaleSet_skuProfile_duplicateVMSizes(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_orchestrated_virtual_machine_scale_set", "test")
	r := OrchestratedVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.skuProfileDuplicateVMSizes(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				// TypeSet automatically deduplicates, so we should only have 1 VM size
				check.That(data.ResourceName).Key("sku_profile.0.vm_sizes.#").HasValue("1"),
			),
		},
		data.ImportStep("os_profile.0.windows_configuration.0.admin_password"),
	})
}

func TestAccOrchestratedVirtualMachineScaleSet_skuProfile_customizeDiffValidation(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_orchestrated_virtual_machine_scale_set", "test")
	r := OrchestratedVirtualMachineScaleSetResource{}

	// TODO - Remove this override when Preview is rolled out to westeurope - currently only supported in EastUS, WestUS, EastUS2, and WestUS2
	data.Locations.Primary = "eastus2"

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.skuProfileWithoutSkuName(data),
			ExpectError: regexp.MustCompile(`"sku_name" is not formatted properly, got ""`),
		},
		{
			Config:      r.skuProfileSkuNameIsNotMix(data),
			ExpectError: regexp.MustCompile("`sku_profile` can only be configured when `sku_name` is set to `Mix`, got `Standard_B1s`"),
		},
		{
			Config:      r.skuProfileNotExist(data),
			ExpectError: regexp.MustCompile("`sku_profile` must be configured when `sku_name` is set to `Mix`"),
		},
		{
			Config:      r.skuProfileRankWithoutPrioritized(data),
			ExpectError: regexp.MustCompile("`rank` can only be set when `allocation_strategy` is `Prioritized`, got `CapacityOptimized`"),
		},
		{
			Config:      r.skuProfilePrioritizedWithoutRank(data),
			ExpectError: regexp.MustCompile("when `allocation_strategy` is `Prioritized`, all `vm_sizes` must have the `rank` field set, `Standard_B1ls` is missing `rank`"),
		},
		{
			Config:      r.skuProfilePrioritizedWithNonConsecutiveRanks(data),
			ExpectError: regexp.MustCompile("the `rank` values must be consecutive starting from 0. Expected rank `1` but got `2`"),
		},
		{
			Config:      r.skuProfileWithInvalidPlatformFaultDomainCount(data),
			ExpectError: regexp.MustCompile("`sku_profile` can only be configured when `platform_fault_domain_count` is set to `1`, got `5`"),
		},
	})
}

func TestAccOrchestratedVirtualMachineScaleSet_skuProfile_forceNewOnRemovalWithSkuNameChange(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_orchestrated_virtual_machine_scale_set", "test")
	r := OrchestratedVirtualMachineScaleSetResource{}

	// TODO - Remove this override when Preview is rolled out to westeurope - currently only supported in EastUS, WestUS, EastUS2, and WestUS2
	data.Locations.Primary = "eastus2"

	data.ResourceTestIgnoreRecreate(t, r, []acceptance.TestStep{
		{
			Config: r.skuProfileBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("os_profile.0.windows_configuration.0.admin_password"),
		{
			Config: r.skuProfileForceNewTransition(data),
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction(data.ResourceName, plancheck.ResourceActionReplace),
				},
			},
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("os_profile.0.windows_configuration.0.admin_password"),
	})
}

func (r OrchestratedVirtualMachineScaleSetResource) skuProfileTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
  }
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-OVMSS-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestVN-%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]
}`, data.RandomInteger, data.Locations.Primary)
}

func (r OrchestratedVirtualMachineScaleSetResource) skuProfileConfig(data acceptance.TestData, skuName string, skuProfileBlock string) string {
	return fmt.Sprintf(`
resource "azurerm_orchestrated_virtual_machine_scale_set" "test" {
  name                = "acctestOVMSS-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku_name = "%[3]s"
%[4]s
  instances                   = 2
  platform_fault_domain_count = 1

  os_profile {
    windows_configuration {
      computer_name_prefix     = "testvm"
      admin_username           = "myadmin"
      admin_password           = "Passwword1234"
      enable_automatic_updates = true
      provision_vm_agent       = true
    }
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  network_interface {
    name    = "TestNetworkProfile-%[1]d"
    primary = true

    ip_configuration {
      name      = "TestIPConfiguration"
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2019-Datacenter"
    version   = "latest"
  }

  tags = {
    environment = "AccTest"
  }
}`, data.RandomInteger, data.Locations.Primary, skuName, skuProfileBlock)
}

func (r OrchestratedVirtualMachineScaleSetResource) skuProfileBasic(data acceptance.TestData) string {
	return r.skuProfileTemplate(data) + "\n" + r.skuProfileConfig(data, "Mix", skuProfileCapacityOptimized())
}

func (r OrchestratedVirtualMachineScaleSetResource) skuProfileUpdate(data acceptance.TestData) string {
	return r.skuProfileTemplate(data) + "\n" + r.skuProfileConfig(data, "Mix", skuProfileLowestPrice())
}

func (r OrchestratedVirtualMachineScaleSetResource) skuProfileWithRank(data acceptance.TestData) string {
	return r.skuProfileTemplate(data) + "\n" + r.skuProfileConfig(data, "Mix", skuProfilePrioritizedWithRank())
}

func (r OrchestratedVirtualMachineScaleSetResource) skuProfileWithoutSkuName(data acceptance.TestData) string {
	return r.skuProfileTemplate(data) + "\n" + r.skuProfileConfig(data, "", skuProfileCapacityOptimized())
}

func (r OrchestratedVirtualMachineScaleSetResource) skuProfileSkuNameIsNotMix(data acceptance.TestData) string {
	return r.skuProfileTemplate(data) + "\n" + r.skuProfileConfig(data, "Standard_B1s", skuProfileCapacityOptimized())
}

func (r OrchestratedVirtualMachineScaleSetResource) skuProfileNotExist(data acceptance.TestData) string {
	return r.skuProfileTemplate(data) + "\n" + r.skuProfileConfig(data, "Mix", "")
}

func (r OrchestratedVirtualMachineScaleSetResource) skuProfileRankWithoutPrioritized(data acceptance.TestData) string {
	skuProfileBlock := `  sku_profile {
    allocation_strategy = "CapacityOptimized"

    vm_sizes {
      name = "Standard_B1ls"
      rank = 0
    }

    vm_sizes {
      name = "Standard_B1s"
      rank = 1
    }
  }`
	return r.skuProfileTemplate(data) + "\n" + r.skuProfileConfig(data, "Mix", skuProfileBlock)
}

func (r OrchestratedVirtualMachineScaleSetResource) skuProfileDuplicateVMSizes(data acceptance.TestData) string {
	skuProfileBlock := `  sku_profile {
    allocation_strategy = "CapacityOptimized"

    vm_sizes {
      name = "Standard_B1ls"
    }

    vm_sizes {
      name = "Standard_B1ls"
    }
  }`
	return r.skuProfileTemplate(data) + "\n" + r.skuProfileConfig(data, "Mix", skuProfileBlock)
}

func (r OrchestratedVirtualMachineScaleSetResource) skuProfilePrioritizedWithoutRank(data acceptance.TestData) string {
	skuProfileBlock := `  sku_profile {
    allocation_strategy = "Prioritized"

    vm_sizes {
      name = "Standard_B1ls"
    }

    vm_sizes {
      name = "Standard_B1s"
    }
  }`
	return r.skuProfileTemplate(data) + "\n" + r.skuProfileConfig(data, "Mix", skuProfileBlock)
}

func (r OrchestratedVirtualMachineScaleSetResource) skuProfilePrioritizedWithNonConsecutiveRanks(data acceptance.TestData) string {
	skuProfileBlock := `  sku_profile {
    allocation_strategy = "Prioritized"

    vm_sizes {
      name = "Standard_B1ls"
      rank = 0
    }

    vm_sizes {
      name = "Standard_B1s"
      rank = 2
    }
  }`
	return r.skuProfileTemplate(data) + "\n" + r.skuProfileConfig(data, "Mix", skuProfileBlock)
}

func (r OrchestratedVirtualMachineScaleSetResource) skuProfileWithInvalidPlatformFaultDomainCount(data acceptance.TestData) string {
	skuProfileBlock := `  sku_profile {
    allocation_strategy = "CapacityOptimized"

    vm_sizes {
      name = "Standard_B1ls"
    }

    vm_sizes {
      name = "Standard_B1s"
    }
  }`
	return r.skuProfileTemplate(data) + "\n" + r.skuProfileConfigWithPlatformFaultDomainCount(data, "Mix", skuProfileBlock, 5)
}

func (r OrchestratedVirtualMachineScaleSetResource) skuProfileConfigWithPlatformFaultDomainCount(data acceptance.TestData, skuName string, skuProfileBlock string, platformFaultDomainCount int) string {
	return fmt.Sprintf(`
resource "azurerm_orchestrated_virtual_machine_scale_set" "test" {
  name                = "acctestOVMSS-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku_name = "%[3]s"
%[4]s
  instances                   = 2
  platform_fault_domain_count = %[5]d

  os_profile {
    windows_configuration {
      computer_name_prefix     = "testvm"
      admin_username           = "myadmin"
      admin_password           = "Passwword1234"
      enable_automatic_updates = true
      provision_vm_agent       = true
    }
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  network_interface {
    name    = "TestNetworkProfile-%[1]d"
    primary = true

    ip_configuration {
      name      = "TestIPConfiguration"
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2019-Datacenter"
    version   = "latest"
  }

  tags = {
    environment = "AccTest"
  }
}`, data.RandomInteger, data.Locations.Primary, skuName, skuProfileBlock, platformFaultDomainCount)
}

// Helper functions for common SKU profile configurations
func skuProfileCapacityOptimized() string {
	return `  sku_profile {
    allocation_strategy = "CapacityOptimized"

    vm_sizes {
      name = "Standard_B1ls"
    }

    vm_sizes {
      name = "Standard_B1s"
    }

    vm_sizes {
      name = "Standard_B2s"
    }
  }`
}

func skuProfileLowestPrice() string {
	return `  sku_profile {
    allocation_strategy = "LowestPrice"

    vm_sizes {
      name = "Standard_B1s"
    }

    vm_sizes {
      name = "Standard_B1ls"
    }
  }`
}

func skuProfilePrioritizedWithRank() string {
	return `  sku_profile {
    allocation_strategy = "Prioritized"

    vm_sizes {
      name = "Standard_B1ls"
      rank = 0
    }

    vm_sizes {
      name = "Standard_B1s"
      rank = 1
    }

    vm_sizes {
      name = "Standard_B2s"
      rank = 2
    }
  }`
}

func (r OrchestratedVirtualMachineScaleSetResource) skuProfileForceNewTransition(data acceptance.TestData) string {
	return r.skuProfileTemplate(data) + "\n" + r.skuProfileConfig(data, "Standard_B1s", "")
}
