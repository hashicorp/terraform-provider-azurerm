// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package mssql_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	customstatecheck "github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/statecheck"
)

func TestAccMssqlVirtualMachine_resourceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_virtual_machine", "test")
	r := MssqlVirtualMachineResource{}

	checkedFields := map[string]struct{}{
		"name":                {},
		"resource_group_name": {},
		"subscription_id":     {},
	}

	data.ResourceIdentityTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			ConfigStateChecks: []statecheck.StateCheck{
				customstatecheck.ExpectAllIdentityFieldsAreChecked("azurerm_mssql_virtual_machine.test", checkedFields),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_mssql_virtual_machine.test", tfjsonpath.New("name"), tfjsonpath.New("virtual_machine_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_mssql_virtual_machine.test", tfjsonpath.New("resource_group_name"), tfjsonpath.New("virtual_machine_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_mssql_virtual_machine.test", tfjsonpath.New("subscription_id"), tfjsonpath.New("virtual_machine_id")),
			},
		},
		data.ImportBlockWithResourceIdentityStep(false),
		data.ImportBlockWithIDStep(false),
	}, false)
}
