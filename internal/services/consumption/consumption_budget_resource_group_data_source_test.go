// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consumption_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type ConsumptionBudgetResourceGroupDataSource struct{}

func TestAccDataSourceConsumptionBudgetResourceGroup_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_consumption_budget_resource_group", "test")
	r := ConsumptionBudgetResourceGroupDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("resource_group_id").Exists(),
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

func (d ConsumptionBudgetResourceGroupDataSource) basic(data acceptance.TestData) string {
	config := ConsumptionBudgetResourceGroupResource{}.basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_consumption_budget_resource_group" "test" {
  name              = azurerm_consumption_budget_resource_group.test.name
  resource_group_id = azurerm_consumption_budget_resource_group.test.resource_group_id
}
`, config)
}
