// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package durabletask_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	customstatecheck "github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/statecheck"
)

func TestAccDurableTaskHub_resourceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_durable_task_hub", "test")
	r := TaskHubResource{}

	checkedFields := map[string]struct{}{
		"subscription_id":     {},
		"name":                {},
		"resource_group_name": {},
		"scheduler_name":      {},
	}

	data.ResourceIdentityTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			ConfigStateChecks: []statecheck.StateCheck{
				customstatecheck.ExpectAllIdentityFieldsAreChecked("azurerm_durable_task_hub.test", checkedFields),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_durable_task_hub.test", tfjsonpath.New("name"), tfjsonpath.New("name")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_durable_task_hub.test", tfjsonpath.New("resource_group_name"), tfjsonpath.New("scheduler_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_durable_task_hub.test", tfjsonpath.New("scheduler_name"), tfjsonpath.New("scheduler_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_durable_task_hub.test", tfjsonpath.New("subscription_id"), tfjsonpath.New("scheduler_id")),
			},
		},
		data.ImportBlockWithResourceIdentityStep(false),
		data.ImportBlockWithIDStep(false),
	}, false)
}
