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

func TestAccServicebusSubscriptionRule_resourceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_subscription_rule", "test")
	r := ServicebusSubscriptionRuleResource{}

	checkedFields := map[string]struct{}{
		"subscription_id":     {},
		"name":                {},
		"namespace_name":      {},
		"resource_group_name": {},
		"subscription_name":   {},
		"topic_name":          {},
	}

	data.ResourceIdentityTest(t, []acceptance.TestStep{
		{
			Config: r.basicSqlFilter(data),
			ConfigStateChecks: []statecheck.StateCheck{
				customstatecheck.ExpectAllIdentityFieldsAreChecked("azurerm_servicebus_subscription_rule.test", checkedFields),
				statecheck.ExpectIdentityValue("azurerm_servicebus_subscription_rule.test", tfjsonpath.New("subscription_id"), knownvalue.StringExact(data.Subscriptions.Primary)),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_servicebus_subscription_rule.test", tfjsonpath.New("name"), tfjsonpath.New("name")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_servicebus_subscription_rule.test", tfjsonpath.New("namespace_name"), tfjsonpath.New("subscription_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_servicebus_subscription_rule.test", tfjsonpath.New("resource_group_name"), tfjsonpath.New("subscription_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_servicebus_subscription_rule.test", tfjsonpath.New("subscription_name"), tfjsonpath.New("subscription_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_servicebus_subscription_rule.test", tfjsonpath.New("topic_name"), tfjsonpath.New("subscription_id")),
			},
		},
		data.ImportBlockWithResourceIdentityStep(false),
		data.ImportBlockWithIDStep(false),
	}, false)
}
