// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package subscription_policy_exemption_test

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
)

type SubscriptionAssignmentTestResource struct{}

func (r SubscriptionAssignmentTestResource) template() string {
	return `data "azurerm_subscription" "test" {}`
}

func (r SubscriptionAssignmentTestResource) withBuiltInPolicySetBasic(data acceptance.TestData) string {
	template := r.template()
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

data "azurerm_policy_set_definition" "test" {
  display_name = "Audit machines with insecure password security settings"
}

resource "azurerm_subscription_policy_assignment" "test" {
  name                 = "acctestpa-sub-%[2]d"
  subscription_id      = data.azurerm_subscription.test.id
  policy_definition_id = data.azurerm_policy_set_definition.test.id
  location             = %[3]q

  identity {
    type = "SystemAssigned"
  }
}
`, template, data.RandomInteger, data.Locations.Primary)
}
