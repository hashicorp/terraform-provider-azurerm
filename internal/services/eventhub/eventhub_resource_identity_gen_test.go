// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package eventhub_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	customstatecheck "github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/statecheck"
)

func TestAccEventhub_resourceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub", "test")
	r := EventhubResource{}

	checkedFields := map[string]struct{}{
		"name":                {},
		"namespace_name":      {},
		"resource_group_name": {},
		"subscription_id":     {},
	}

	data.ResourceIdentityTest(t, []acceptance.TestStep{
		{
			Config: r.standard(data),
			ConfigStateChecks: []statecheck.StateCheck{
				customstatecheck.ExpectAllIdentityFieldsAreChecked("azurerm_eventhub.test", checkedFields),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_eventhub.test", tfjsonpath.New("name"), tfjsonpath.New("name")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_eventhub.test", tfjsonpath.New("namespace_name"), tfjsonpath.New("namespace_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_eventhub.test", tfjsonpath.New("resource_group_name"), tfjsonpath.New("namespace_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_eventhub.test", tfjsonpath.New("subscription_id"), tfjsonpath.New("namespace_id")),
			},
		},
		data.ImportBlockWithResourceIdentityStep(false),
		data.ImportBlockWithIDStep(false),
	}, false)
}
