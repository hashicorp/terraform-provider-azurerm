// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

func TestAccOrchestratedVirtualMachineScaleSet_legacyBasicZonal(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_orchestrated_virtual_machine_scale_set", "test")
	r := OrchestratedVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.legacyBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("extensions_time_budget", "max_bid_price", "priority"),
	})
}

func TestAccOrchestratedVirtualMachineScaleSet_legacyUpdateZonal(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_orchestrated_virtual_machine_scale_set", "test")
	r := OrchestratedVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.legacyBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("extensions_time_budget", "max_bid_price", "priority"),
		{
			Config: r.legacyUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("extensions_time_budget", "max_bid_price", "priority"),
		{
			Config: r.legacyBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("extensions_time_budget", "max_bid_price", "priority"),
	})
}

func TestAccOrchestratedVirtualMachineScaleSet_legacyBasicNonZonal(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_orchestrated_virtual_machine_scale_set", "test")
	r := OrchestratedVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.legacyBasicNonZonal(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("extensions_time_budget", "max_bid_price", "priority"),
	})
}

func TestAccOrchestratedVirtualMachineScaleSet_legacyRequiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_orchestrated_virtual_machine_scale_set", "test")
	r := OrchestratedVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.legacyBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.legacyRequiresImport),
	})
}

// legacy test cases for backward compatibility validation
func (r OrchestratedVirtualMachineScaleSetResource) legacyBasic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_orchestrated_virtual_machine_scale_set" "test" {
  name                = "acctestOVMSS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  platform_fault_domain_count = 1

  zones = ["1"]

  tags = {
    ENV = "Test"
  }
}
`, r.legacyTemplate(data), data.RandomInteger)
}

func (r OrchestratedVirtualMachineScaleSetResource) legacyUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_orchestrated_virtual_machine_scale_set" "test" {
  name                = "acctestOVMSS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  platform_fault_domain_count = 1

  zones = ["1"]

  tags = {
    ENV = "Test",
    FOO = "Bar"
  }
}
`, r.legacyTemplate(data), data.RandomInteger)
}

func (r OrchestratedVirtualMachineScaleSetResource) legacyRequiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_orchestrated_virtual_machine_scale_set" "import" {
  name                = azurerm_orchestrated_virtual_machine_scale_set.test.name
  location            = azurerm_orchestrated_virtual_machine_scale_set.test.location
  resource_group_name = azurerm_orchestrated_virtual_machine_scale_set.test.resource_group_name

  platform_fault_domain_count = azurerm_orchestrated_virtual_machine_scale_set.test.platform_fault_domain_count
}
`, r.legacyBasic(data))
}

func (r OrchestratedVirtualMachineScaleSetResource) legacyBasicNonZonal(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_orchestrated_virtual_machine_scale_set" "test" {
  name                = "acctestOVMSS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  platform_fault_domain_count = 2

  tags = {
    ENV = "Test"
  }
}
`, r.legacyTemplate(data), data.RandomInteger)
}

func (OrchestratedVirtualMachineScaleSetResource) legacyTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-OVMSS-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}
