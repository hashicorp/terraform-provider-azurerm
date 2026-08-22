// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package policy_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/guestconfiguration/2024-04-05/guestconfigurationhcrpassignments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type PolicyArcMachineConfigurationAssignmentResource struct{}

func TestAccPolicyArcMachineConfigurationAssignment_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_policy_arc_machine_configuration_assignment", "test")
	r := PolicyArcMachineConfigurationAssignmentResource{}

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

func TestAccPolicyArcMachineConfigurationAssignment_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_policy_arc_machine_configuration_assignment", "test")
	r := PolicyArcMachineConfigurationAssignmentResource{}

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

func TestAccPolicyArcMachineConfigurationAssignment_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_policy_arc_machine_configuration_assignment", "test")
	r := PolicyArcMachineConfigurationAssignmentResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateGuestConfiguration(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r PolicyArcMachineConfigurationAssignmentResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := guestconfigurationhcrpassignments.ParseProviders2GuestConfigurationAssignmentID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.Policy.GuestConfigurationHCRPAssignmentsClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return pointer.To(resp.Model != nil), nil
}

func (r PolicyArcMachineConfigurationAssignmentResource) template(data acceptance.TestData) string {
	rgid := data.RandomInteger
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestrg-arcgc-%d"
  location = "%s"
}

resource "azurerm_arc_machine" "test" {
  name                = "acctestrg-arcgc-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  kind                = "SCVMM"

  identity {
    type = "SystemAssigned"
  }

  tags = {
    foo = "bar"
  }
}
`, rgid, data.Locations.Primary, rgid)
}

func (r PolicyArcMachineConfigurationAssignmentResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_policy_arc_machine_configuration_assignment" "test" {
  name     = "WhitelistedApplication"
  location = azurerm_arc_machine.test.location

  machine_id = azurerm_arc_machine.test.id

  configuration {
    version = "1.*"

    parameter {
      name  = "[InstalledApplication]bwhitelistedapp;Name"
      value = "NotePad,sql"
    }
  }
}
`, r.template(data))
}

func (r PolicyArcMachineConfigurationAssignmentResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_policy_arc_machine_configuration_assignment" "import" {
  name       = azurerm_policy_arc_machine_configuration_assignment.test.name
  location   = azurerm_policy_arc_machine_configuration_assignment.test.location
  machine_id = azurerm_policy_arc_machine_configuration_assignment.test.machine_id

  configuration {
    version = "1.*"

    parameter {
      name  = "[InstalledApplication]bwhitelistedapp;Name"
      value = "NotePad,sql"
    }
  }
}
`, r.basic(data))
}

func (r PolicyArcMachineConfigurationAssignmentResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_policy_arc_machine_configuration_assignment" "test" {
  name       = "WhitelistedApplication"
  location   = azurerm_arc_machine.test.location
  machine_id = azurerm_arc_machine.test.id

  configuration {
    version         = "1.1.1.1"
    assignment_type = "ApplyAndAutoCorrect"
    content_hash    = upper("db4d5cd43c59c756f9beb1f029c858bc341587bf75332288270a26493565f058")
    content_uri     = "https://testcontenturi/package"

    parameter {
      name  = "[InstalledApplication]bwhitelistedapp;Name"
      value = "NotePad,sql"
    }
  }
}
`, r.template(data))
}

func (r PolicyArcMachineConfigurationAssignmentResource) updateGuestConfiguration(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_policy_arc_machine_configuration_assignment" "test" {
  name       = "WhitelistedApplication"
  location   = azurerm_arc_machine.test.location
  machine_id = azurerm_arc_machine.test.id

  configuration {
    version         = "1.1.1.1"
    assignment_type = "Audit"
    content_hash    = upper("cde01f651f3a3055834753d42d73b44e2a505844ac34f9ccc35d3d6dfffcb2e4")
    content_uri     = "https://testcontenturi/package2"

    parameter {
      name  = "[InstalledApplication]bwhitelistedapp;Name"
      value = "NotePad,sql"
    }
  }
}
`, r.template(data))
}
