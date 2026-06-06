// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package compute_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
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

func TestOrchestratedVirtualMachineScaleSet_legacyEncryptionAtHost(t *testing.T) {
	r := OrchestratedVirtualMachineScaleSetResource{}
	resource.UnitTest(t, resource.TestCase{
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		Steps: []resource.TestStep{
			{
				Config:      r.legacyEncryptionAtHostConfig(),
				ExpectError: regexp.MustCompile("`encryption_at_host_enabled,sku_name`"),
			},
		},
	})
}

func TestOrchestratedVirtualMachineScaleSet_legacyZoneBalance(t *testing.T) {
	r := OrchestratedVirtualMachineScaleSetResource{}
	resource.UnitTest(t, resource.TestCase{
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		Steps: []resource.TestStep{
			{
				Config:      r.legacyZoneBalanceConfig(),
				ExpectError: regexp.MustCompile("`sku_name,zone_balance`"),
			},
		},
	})
}

func TestOrchestratedVirtualMachineScaleSet_legacyUserDataBase64(t *testing.T) {
	r := OrchestratedVirtualMachineScaleSetResource{}
	resource.UnitTest(t, resource.TestCase{
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		Steps: []resource.TestStep{
			{
				Config:      r.legacyUserDataBase64Config(),
				ExpectError: regexp.MustCompile("`sku_name,user_data_base64`"),
			},
		},
	})
}

func TestOrchestratedVirtualMachineScaleSet_legacyLicenseType(t *testing.T) {
	r := OrchestratedVirtualMachineScaleSetResource{}
	resource.UnitTest(t, resource.TestCase{
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		Steps: []resource.TestStep{
			{
				Config:      r.legacyLicenseTypeConfig(),
				ExpectError: regexp.MustCompile("`license_type,sku_name`"),
			},
		},
	})
}

func TestOrchestratedVirtualMachineScaleSet_legacyPriority(t *testing.T) {
	r := OrchestratedVirtualMachineScaleSetResource{}
	resource.UnitTest(t, resource.TestCase{
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		Steps: []resource.TestStep{
			{
				Config:      r.legacyPriorityConfig(),
				ExpectError: regexp.MustCompile("`priority,sku_name`"),
			},
		},
	})
}

func (r OrchestratedVirtualMachineScaleSetResource) legacyEncryptionAtHostConfig() string {
	return `
provider "azurerm" {
  features {}
  subscription_id = "00000000-0000-0000-0000-000000000000"
}

resource "azurerm_orchestrated_virtual_machine_scale_set" "test" {
  name                = "acctestOVMSS"
  location            = "eastus"
  resource_group_name = "acctestRG"

  platform_fault_domain_count = 1

  encryption_at_host_enabled = true
}
`
}

func (r OrchestratedVirtualMachineScaleSetResource) legacyZoneBalanceConfig() string {
	return `
provider "azurerm" {
  features {}
  subscription_id = "00000000-0000-0000-0000-000000000000"
}

resource "azurerm_orchestrated_virtual_machine_scale_set" "test" {
  name                = "acctestOVMSS"
  location            = "eastus"
  resource_group_name = "acctestRG"

  platform_fault_domain_count = 1

  zone_balance = true

  zones = ["1", "2", "3"]
}
`
}

func (r OrchestratedVirtualMachineScaleSetResource) legacyUserDataBase64Config() string {
	return `
provider "azurerm" {
  features {}
  subscription_id = "00000000-0000-0000-0000-000000000000"
}

resource "azurerm_orchestrated_virtual_machine_scale_set" "test" {
  name                = "acctestOVMSS"
  location            = "eastus"
  resource_group_name = "acctestRG"

  platform_fault_domain_count = 1

  user_data_base64 = "dGVzdA=="
}
`
}

func (r OrchestratedVirtualMachineScaleSetResource) legacyLicenseTypeConfig() string {
	return `
provider "azurerm" {
  features {}
  subscription_id = "00000000-0000-0000-0000-000000000000"
}

resource "azurerm_orchestrated_virtual_machine_scale_set" "test" {
  name                = "acctestOVMSS"
  location            = "eastus"
  resource_group_name = "acctestRG"

  platform_fault_domain_count = 1

  license_type = "Windows_Server"
}
`
}

func (r OrchestratedVirtualMachineScaleSetResource) legacyPriorityConfig() string {
	return `
provider "azurerm" {
  features {}
  subscription_id = "00000000-0000-0000-0000-000000000000"
}

resource "azurerm_orchestrated_virtual_machine_scale_set" "test" {
  name                = "acctestOVMSS"
  location            = "eastus"
  resource_group_name = "acctestRG"

  platform_fault_domain_count = 1

  priority = "Spot"
}
`
}
