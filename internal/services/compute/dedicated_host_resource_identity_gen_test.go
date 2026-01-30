// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package compute_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	customstatecheck "github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/statecheck"
)

func TestAccDedicatedHost_resourceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dedicated_host", "test")
	r := DedicatedHostResource{}

	checkedFields := map[string]struct{}{
		"name":                {},
		"host_group_name":     {},
		"resource_group_name": {},
		"subscription_id":     {},
	}

	data.ResourceIdentityTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			ConfigStateChecks: []statecheck.StateCheck{
				customstatecheck.ExpectAllIdentityFieldsAreChecked("azurerm_dedicated_host.test", checkedFields),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_dedicated_host.test", tfjsonpath.New("name"), tfjsonpath.New("name")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_dedicated_host.test", tfjsonpath.New("host_group_name"), tfjsonpath.New("dedicated_host_group_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_dedicated_host.test", tfjsonpath.New("resource_group_name"), tfjsonpath.New("dedicated_host_group_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_dedicated_host.test", tfjsonpath.New("subscription_id"), tfjsonpath.New("dedicated_host_group_id")),
			},
		},
		data.ImportBlockWithResourceIdentityStep(false),
		data.ImportBlockWithIDStep(false),
	}, false)
}
