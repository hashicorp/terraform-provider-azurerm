// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package maintenance_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	customstatecheck "github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/statecheck"
)

func TestAccMaintenanceAssignmentArcMachine_resourceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_maintenance_assignment_arc_machine", "test")
	r := MaintenanceAssignmentArcMachineResource{}

	checkedFields := map[string]struct{}{
		"scope": {},
		"name":  {},
	}

	data.ResourceIdentityTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			ConfigStateChecks: []statecheck.StateCheck{
				customstatecheck.ExpectAllIdentityFieldsAreChecked("azurerm_maintenance_assignment_arc_machine.test", checkedFields),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_maintenance_assignment_arc_machine.test", tfjsonpath.New("scope"), tfjsonpath.New("arc_machine_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPathCaseInsensitive("azurerm_maintenance_assignment_arc_machine.test", tfjsonpath.New("name"), tfjsonpath.New("maintenance_configuration_id")),
			},
		},
		data.ImportBlockWithResourceIdentityStep(false),
		data.ImportBlockWithIDStep(false),
	}, false)
}
