// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package automation_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	customstatecheck "github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/statecheck"
)

func TestAccAutomationRuntimeEnvironmentPackage_resourceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_runtime_environment_package", "test")
	r := AutomationRuntimeEnvironmentPackageResource{}

	checkedFields := map[string]struct{}{
		"name":                     {},
		"automation_account_name":  {},
		"resource_group_name":      {},
		"runtime_environment_name": {},
		"subscription_id":          {},
	}

	data.ResourceIdentityTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			ConfigStateChecks: []statecheck.StateCheck{
				customstatecheck.ExpectAllIdentityFieldsAreChecked("azurerm_automation_runtime_environment_package.test", checkedFields),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_automation_runtime_environment_package.test", tfjsonpath.New("name"), tfjsonpath.New("name")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_automation_runtime_environment_package.test", tfjsonpath.New("automation_account_name"), tfjsonpath.New("automation_runtime_environment_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_automation_runtime_environment_package.test", tfjsonpath.New("resource_group_name"), tfjsonpath.New("automation_runtime_environment_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_automation_runtime_environment_package.test", tfjsonpath.New("runtime_environment_name"), tfjsonpath.New("automation_runtime_environment_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_automation_runtime_environment_package.test", tfjsonpath.New("subscription_id"), tfjsonpath.New("automation_runtime_environment_id")),
			},
		},
		data.ImportBlockWithResourceIdentityStep(true),
		data.ImportBlockWithIDStep(true),
	}, false)
}
