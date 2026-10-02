// Copyright IBM Corp. 2014, 2025
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

func sequentialResourceTest(t *testing.T, data acceptance.TestData, testResource OrchestratedVirtualMachineScaleSetResource, steps []acceptance.TestStep) {
	refreshStep := acceptance.TestStep{
		RefreshState: true,
	}

	newSteps := make([]acceptance.TestStep, 0)
	for _, step := range steps {
		if !step.ImportState {
			newSteps = append(newSteps, step)
		} else {
			newSteps = append(newSteps, refreshStep)
			newSteps = append(newSteps, step)
		}
	}

	data.ResourceSequentialTest(t, testResource, newSteps)
}

func TestAccOrchestratedVirtualMachineScaleSet_skuProfile_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_orchestrated_virtual_machine_scale_set", "test")
	r := OrchestratedVirtualMachineScaleSetResource{}

	sequentialResourceTest(t, data, r, []acceptance.TestStep{
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

	sequentialResourceTest(t, data, r, []acceptance.TestStep{
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

	sequentialResourceTest(t, data, r, []acceptance.TestStep{
		{
			Config: r.skuProfileWithRank(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("os_profile.0.windows_configuration.0.admin_password"),
	})
}

func TestAccOrchestratedVirtualMachineScaleSet_skuProfile_withFiveVMSizesOutOfOrderDuplicateRanks(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_orchestrated_virtual_machine_scale_set", "test")
	r := OrchestratedVirtualMachineScaleSetResource{}

	sequentialResourceTest(t, data, r, []acceptance.TestStep{
		{
			Config: r.skuProfileWithFiveVMSizesOutOfOrderDuplicateRanks(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sku_profile.0.virtual_machine_size.#").HasValue("5"),
			),
		},
		data.ImportStep("os_profile.0.windows_configuration.0.admin_password"),
	})
}

func TestAccOrchestratedVirtualMachineScaleSet_skuProfile_customizeDiffValidation(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_orchestrated_virtual_machine_scale_set", "test")
	r := OrchestratedVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.skuProfileSkuNameIsNotMix(data),
			ExpectError: regexp.MustCompile("`sku_profile` can only be configured when `sku_name` is set to `Mix`, got `Standard_D2s_v3`"),
		},
		{
			Config:      r.skuProfileNotExist(data),
			ExpectError: regexp.MustCompile("`sku_profile` must be configured when `sku_name` is set to `Mix`"),
		},
		{
			Config:      r.skuProfileRankWithoutPrioritized(data),
			ExpectError: regexp.MustCompile("`rank` can only be set when `allocation_strategy` is `Prioritized`, got `CapacityOptimized`"),
		},
	})
}

func TestAccOrchestratedVirtualMachineScaleSet_skuProfile_forceNewOnRemovalWithSkuNameChange(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_orchestrated_virtual_machine_scale_set", "test")
	r := OrchestratedVirtualMachineScaleSetResource{}

	sequentialResourceTest(t, data, r, []acceptance.TestStep{
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
  instances                   = 1
  platform_fault_domain_count = 1

  os_profile {
    windows_configuration {
      computer_name_prefix      = "testvm"
      admin_username            = "myadmin"
      admin_password            = "Passwword1234"
      automatic_updates_enabled = true
      provision_vm_agent        = true
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

func (r OrchestratedVirtualMachineScaleSetResource) skuProfileWithFiveVMSizesOutOfOrderDuplicateRanks(data acceptance.TestData) string {
	return r.skuProfileTemplate(data) + "\n" + r.skuProfileConfig(data, "Mix", skuProfilePrioritizedWithFiveVMSizesOutOfOrderDuplicateRanks())
}

func (r OrchestratedVirtualMachineScaleSetResource) skuProfileSkuNameIsNotMix(data acceptance.TestData) string {
	return r.skuProfileTemplate(data) + "\n" + r.skuProfileConfig(data, "Standard_D2s_v3", skuProfileCapacityOptimized())
}

func (r OrchestratedVirtualMachineScaleSetResource) skuProfileNotExist(data acceptance.TestData) string {
	return r.skuProfileTemplate(data) + "\n" + r.skuProfileConfig(data, "Mix", "")
}

func (r OrchestratedVirtualMachineScaleSetResource) skuProfileRankWithoutPrioritized(data acceptance.TestData) string {
	skuProfileBlock := `  sku_profile {
    allocation_strategy = "CapacityOptimized"

    virtual_machine_size {
      name = "Standard_D2s_v3"
      rank = 1
    }

    virtual_machine_size {
      name = "Standard_D4s_v3"
      rank = 2
    }
  }`
	return r.skuProfileTemplate(data) + "\n" + r.skuProfileConfig(data, "Mix", skuProfileBlock)
}

// Helper functions for common SKU profile configurations
func skuProfileCapacityOptimized() string {
	return `  sku_profile {
    allocation_strategy = "CapacityOptimized"

    virtual_machine_size {
	  name = "Standard_A2_v2"
    }

    virtual_machine_size {
	  name = "Standard_D2s_v3"
    }

    virtual_machine_size {
	  name = "Standard_F2s_v2"
    }
  }`
}

func skuProfileLowestPrice() string {
	return `  sku_profile {
    allocation_strategy = "LowestPrice"

    virtual_machine_size {
	  name = "Standard_D2s_v3"
    }

    virtual_machine_size {
	  name = "Standard_A2_v2"
    }
  }`
}

func skuProfilePrioritizedWithRank() string {
	return `  sku_profile {
    allocation_strategy = "Prioritized"

    virtual_machine_size {
      name = "Standard_A2_v2"
      rank = 1
    }

    virtual_machine_size {
      name = "Standard_D2s_v3"
      rank = 3
    }

    virtual_machine_size {
      name = "Standard_F2s_v2"
      rank = 3
    }
  }`
}

func skuProfilePrioritizedWithFiveVMSizesOutOfOrderDuplicateRanks() string {
	return `  sku_profile {
    allocation_strategy = "Prioritized"

    virtual_machine_size {
      name = "Standard_A2_v2"
      rank = 3
    }

    virtual_machine_size {
      name = "Standard_D2s_v3"
      rank = 1
    }

    virtual_machine_size {
      name = "Standard_F2s_v2"
      rank = 2
    }

    virtual_machine_size {
      name = "Standard_F4s_v2"
      rank = 3
    }

    virtual_machine_size {
      name = "Standard_E2s_v3"
      rank = 1
    }
  }`
}

func (r OrchestratedVirtualMachineScaleSetResource) skuProfileForceNewTransition(data acceptance.TestData) string {
	return r.skuProfileTemplate(data) + "\n" + r.skuProfileConfig(data, "Standard_D2s_v3", "")
}
