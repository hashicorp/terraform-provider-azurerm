// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package maintenance_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/maintenance/2022-07-01-preview/configurationassignments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type MaintenanceAssignmentVirtualMachineResource struct{}

func TestAccMaintenanceAssignmentVirtualMachine_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_maintenance_assignment_virtual_machine", "test")
	r := MaintenanceAssignmentVirtualMachineResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		// location not returned by list rest api
		data.ImportStep("location"),
	})
}

func TestAccMaintenanceAssignmentVirtualMachine_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_maintenance_assignment_virtual_machine", "test")
	r := MaintenanceAssignmentVirtualMachineResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccMaintenanceAssignmentVirtualMachine_linkMultipleMaintenanceAssignmentsToOneVM(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_maintenance_assignment_virtual_machine", "test")
	r := MaintenanceAssignmentVirtualMachineResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.linkMultipleMaintenanceAssignmentsToOneVM(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		// location not returned by list rest api
		data.ImportStep("location"),
	})
}

func (MaintenanceAssignmentVirtualMachineResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := configurationassignments.ParseScopedConfigurationAssignmentID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Maintenance.ConfigurationAssignmentsClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %v", id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (r MaintenanceAssignmentVirtualMachineResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_maintenance_assignment_virtual_machine" "test" {
  location                     = azurerm_resource_group.test.location
  maintenance_configuration_id = azurerm_maintenance_configuration.test.id
  virtual_machine_id           = azurerm_linux_virtual_machine.test.id
}
`, r.template(data))
}

func (r MaintenanceAssignmentVirtualMachineResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_maintenance_assignment_virtual_machine" "import" {
  location                     = azurerm_maintenance_assignment_virtual_machine.test.location
  maintenance_configuration_id = azurerm_maintenance_assignment_virtual_machine.test.maintenance_configuration_id
  virtual_machine_id           = azurerm_maintenance_assignment_virtual_machine.test.virtual_machine_id
}
`, r.basic(data))
}

func (MaintenanceAssignmentVirtualMachineResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-maint-%[1]d"
  location = "%[2]s"
}

resource "azurerm_maintenance_configuration" "test" {
  name                = "acctest-MC%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  scope               = "SQLDB"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestnw-%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_network_interface" "test" {
  name                = "acctni-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_linux_virtual_machine" "test" {
  name                = "acctestVM-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_F2"
  admin_username      = "adminuser"
  admin_password      = "P@$$w0rd1234!"

  disable_password_authentication = false

  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r MaintenanceAssignmentVirtualMachineResource) linkMultipleMaintenanceAssignmentsToOneVM(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_maintenance_configuration" "test2" {
  name                = "acctest-MC2%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  scope               = "SQLDB"
}

resource "azurerm_maintenance_assignment_virtual_machine" "test" {
  location                     = azurerm_resource_group.test.location
  maintenance_configuration_id = azurerm_maintenance_configuration.test.id
  virtual_machine_id           = azurerm_linux_virtual_machine.test.id
}

resource "azurerm_maintenance_assignment_virtual_machine" "test2" {
  location                     = azurerm_resource_group.test.location
  maintenance_configuration_id = azurerm_maintenance_configuration.test2.id
  virtual_machine_id           = azurerm_linux_virtual_machine.test.id
}
`, r.template(data), data.RandomInteger)
}
