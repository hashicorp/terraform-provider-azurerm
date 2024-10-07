// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package systemcentervirtualmachinemanager_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/systemcentervirtualmachinemanager/2023-10-07/virtualmachinetemplates"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type SystemCenterVirtualMachineManagerVirtualMachineTemplateResource struct{}

func TestAccSystemCenterVirtualMachineManagerVirtualMachineTemplateSequential(t *testing.T) {
	// NOTE: this is a combined test rather than separate split out tests because only one System Center Virtual Machine Manager Server can be onboarded at a time on a given Custom Location

	if os.Getenv("ARM_TEST_CUSTOM_LOCATION_ID") == "" || os.Getenv("ARM_TEST_FQDN") == "" || os.Getenv("ARM_TEST_USERNAME") == "" || os.Getenv("ARM_TEST_PASSWORD") == "" {
		t.Skip("Skipping as one of `ARM_TEST_CUSTOM_LOCATION_ID`, `ARM_TEST_FQDN`, `ARM_TEST_USERNAME`, `ARM_TEST_PASSWORD` was not specified")
	}

	acceptance.RunTestsInSequence(t, map[string]map[string]func(t *testing.T){
		"scvmmVirtualMachineTemplate": {
			"basic":          testAccSystemCenterVirtualMachineManagerVirtualMachineTemplate_basic,
			"requiresImport": testAccSystemCenterVirtualMachineManagerVirtualMachineTemplate_requiresImport,
			"complete":       testAccSystemCenterVirtualMachineManagerVirtualMachineTemplate_complete,
			"update":         testAccSystemCenterVirtualMachineManagerVirtualMachineTemplate_update,
		},
	})
}

func testAccSystemCenterVirtualMachineManagerVirtualMachineTemplate_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_system_center_virtual_machine_manager_virtual_machine_template", "test")
	r := SystemCenterVirtualMachineManagerVirtualMachineTemplateResource{}

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

func testAccSystemCenterVirtualMachineManagerVirtualMachineTemplate_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_system_center_virtual_machine_manager_virtual_machine_template", "test")
	r := SystemCenterVirtualMachineManagerVirtualMachineTemplateResource{}

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

func testAccSystemCenterVirtualMachineManagerVirtualMachineTemplate_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_system_center_virtual_machine_manager_virtual_machine_template", "test")
	r := SystemCenterVirtualMachineManagerVirtualMachineTemplateResource{}

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

func testAccSystemCenterVirtualMachineManagerVirtualMachineTemplate_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_system_center_virtual_machine_manager_virtual_machine_template", "test")
	r := SystemCenterVirtualMachineManagerVirtualMachineTemplateResource{}

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

func (r SystemCenterVirtualMachineManagerVirtualMachineTemplateResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := virtualmachinetemplates.ParseVirtualMachineTemplateID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.SystemCenterVirtualMachineManager.VirtualMachineTemplates.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (r SystemCenterVirtualMachineManagerVirtualMachineTemplateResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_system_center_virtual_machine_manager_virtual_machine_template" "test" {
  name                                                           = "acctest-scvmmc-%d"
  location                                                       = azurerm_resource_group.test.location
  resource_group_name                                            = azurerm_resource_group.test.name
  custom_location_id                                             = azurerm_system_center_virtual_machine_manager_server.test.custom_location_id
  system_center_virtual_machine_manager_server_inventory_item_id = data.azurerm_system_center_virtual_machine_manager_inventory_items.test.inventory_items[0].id
}
`, r.template(data), data.RandomInteger)
}

func (r SystemCenterVirtualMachineManagerVirtualMachineTemplateResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_system_center_virtual_machine_manager_virtual_machine_template" "import" {
  name                                                           = azurerm_system_center_virtual_machine_manager_virtual_machine_template.test.name
  location                                                       = azurerm_system_center_virtual_machine_manager_virtual_machine_template.test.location
  resource_group_name                                            = azurerm_system_center_virtual_machine_manager_virtual_machine_template.test.resource_group_name
  custom_location_id                                             = azurerm_system_center_virtual_machine_manager_virtual_machine_template.test.custom_location_id
  system_center_virtual_machine_manager_server_inventory_item_id = azurerm_system_center_virtual_machine_manager_virtual_machine_template.test.system_center_virtual_machine_manager_server_inventory_item_id
}
`, r.basic(data))
}

func (r SystemCenterVirtualMachineManagerVirtualMachineTemplateResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_system_center_virtual_machine_manager_virtual_machine_template" "test" {
  name                                                           = "acctest-scvmmc-%d"
  location                                                       = azurerm_resource_group.test.location
  resource_group_name                                            = azurerm_resource_group.test.name
  custom_location_id                                             = azurerm_system_center_virtual_machine_manager_server.test.custom_location_id
  system_center_virtual_machine_manager_server_inventory_item_id = data.azurerm_system_center_virtual_machine_manager_inventory_items.test.inventory_items[0].id

  tags = {
    env = "Test"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r SystemCenterVirtualMachineManagerVirtualMachineTemplateResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_system_center_virtual_machine_manager_virtual_machine_template" "test" {
  name                                                           = "acctest-scvmmc-%d"
  location                                                       = azurerm_resource_group.test.location
  resource_group_name                                            = azurerm_resource_group.test.name
  custom_location_id                                             = azurerm_system_center_virtual_machine_manager_server.test.custom_location_id
  system_center_virtual_machine_manager_server_inventory_item_id = data.azurerm_system_center_virtual_machine_manager_inventory_items.test.inventory_items[0].id

  tags = {
    env = "Test2"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r SystemCenterVirtualMachineManagerVirtualMachineTemplateResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestrg-scvmmvmtemplate-%d"
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

data "azurerm_system_center_virtual_machine_manager_inventory_items" "test" {
  inventory_type                                  = "VirtualMachineTemplate"
  system_center_virtual_machine_manager_server_id = azurerm_system_center_virtual_machine_manager_server.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, os.Getenv("ARM_TEST_CUSTOM_LOCATION_ID"), os.Getenv("ARM_TEST_FQDN"), os.Getenv("ARM_TEST_USERNAME"), os.Getenv("ARM_TEST_PASSWORD"))
}
