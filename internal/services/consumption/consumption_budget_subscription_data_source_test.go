// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consumption_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type ConsumptionBudgetSubscriptionDataSource struct{}

func TestAccDataSourceConsumptionBudgetSubscription_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_consumption_budget_subscription", "test")
	r := ConsumptionBudgetSubscriptionDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("subscription_id").HasValue("/subscriptions/"+data.Client().SubscriptionID),
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("amount").HasValue("1000"),
				check.That(data.ResourceName).Key("time_grain").HasValue("Monthly"),
				check.That(data.ResourceName).Key("time_period.#").Exists(),
				check.That(data.ResourceName).Key("time_period.0.start_date").Exists(),
				check.That(data.ResourceName).Key("filter.#").Exists(),
				check.That(data.ResourceName).Key("filter.0.tag.0.name").HasValue("foo"),
				check.That(data.ResourceName).Key("filter.0.tag.0.values.#").Exists(),
				check.That(data.ResourceName).Key("notification.#").Exists(),
				check.That(data.ResourceName).Key("notification.0.threshold").HasValue("90"),
				check.That(data.ResourceName).Key("notification.0.operator").HasValue("EqualTo"),
				check.That(data.ResourceName).Key("notification.0.enabled").Exists(),
				check.That(data.ResourceName).Key("notification.0.contact_emails.0").HasValue("foo@example.com"),
				check.That(data.ResourceName).Key("notification.0.contact_emails.1").HasValue("bar@example.com"),
			),
		},
	})
}

func (d ConsumptionBudgetSubscriptionDataSource) basic(data acceptance.TestData) string {
	config := ConsumptionBudgetSubscriptionResource{}.basic(data)
	return fmt.Sprintf(`
  %s

data "azurerm_consumption_budget_subscription" "test" {
  name            = azurerm_consumption_budget_subscription.test.name
  subscription_id = azurerm_consumption_budget_subscription.test.subscription_id
}
`, config)
}
