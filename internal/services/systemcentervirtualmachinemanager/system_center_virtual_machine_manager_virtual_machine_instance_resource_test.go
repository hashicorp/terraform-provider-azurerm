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

type SystemCenterVirtualMachineManagerVirtualMachineInstanceResource struct{}

func TestAccSystemCenterVirtualMachineManagerVirtualMachineInstanceSequential(t *testing.T) {
	// NOTE: this is a combined test rather than separate split out tests because only one System Center Virtual Machine Manager Server can be onboarded at a time on a given Custom Location

	if os.Getenv("ARM_TEST_CUSTOM_LOCATION_ID") == "" || os.Getenv("ARM_TEST_FQDN") == "" || os.Getenv("ARM_TEST_USERNAME") == "" || os.Getenv("ARM_TEST_PASSWORD") == "" || os.Getenv("ARM_TEST_VIRTUAL_NETWORK_INVENTORY_NAME") == "" {
		t.Skip("Skipping as one of `ARM_TEST_CUSTOM_LOCATION_ID`, `ARM_TEST_FQDN`, `ARM_TEST_USERNAME`, `ARM_TEST_PASSWORD`, `ARM_TEST_VIRTUAL_NETWORK_INVENTORY_NAME` was not specified")
	}

	acceptance.RunTestsInSequence(t, map[string]map[string]func(t *testing.T){
		"scvmmVirtualMachineInstance": {
			"basic":           testAccSystemCenterVirtualMachineManagerVirtualMachineInstance_basic,
			"requiresImport":  testAccSystemCenterVirtualMachineManagerVirtualMachineInstance_requiresImport,
			"complete":        testAccSystemCenterVirtualMachineManagerVirtualMachineInstance_complete,
			"update":          testAccSystemCenterVirtualMachineManagerVirtualMachineInstance_update,
			"inventoryItemId": testAccSystemCenterVirtualMachineManagerVirtualMachineInstance_inventoryItemId,
		},
	})
}

func testAccSystemCenterVirtualMachineManagerVirtualMachineInstance_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_system_center_virtual_machine_manager_virtual_machine_instance", "test")
	r := SystemCenterVirtualMachineManagerVirtualMachineInstanceResource{}

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

func testAccSystemCenterVirtualMachineManagerVirtualMachineInstance_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_system_center_virtual_machine_manager_virtual_machine_instance", "test")
	r := SystemCenterVirtualMachineManagerVirtualMachineInstanceResource{}

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

func testAccSystemCenterVirtualMachineManagerVirtualMachineInstance_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_system_center_virtual_machine_manager_virtual_machine_instance", "test")
	r := SystemCenterVirtualMachineManagerVirtualMachineInstanceResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("operating_system.0.admin_password"),
	})
}

func testAccSystemCenterVirtualMachineManagerVirtualMachineInstance_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_system_center_virtual_machine_manager_virtual_machine_instance", "test")
	r := SystemCenterVirtualMachineManagerVirtualMachineInstanceResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("operating_system.0.admin_password"),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("operating_system.0.admin_password"),
	})
}

