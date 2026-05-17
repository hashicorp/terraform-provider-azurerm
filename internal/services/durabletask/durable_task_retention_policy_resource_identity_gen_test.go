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

func TestAccDurableTaskRetentionPolicy_resourceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_durable_task_retention_policy", "test")
	r := RetentionPolicyResource{}

	checkedFields := map[string]struct{}{
		"subscription_id":     {},
		"resource_group_name": {},
		"scheduler_name":      {},
	}

	data.ResourceIdentityTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			ConfigStateChecks: []statecheck.StateCheck{
				customstatecheck.ExpectAllIdentityFieldsAreChecked("azurerm_durable_task_retention_policy.test", checkedFields),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_durable_task_retention_policy.test", tfjsonpath.New("resource_group_name"), tfjsonpath.New("scheduler_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_durable_task_retention_policy.test", tfjsonpath.New("scheduler_name"), tfjsonpath.New("scheduler_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_durable_task_retention_policy.test", tfjsonpath.New("subscription_id"), tfjsonpath.New("scheduler_id")),
			},
		},
		data.ImportBlockWithResourceIdentityStep(false),
		data.ImportBlockWithIDStep(false),
	}, false)
}
