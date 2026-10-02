// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package eventhub_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	customstatecheck "github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/statecheck"
)

func TestAccEventhubConsumerGroup_resourceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub_consumer_group", "test")
	r := EventhubConsumerGroupResource{}

	checkedFields := map[string]struct{}{
		"subscription_id":     {},
		"eventhub_name":       {},
		"name":                {},
		"namespace_name":      {},
		"resource_group_name": {},
	}

	data.ResourceIdentityTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			ConfigStateChecks: []statecheck.StateCheck{
				customstatecheck.ExpectAllIdentityFieldsAreChecked("azurerm_eventhub_consumer_group.test", checkedFields),
				statecheck.ExpectIdentityValue("azurerm_eventhub_consumer_group.test", tfjsonpath.New("subscription_id"), knownvalue.StringExact(data.Subscriptions.Primary)),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_eventhub_consumer_group.test", tfjsonpath.New("eventhub_name"), tfjsonpath.New("eventhub_name")),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_eventhub_consumer_group.test", tfjsonpath.New("name"), tfjsonpath.New("name")),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_eventhub_consumer_group.test", tfjsonpath.New("namespace_name"), tfjsonpath.New("namespace_name")),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_eventhub_consumer_group.test", tfjsonpath.New("resource_group_name"), tfjsonpath.New("resource_group_name")),
			},
		},
		data.ImportBlockWithResourceIdentityStep(false),
		data.ImportBlockWithIDStep(false),
	}, false)
}
