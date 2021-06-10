package consumption_test

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type BudgetDataSource struct{}

func TestAccDataSourceBudget_current(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_budget", "current")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: BudgetDataSource{}.currentConfig(),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("subscription_id").HasValue(data.Client().SubscriptionID),
				check.That(data.ResourceName).Key("tenant_id").Exists(),
			),
		},
	})
}

func (d BudgetDataSource) currentConfig() string {
	return `
provider "azurerm" {
  features {}
}

data "azurerm_budget" "current" {}
`
}
