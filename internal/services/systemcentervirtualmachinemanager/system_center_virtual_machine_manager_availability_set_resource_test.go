// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package systemcentervirtualmachinemanager_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/systemcentervirtualmachinemanager/2023-10-07/availabilitysets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type SystemCenterVirtualMachineManagerAvailabilitySetResource struct{}

func TestAccSystemCenterVirtualMachineManagerAvailabilitySetSequential(t *testing.T) {
	// NOTE: this is a combined test rather than separate split out tests because only one System Center Virtual Machine Manager Server can be onboarded at a time on a given Custom Location

	if os.Getenv("ARM_TEST_CUSTOM_LOCATION_ID") == "" || os.Getenv("ARM_TEST_FQDN") == "" || os.Getenv("ARM_TEST_USERNAME") == "" || os.Getenv("ARM_TEST_PASSWORD") == "" {
		t.Skip("Skipping as one of `ARM_TEST_CUSTOM_LOCATION_ID`, `ARM_TEST_FQDN`, `ARM_TEST_USERNAME`, `ARM_TEST_PASSWORD` was not specified")
	}

	acceptance.RunTestsInSequence(t, map[string]map[string]func(t *testing.T){
		"scvmmAvailabilitySet": {
			"basic":          testAccSystemCenterVirtualMachineManagerAvailabilitySet_basic,
			"requiresImport": testAccSystemCenterVirtualMachineManagerAvailabilitySet_requiresImport,
			"complete":       testAccSystemCenterVirtualMachineManagerAvailabilitySet_complete,
			"update":         testAccSystemCenterVirtualMachineManagerAvailabilitySet_update,
		},
	})
}

func testAccSystemCenterVirtualMachineManagerAvailabilitySet_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_system_center_virtual_machine_manager_availability_set", "test")
	r := SystemCenterVirtualMachineManagerAvailabilitySetResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccSystemCenterVirtualMachineManagerAvailabilitySet_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_system_center_virtual_machine_manager_availability_set", "test")
	r := SystemCenterVirtualMachineManagerAvailabilitySetResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func testAccSystemCenterVirtualMachineManagerAvailabilitySet_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_system_center_virtual_machine_manager_availability_set", "test")
	r := SystemCenterVirtualMachineManagerAvailabilitySetResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccSystemCenterVirtualMachineManagerAvailabilitySet_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_system_center_virtual_machine_manager_availability_set", "test")
	r := SystemCenterVirtualMachineManagerAvailabilitySetResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r SystemCenterVirtualMachineManagerAvailabilitySetResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := availabilitysets.ParseAvailabilitySetID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.SystemCenterVirtualMachineManager.AvailabilitySets.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (r SystemCenterVirtualMachineManagerAvailabilitySetResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_system_center_virtual_machine_manager_availability_set" "test" {
  name                                            = "acctest-scvmmas-%d"
  resource_group_name                             = azurerm_resource_group.test.name
  location                                        = azurerm_resource_group.test.location
  custom_location_id                              = azurerm_system_center_virtual_machine_manager_server.test.custom_location_id
  system_center_virtual_machine_manager_server_id = azurerm_system_center_virtual_machine_manager_server.test.id
}
`, r.template(data), data.RandomInteger)
}

func (r SystemCenterVirtualMachineManagerAvailabilitySetResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_system_center_virtual_machine_manager_availability_set" "import" {
  name                                            = azurerm_system_center_virtual_machine_manager_availability_set.test.name
  resource_group_name                             = azurerm_system_center_virtual_machine_manager_availability_set.test.resource_group_name
  location                                        = azurerm_system_center_virtual_machine_manager_availability_set.test.location
  custom_location_id                              = azurerm_system_center_virtual_machine_manager_availability_set.test.custom_location_id
  system_center_virtual_machine_manager_server_id = azurerm_system_center_virtual_machine_manager_availability_set.test.system_center_virtual_machine_manager_server_id
}
`, r.basic(data))
}

func (r SystemCenterVirtualMachineManagerAvailabilitySetResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_system_center_virtual_machine_manager_availability_set" "test" {
  name                                            = "acctest-scvmmas-%d"
  resource_group_name                             = azurerm_resource_group.test.name
  location                                        = azurerm_resource_group.test.location
  custom_location_id                              = azurerm_system_center_virtual_machine_manager_server.test.custom_location_id
  system_center_virtual_machine_manager_server_id = azurerm_system_center_virtual_machine_manager_server.test.id

  tags = {
    env = "Test"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r SystemCenterVirtualMachineManagerAvailabilitySetResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_system_center_virtual_machine_manager_availability_set" "test" {
  name                                            = "acctest-scvmmas-%d"
  resource_group_name                             = azurerm_resource_group.test.name
  location                                        = azurerm_resource_group.test.location
  custom_location_id                              = azurerm_system_center_virtual_machine_manager_server.test.custom_location_id
  system_center_virtual_machine_manager_server_id = azurerm_system_center_virtual_machine_manager_server.test.id

  tags = {
    env = "Test2"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r SystemCenterVirtualMachineManagerAvailabilitySetResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestrg-scvmmas-%d"
  location = "%s"
}

resource "azurerm_system_center_virtual_machine_manager_server" "test" {
  name                = "acctest-scvmmms-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  custom_location_id  = "%s"
  fqdn                = "%s"
  username            = "%s"
  password            = "%s"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, os.Getenv("ARM_TEST_CUSTOM_LOCATION_ID"), os.Getenv("ARM_TEST_FQDN"), os.Getenv("ARM_TEST_USERNAME"), os.Getenv("ARM_TEST_PASSWORD"))
}
