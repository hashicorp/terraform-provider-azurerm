// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package maintenance_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/maintenance/2023-04-01/configurationassignments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type MaintenanceAssignmentArcMachineResource struct{}

func TestAccMaintenanceAssignmentArcMachine_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_maintenance_assignment_arc_machine", "test")
	r := MaintenanceAssignmentArcMachineResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMaintenanceAssignmentArcMachine_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_maintenance_assignment_arc_machine", "test")
	r := MaintenanceAssignmentArcMachineResource{}

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

func (MaintenanceAssignmentArcMachineResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := configurationassignments.ParseScopedConfigurationAssignmentID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Maintenance.ConfigurationAssignmentsClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r MaintenanceAssignmentArcMachineResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_maintenance_assignment_arc_machine" "test" {
  arc_machine_id               = azurerm_arc_machine.test.id
  maintenance_configuration_id = azurerm_maintenance_configuration.test.id
}
`, r.template(data))
}

func (r MaintenanceAssignmentArcMachineResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_maintenance_assignment_arc_machine" "import" {
  arc_machine_id               = azurerm_maintenance_assignment_arc_machine.test.arc_machine_id
  maintenance_configuration_id = azurerm_maintenance_assignment_arc_machine.test.maintenance_configuration_id
}
`, r.basic(data))
}

func (MaintenanceAssignmentArcMachineResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
  resource_providers_to_register = ["Microsoft.HybridCompute"]
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-maint-%[1]d"
  location = "%[2]s"
}

resource "azurerm_arc_machine" "test" {
  name                = "acctest-arc-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  kind                = "SCVMM"
}

resource "azurerm_maintenance_configuration" "test" {
  name                     = "acctest-mc%[1]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  scope                    = "InGuestPatch"
  in_guest_user_patch_mode = "User"

  window {
    start_date_time = formatdate("YYYY-MM-DD hh:mm", timestamp())
    duration        = "02:00"
    time_zone       = "Greenwich Standard Time"
    recur_every     = "1Day"
  }

  install_patches {
    reboot = "Always"

    windows {
      classifications_to_include = ["Critical"]
    }
  }

  lifecycle {
    ignore_changes = [
      window[0].start_date_time
    ]
  }
}
`, data.RandomInteger, data.Locations.Primary)
}
