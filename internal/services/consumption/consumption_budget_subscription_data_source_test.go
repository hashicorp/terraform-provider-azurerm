package consumption_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type ConsumptionBudgetSubscriptionDataSource struct{}

func TestAccBudgetSubscriptionDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_consumption_budget_subscription", "test")
	r := ConsumptionBudgetSubscriptionDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("subscription_id").HasValue(data.Client().SubscriptionID),
				check.That(data.ResourceName).Key("name").Exists(),
			),
		},
	})
}

func (d ConsumptionBudgetSubscriptionDataSource) basic(data acceptance.TestData) string {
	config := ConsumptionBudgetResourceGroupResource{}.basic(data)
	return fmt.Sprintf(`
  %s

data "azurerm_consumption_budget_subscription" "test" {
  name            = azurerm_subscription.test.name
  subscription_id = azurerm_subscription.test.subscription_id
}
`, config)
}
