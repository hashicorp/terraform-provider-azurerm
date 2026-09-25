// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package servicebus_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	customstatecheck "github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/statecheck"
)

func TestAccServicebusQueue_resourceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_queue", "test")
	r := ServicebusQueueResource{}

	checkedFields := map[string]struct{}{
		"subscription_id":     {},
		"name":                {},
		"namespace_name":      {},
		"resource_group_name": {},
	}

	data.ResourceIdentityTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			ConfigStateChecks: []statecheck.StateCheck{
				customstatecheck.ExpectAllIdentityFieldsAreChecked("azurerm_servicebus_queue.test", checkedFields),
				statecheck.ExpectIdentityValue("azurerm_servicebus_queue.test", tfjsonpath.New("subscription_id"), knownvalue.StringExact(data.Subscriptions.Primary)),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_servicebus_queue.test", tfjsonpath.New("name"), tfjsonpath.New("name")),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_servicebus_queue.test", tfjsonpath.New("namespace_name"), tfjsonpath.New("namespace_name")),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_servicebus_queue.test", tfjsonpath.New("resource_group_name"), tfjsonpath.New("resource_group_name")),
			},
		},
		data.ImportBlockWithResourceIdentityStep(false),
		data.ImportBlockWithIDStep(false),
	}, false)
}
