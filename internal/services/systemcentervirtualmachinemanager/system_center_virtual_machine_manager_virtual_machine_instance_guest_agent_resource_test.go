// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package systemcentervirtualmachinemanager_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/systemcentervirtualmachinemanager/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type SystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentResource struct{}

func TestAccSystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentSequential(t *testing.T) {
	// NOTE: this is a combined test rather than separate split out tests because the testing is against the same Hybrid Machine

	if os.Getenv("ARM_TEST_CUSTOM_LOCATION_ID") == "" || os.Getenv("ARM_TEST_FQDN") == "" || os.Getenv("ARM_TEST_USERNAME") == "" || os.Getenv("ARM_TEST_PASSWORD") == "" {
		t.Skip("Skipping as one of `ARM_TEST_CUSTOM_LOCATION_ID`, `ARM_TEST_FQDN`, `ARM_TEST_USERNAME`, `ARM_TEST_PASSWORD` was not specified")
	}

	acceptance.RunTestsInSequence(t, map[string]map[string]func(t *testing.T){
		"scvmmInstanceGuestAgent": {
			"basic":          testAccSystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgent_basic,
			"requiresImport": testAccSystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgent_requiresImport,
			"complete":       testAccSystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgent_complete,
		},
	})
}

func testAccSystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgent_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_system_center_virtual_machine_manager_virtual_machine_instance_guest_agent", "test")
	r := SystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentResource{}

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

func testAccSystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgent_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_system_center_virtual_machine_manager_virtual_machine_instance_guest_agent", "test")
	r := SystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentResource{}

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

func testAccSystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgent_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_system_center_virtual_machine_manager_virtual_machine_instance_guest_agent", "test")
	r := SystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentResource{}

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

func (r SystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.SystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.SystemCenterVirtualMachineManager.GuestAgents.Get(ctx, commonids.NewScopeID(id.Scope))
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r SystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_system_center_virtual_machine_manager_virtual_machine_instance_guest_agent" "test" {
  scoped_resource_id = azurerm_arc_machine.test.id
  username           = "Administrator"
  password           = "AdminPassword123!"

  depends_on = [azurerm_system_center_virtual_machine_manager_virtual_machine_instance.test]
}
`, r.template(data))
}

func (r SystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_system_center_virtual_machine_manager_virtual_machine_instance_guest_agent" "import" {
  scoped_resource_id = azurerm_system_center_virtual_machine_manager_virtual_machine_instance_guest_agent.test.scoped_resource_id
  username           = "Administrator"
  password           = "AdminPassword123!"
}
`, r.basic(data))
}

func (r SystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_system_center_virtual_machine_manager_virtual_machine_instance_guest_agent" "test" {
  scoped_resource_id  = azurerm_arc_machine.test.id
  provisioning_action = "install"
  username            = "Administrator"
  password            = "AdminPassword123!"

  depends_on = [azurerm_system_center_virtual_machine_manager_virtual_machine_instance.test]
}
`, r.template(data))
}

func (r SystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestrg-scvmmvmiga-%d"
  location = "%s"
}

resource "azurerm_arc_machine" "test" {
  name                = "acctest-arcmachine-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  kind                = "SCVMM"

  identity {
    type = "SystemAssigned"
  }
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
  inventory_type                                  = "Cloud"
  system_center_virtual_machine_manager_server_id = azurerm_system_center_virtual_machine_manager_server.test.id
}

resource "azurerm_system_center_virtual_machine_manager_cloud" "test" {
  name                                                           = "acctest-scvmmc-%d"
  location                                                       = azurerm_resource_group.test.location
  resource_group_name                                            = azurerm_resource_group.test.name
  custom_location_id                                             = azurerm_system_center_virtual_machine_manager_server.test.custom_location_id
  system_center_virtual_machine_manager_server_inventory_item_id = data.azurerm_system_center_virtual_machine_manager_inventory_items.test.inventory_items[0].id
}

data "azurerm_system_center_virtual_machine_manager_inventory_items" "test2" {
  inventory_type                                  = "VirtualMachineTemplate"
  system_center_virtual_machine_manager_server_id = azurerm_system_center_virtual_machine_manager_server.test.id
}

resource "azurerm_system_center_virtual_machine_manager_virtual_machine_template" "test" {
  name                                                           = "acctest-scvmmvmt-%d"
  location                                                       = azurerm_resource_group.test.location
  resource_group_name                                            = azurerm_resource_group.test.name
  custom_location_id                                             = azurerm_system_center_virtual_machine_manager_server.test.custom_location_id
  system_center_virtual_machine_manager_server_inventory_item_id = data.azurerm_system_center_virtual_machine_manager_inventory_items.test2.inventory_items[0].id
}

resource "azurerm_system_center_virtual_machine_manager_virtual_machine_instance" "test" {
  scoped_resource_id = azurerm_arc_machine.test.id
  custom_location_id = azurerm_system_center_virtual_machine_manager_server.test.custom_location_id

  infrastructure {
    checkpoint_type                                                 = "Standard"
    system_center_virtual_machine_manager_cloud_id                  = azurerm_system_center_virtual_machine_manager_cloud.test.id
    system_center_virtual_machine_manager_template_id               = azurerm_system_center_virtual_machine_manager_virtual_machine_template.test.id
    system_center_virtual_machine_manager_virtual_machine_server_id = azurerm_system_center_virtual_machine_manager_server.test.id
  }

  operating_system {
    admin_password = "AdminPassword123!"
  }

  lifecycle {
    // Service API always provisions a virtual disk with bus type IDE, hardware, network interface per Virtual Machine Template by default
    ignore_changes = [storage_disk, hardware, network_interface, operating_system.0.computer_name]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, os.Getenv("ARM_TEST_CUSTOM_LOCATION_ID"), os.Getenv("ARM_TEST_FQDN"), os.Getenv("ARM_TEST_USERNAME"), os.Getenv("ARM_TEST_PASSWORD"), data.RandomInteger, data.RandomInteger)
}
