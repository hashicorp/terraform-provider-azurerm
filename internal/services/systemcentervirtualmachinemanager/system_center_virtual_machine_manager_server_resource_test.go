// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package systemcentervirtualmachinemanager_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/systemcentervirtualmachinemanager/2023-10-07/vmmservers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type SystemCenterVirtualMachineManagerServerResource struct{}

func TestAccSystemCenterVirtualMachineManagerServerSequential(t *testing.T) {
	// NOTE: this is a combined test rather than separate split out tests because only one System Center Virtual Machine Manager Server can be onboarded at a time on a given Custom Location

	if os.Getenv("ARM_TEST_CUSTOM_LOCATION_ID") == "" || os.Getenv("ARM_TEST_FQDN") == "" || os.Getenv("ARM_TEST_PORT") == "" || os.Getenv("ARM_TEST_USERNAME") == "" || os.Getenv("ARM_TEST_PASSWORD") == "" {
		t.Skip("Skipping as one of `ARM_TEST_CUSTOM_LOCATION_ID`, `ARM_TEST_FQDN`, `ARM_TEST_PORT`, `ARM_TEST_USERNAME`, `ARM_TEST_PASSWORD` was not specified")
	}

	acceptance.RunTestsInSequence(t, map[string]map[string]func(t *testing.T){
		"scvmmServer": {
			"basic":          testAccSystemCenterVirtualMachineManagerServer_basic,
			"requiresImport": testAccSystemCenterVirtualMachineManagerServer_requiresImport,
			"complete":       testAccSystemCenterVirtualMachineManagerServer_complete,
			"update":         testAccSystemCenterVirtualMachineManagerServer_update,
		},
	})
}

func testAccSystemCenterVirtualMachineManagerServer_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_system_center_virtual_machine_manager_server", "test")
	r := SystemCenterVirtualMachineManagerServerResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("password"),
	})
}

func testAccSystemCenterVirtualMachineManagerServer_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_system_center_virtual_machine_manager_server", "test")
	r := SystemCenterVirtualMachineManagerServerResource{}

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

func testAccSystemCenterVirtualMachineManagerServer_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_system_center_virtual_machine_manager_server", "test")
	r := SystemCenterVirtualMachineManagerServerResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("password"),
	})
}

func testAccSystemCenterVirtualMachineManagerServer_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_system_center_virtual_machine_manager_server", "test")
	r := SystemCenterVirtualMachineManagerServerResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("password"),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("password"),
	})
}

func (r SystemCenterVirtualMachineManagerServerResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := vmmservers.ParseVMmServerID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.SystemCenterVirtualMachineManager.VMmServers.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (r SystemCenterVirtualMachineManagerServerResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
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
`, r.template(data), data.RandomInteger, os.Getenv("ARM_TEST_CUSTOM_LOCATION_ID"), os.Getenv("ARM_TEST_FQDN"), os.Getenv("ARM_TEST_USERNAME"), os.Getenv("ARM_TEST_PASSWORD"))
}

func (r SystemCenterVirtualMachineManagerServerResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_system_center_virtual_machine_manager_server" "import" {
  name                = azurerm_system_center_virtual_machine_manager_server.test.name
  resource_group_name = azurerm_system_center_virtual_machine_manager_server.test.resource_group_name
  location            = azurerm_system_center_virtual_machine_manager_server.test.location
  custom_location_id  = azurerm_system_center_virtual_machine_manager_server.test.custom_location_id
  fqdn                = azurerm_system_center_virtual_machine_manager_server.test.fqdn
  username            = "%s"
  password            = "%s"
}
`, r.basic(data), os.Getenv("ARM_TEST_USERNAME"), os.Getenv("ARM_TEST_PASSWORD"))
}

func (r SystemCenterVirtualMachineManagerServerResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_system_center_virtual_machine_manager_server" "test" {
  name                = "acctest-scvmmms-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  custom_location_id  = "%s"
  fqdn                = "%s"
  port                = tonumber("%s")
  username            = "%s"
  password            = "%s"

  tags = {
    Env = "Test"
  }
}
`, r.template(data), data.RandomInteger, os.Getenv("ARM_TEST_CUSTOM_LOCATION_ID"), os.Getenv("ARM_TEST_FQDN"), os.Getenv("ARM_TEST_PORT"), os.Getenv("ARM_TEST_USERNAME"), os.Getenv("ARM_TEST_PASSWORD"))
}

func (r SystemCenterVirtualMachineManagerServerResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_system_center_virtual_machine_manager_server" "test" {
  name                = "acctest-scvmmms-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  custom_location_id  = "%s"
  fqdn                = "%s"
  port                = tonumber("%s")
  username            = "%s"
  password            = "%s"

  tags = {
    Env = "Test2"
  }
}
`, r.template(data), data.RandomInteger, os.Getenv("ARM_TEST_CUSTOM_LOCATION_ID"), os.Getenv("ARM_TEST_FQDN"), os.Getenv("ARM_TEST_PORT"), os.Getenv("ARM_TEST_USERNAME"), os.Getenv("ARM_TEST_PASSWORD"))
}

func (r SystemCenterVirtualMachineManagerServerResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestrg-scvmmms-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}
