package consumption_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type BudgetSubscriptionDataSource struct{}

func TestAccBudgetSubscriptionDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_consumption_budget_subscription", "test")
	r := BudgetSubscriptionDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.template(data),
		},
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("subscription_id").HasValue(data.Client().SubscriptionID),
				check.That(data.ResourceName).Key("name").HasValue(fmt.Sprintf("acctestconsumptionbudget-%d", data.RandomInteger)),
			),
		},
	})
}

func (d BudgetSubscriptionDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_subscription" "current" {}

data "azurerm_consumption_budget_subscription" "test" {
  name            = "acctestconsumptionbudget-%d"
  subscription_id = data.azurerm_subscription.current.subscription_id
}
`, data.RandomInteger)
}

func (BudgetSubscriptionDataSource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_subscription" "current" {}

resource "azurerm_consumption_budget_subscription" "test" {
  name            = "acctestconsumptionbudget-%d"
  subscription_id = data.azurerm_subscription.current.subscription_id

  amount     = 1000
  time_grain = "Monthly"

  time_period {
    start_date = "%s"
  }

  filter {
    tag {
      name = "foo"
      values = [
        "bar",
      ]
    }
  }

  notification {
    enabled   = true
    threshold = 90.0
    operator  = "EqualTo"

    contact_emails = [
      "foo@example.com",
      "bar@example.com",
    ]
  }
}
`, data.RandomInteger, consumptionBudgetTestStartDate().Format(time.RFC3339))
}