func testAccSystemCenterVirtualMachineManagerVirtualMachineInstance_inventoryItemId(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_system_center_virtual_machine_manager_virtual_machine_instance", "test")
	r := SystemCenterVirtualMachineManagerVirtualMachineInstanceResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.inventoryItemId(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r SystemCenterVirtualMachineManagerVirtualMachineInstanceResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.SystemCenterVirtualMachineManagerVirtualMachineInstanceID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.SystemCenterVirtualMachineManager.VirtualMachineInstances.Get(ctx, commonids.NewScopeID(id.Scope))
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r SystemCenterVirtualMachineManagerVirtualMachineInstanceResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
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
    computer_name = "testComputer"
  }

  hardware {
    cpu_count    = 1
    memory_in_mb = 1024
  }

  network_interface {
    name              = "testNIC"
    ipv4_address_type = "Dynamic"
    ipv6_address_type = "Dynamic"
    mac_address_type  = "Dynamic"
  }

  lifecycle {
    // Service API always provisions a virtual disk with bus type IDE per Virtual Machine Template by default
    ignore_changes = [storage_disk]
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r SystemCenterVirtualMachineManagerVirtualMachineInstanceResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_system_center_virtual_machine_manager_virtual_machine_instance" "import" {
  scoped_resource_id = azurerm_system_center_virtual_machine_manager_virtual_machine_instance.test.scoped_resource_id
  custom_location_id = azurerm_system_center_virtual_machine_manager_virtual_machine_instance.test.custom_location_id

  infrastructure {
    checkpoint_type                                                 = "Standard"
    system_center_virtual_machine_manager_cloud_id                  = azurerm_system_center_virtual_machine_manager_cloud.test.id
    system_center_virtual_machine_manager_template_id               = azurerm_system_center_virtual_machine_manager_virtual_machine_template.test.id
    system_center_virtual_machine_manager_virtual_machine_server_id = azurerm_system_center_virtual_machine_manager_server.test.id
  }

  operating_system {
    computer_name = "testComputer"
  }

  hardware {
    cpu_count    = 1
    memory_in_mb = 1024
  }

  network_interface {
    name              = "testNIC"
    ipv4_address_type = "Dynamic"
    ipv6_address_type = "Dynamic"
    mac_address_type  = "Dynamic"
  }
}
`, r.basic(data))
}

func (r SystemCenterVirtualMachineManagerVirtualMachineInstanceResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
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

data "azurerm_system_center_virtual_machine_manager_inventory_items" "test3" {
  inventory_type                                  = "VirtualNetwork"
  system_center_virtual_machine_manager_server_id = azurerm_system_center_virtual_machine_manager_server.test.id
}

resource "azurerm_system_center_virtual_machine_manager_virtual_network" "test" {
  name                                                           = "acctest-scvmmvnet-%d"
  location                                                       = azurerm_resource_group.test.location
  resource_group_name                                            = azurerm_resource_group.test.name
  custom_location_id                                             = azurerm_system_center_virtual_machine_manager_server.test.custom_location_id
  system_center_virtual_machine_manager_server_inventory_item_id = [for item in data.azurerm_system_center_virtual_machine_manager_inventory_items.test3.inventory_items : item.id if item.name == "%s"][0]
}

resource "azurerm_system_center_virtual_machine_manager_availability_set" "test" {
  name                                            = "acctest-scvmmas-%d"
  resource_group_name                             = azurerm_resource_group.test.name
  location                                        = azurerm_resource_group.test.location
  custom_location_id                              = azurerm_system_center_virtual_machine_manager_server.test.custom_location_id
  system_center_virtual_machine_manager_server_id = azurerm_system_center_virtual_machine_manager_server.test.id
}

resource "azurerm_system_center_virtual_machine_manager_virtual_machine_instance" "test" {
  scoped_resource_id = azurerm_arc_machine.test.id
  custom_location_id = azurerm_system_center_virtual_machine_manager_server.test.custom_location_id

  infrastructure {
    checkpoint_type                                                 = "Production"
    system_center_virtual_machine_manager_cloud_id                  = azurerm_system_center_virtual_machine_manager_cloud.test.id
    system_center_virtual_machine_manager_template_id               = azurerm_system_center_virtual_machine_manager_virtual_machine_template.test.id
    system_center_virtual_machine_manager_virtual_machine_server_id = azurerm_system_center_virtual_machine_manager_server.test.id
  }

  operating_system {
    computer_name  = "testComputer"
    admin_password = "AdminPassword123!"
  }

  hardware {
    limit_cpu_for_migration_enabled = false
    cpu_count                       = 1
    memory_in_mb                    = 1024
    dynamic_memory_min_in_mb        = 32
    dynamic_memory_max_in_mb        = 1024
  }

  network_interface {
    name               = "testNIC%s"
    virtual_network_id = azurerm_system_center_virtual_machine_manager_virtual_network.test.id
    ipv4_address_type  = "Dynamic"
    ipv6_address_type  = "Dynamic"
    mac_address_type   = "Dynamic"
  }

  storage_disk {
    name         = "testSD%s"
    bus_type     = "SCSI"
    vhd_type     = "Dynamic"
    bus          = 1
    lun          = 1
    disk_size_gb = 30
  }

  system_center_virtual_machine_manager_availability_set_ids = [azurerm_system_center_virtual_machine_manager_availability_set.test.id]

  lifecycle {
    // Service API always provisions a virtual disk with bus type IDE per Virtual Machine Template by default, so it has to be ignored
    ignore_changes = [storage_disk]
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger, data.RandomInteger, os.Getenv("ARM_TEST_VIRTUAL_NETWORK_INVENTORY_NAME"), data.RandomInteger, data.RandomString, data.RandomString)
}

func (r SystemCenterVirtualMachineManagerVirtualMachineInstanceResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
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

data "azurerm_system_center_virtual_machine_manager_inventory_items" "test3" {
  inventory_type                                  = "VirtualNetwork"
  system_center_virtual_machine_manager_server_id = azurerm_system_center_virtual_machine_manager_server.test.id
}

resource "azurerm_system_center_virtual_machine_manager_virtual_network" "test" {
  name                                                           = "acctest-scvmmvnet-%d"
  location                                                       = azurerm_resource_group.test.location
  resource_group_name                                            = azurerm_resource_group.test.name
  custom_location_id                                             = azurerm_system_center_virtual_machine_manager_server.test.custom_location_id
  system_center_virtual_machine_manager_server_inventory_item_id = [for item in data.azurerm_system_center_virtual_machine_manager_inventory_items.test3.inventory_items : item.id if item.name == "%s"][0]
}

resource "azurerm_system_center_virtual_machine_manager_availability_set" "test" {
  name                                            = "acctest-scvmmas-%d"
  resource_group_name                             = azurerm_resource_group.test.name
  location                                        = azurerm_resource_group.test.location
  custom_location_id                              = azurerm_system_center_virtual_machine_manager_server.test.custom_location_id
  system_center_virtual_machine_manager_server_id = azurerm_system_center_virtual_machine_manager_server.test.id
}

resource "azurerm_system_center_virtual_machine_manager_availability_set" "test2" {
  name                                            = "acctest-scvmmas2-%d"
  resource_group_name                             = azurerm_resource_group.test.name
  location                                        = azurerm_resource_group.test.location
  custom_location_id                              = azurerm_system_center_virtual_machine_manager_server.test.custom_location_id
  system_center_virtual_machine_manager_server_id = azurerm_system_center_virtual_machine_manager_server.test.id
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
    computer_name  = "testComputer"
    admin_password = "AdminPassword123!"
  }

  hardware {
    limit_cpu_for_migration_enabled = false
    cpu_count                       = 1
    memory_in_mb                    = 1048
    dynamic_memory_min_in_mb        = 64
    dynamic_memory_max_in_mb        = 1048
  }

  network_interface {
    name               = "testNIC2%s"
    virtual_network_id = azurerm_system_center_virtual_machine_manager_virtual_network.test.id
    ipv4_address_type  = "Static"
    ipv6_address_type  = "Static"
    mac_address_type   = "Static"
  }

  storage_disk {
    name         = "testSD%s"
    bus_type     = "SCSI"
    vhd_type     = "Dynamic"
    bus          = 1
    lun          = 1
    disk_size_gb = 30
  }

  system_center_virtual_machine_manager_availability_set_ids = [azurerm_system_center_virtual_machine_manager_availability_set.test.id, azurerm_system_center_virtual_machine_manager_availability_set.test2.id]

  lifecycle {
    // Service API always provisions a virtual disk with bus type IDE per Virtual Machine Template by default, so it has to be ignored
    ignore_changes = [storage_disk]
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger, data.RandomInteger, os.Getenv("ARM_TEST_VIRTUAL_NETWORK_INVENTORY_NAME"), data.RandomInteger, data.RandomInteger, data.RandomString, data.RandomString)
}

func (r SystemCenterVirtualMachineManagerVirtualMachineInstanceResource) inventoryItemId(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

data "azurerm_system_center_virtual_machine_manager_inventory_items" "test4" {
  inventory_type                                  = "VirtualMachine"
  system_center_virtual_machine_manager_server_id = azurerm_system_center_virtual_machine_manager_server.test.id
}

resource "azurerm_system_center_virtual_machine_manager_virtual_machine_instance" "test" {
  scoped_resource_id = azurerm_arc_machine.test.id
  custom_location_id = azurerm_system_center_virtual_machine_manager_server.test.custom_location_id

  infrastructure {
    checkpoint_type                                                 = "Standard"
    system_center_virtual_machine_manager_virtual_machine_server_id = azurerm_system_center_virtual_machine_manager_server.test.id
    system_center_virtual_machine_manager_inventory_item_id         = data.azurerm_system_center_virtual_machine_manager_inventory_items.test4.inventory_items[0].id
  }

  lifecycle {
    // Service API provisions VM Instance based on the existing VM Instance that includes hardware, network_interface, operating_system and storage_disk
    ignore_changes = [hardware, network_interface, operating_system, storage_disk]
  }
}
`, r.template(data))
}

func (r SystemCenterVirtualMachineManagerVirtualMachineInstanceResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestrg-scvmmvmi-%d"
  location = "%s"
}

resource "azurerm_arc_machine" "test" {
  name                = "acctest-arcmachine-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  kind                = "SCVMM"
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, os.Getenv("ARM_TEST_CUSTOM_LOCATION_ID"), os.Getenv("ARM_TEST_FQDN"), os.Getenv("ARM_TEST_USERNAME"), os.Getenv("ARM_TEST_PASSWORD"))
}
