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
		"name":                {},
		"machine_name":        {},
		"resource_group_name": {},
		"subscription_id":     {},
	}

	data.ResourceIdentityTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			ConfigStateChecks: []statecheck.StateCheck{
				customstatecheck.ExpectAllIdentityFieldsAreChecked("azurerm_maintenance_assignment_arc_machine.test", checkedFields),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_maintenance_assignment_arc_machine.test", tfjsonpath.New("name"), tfjsonpath.New("name")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_maintenance_assignment_arc_machine.test", tfjsonpath.New("machine_name"), tfjsonpath.New("arc_machine_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_maintenance_assignment_arc_machine.test", tfjsonpath.New("resource_group_name"), tfjsonpath.New("arc_machine_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_maintenance_assignment_arc_machine.test", tfjsonpath.New("subscription_id"), tfjsonpath.New("arc_machine_id")),
			},
		},
		data.ImportBlockWithResourceIdentityStep(false),
		data.ImportBlockWithIDStep(false),
	}, false)
}
