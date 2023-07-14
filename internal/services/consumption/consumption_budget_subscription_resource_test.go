// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consumption_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/consumption/2019-10-01/budgets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func consumptionBudgetTestStartDate() time.Time {
	utcNow := time.Now().UTC()
	startDate := time.Date(utcNow.Year(), utcNow.Month(), 1, 0, 0, 0, 0, utcNow.Location())

	return startDate
}

type ConsumptionBudgetSubscriptionResource struct{}

func TestAccConsumptionBudgetSubscription_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_consumption_budget_subscription", "test")
	r := ConsumptionBudgetSubscriptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccConsumptionBudgetSubscription_basicUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_consumption_budget_subscription", "test")
	r := ConsumptionBudgetSubscriptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccConsumptionBudgetSubscription_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_consumption_budget_subscription", "test")
	r := ConsumptionBudgetSubscriptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_consumption_budget_subscription"),
		},
	})
}

func TestAccConsumptionBudgetSubscription_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_consumption_budget_subscription", "test")
	r := ConsumptionBudgetSubscriptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccConsumptionBudgetSubscription_completeUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_consumption_budget_subscription", "test")
	r := ConsumptionBudgetSubscriptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.completeUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (ConsumptionBudgetSubscriptionResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := budgets.ParseScopedBudgetID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Consumption.BudgetsClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (ConsumptionBudgetSubscriptionResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_subscription" "test" {}

resource "azurerm_consumption_budget_subscription" "test" {
  name            = "acctestconsumptionbudgetsubscription-%d"
  subscription_id = data.azurerm_subscription.test.id

  amount     = 1000
  time_grain = "Monthly"

  time_period {
    start_date = "%s"
  }

  filter {
    tag {
      name = "foo"
      values = [
        "bar"
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

func (ConsumptionBudgetSubscriptionResource) basicUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_subscription" "current" {}

resource "azurerm_consumption_budget_subscription" "test" {
  name            = "acctestconsumptionbudgetsubscription-%d"
  subscription_id = data.azurerm_subscription.current.id

  // Changed the amount from 1000 to 2000
  amount     = 3000
  time_grain = "Monthly"

  // Add end_date
  time_period {
    start_date = "%s"
    end_date   = "%s"
  }

  // Remove filter

  // Changed threshold and operator
  notification {
    enabled        = true
    threshold      = 95.0
    threshold_type = "Forecasted"
    operator       = "GreaterThan"

    contact_emails = [
      "foo@example.com",
      "bar@example.com",
    ]
  }
}
`, data.RandomInteger, consumptionBudgetTestStartDate().Format(time.RFC3339), consumptionBudgetTestStartDate().AddDate(1, 1, 0).Format(time.RFC3339))
}

func (ConsumptionBudgetSubscriptionResource) requiresImport(data acceptance.TestData) string {
	template := ConsumptionBudgetSubscriptionResource{}.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_consumption_budget_subscription" "import" {
  name            = azurerm_consumption_budget_subscription.test.name
  subscription_id = azurerm_consumption_budget_subscription.test.subscription_id

  amount     = azurerm_consumption_budget_subscription.test.amount
  time_grain = azurerm_consumption_budget_subscription.test.time_grain

  time_period {
    start_date = "%s"
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
`, template, consumptionBudgetTestStartDate().Format(time.RFC3339))
}

func (ConsumptionBudgetSubscriptionResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_subscription" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_monitor_action_group" "test" {
  name                = "acctestAG-%d"
  resource_group_name = azurerm_resource_group.test.name
  short_name          = "acctestAG"
}

resource "azurerm_consumption_budget_subscription" "test" {
  name            = "acctestconsumptionbudgetsubscription-%d"
  subscription_id = data.azurerm_subscription.current.id

  amount     = 1000
  time_grain = "Monthly"

  time_period {
    start_date = "%s"
    end_date   = "%s"
  }

  filter {
    dimension {
      name = "ResourceGroupName"
      values = [
        azurerm_resource_group.test.name,
      ]
    }

    dimension {
      name = "ResourceId"
      values = [
        azurerm_monitor_action_group.test.id,
      ]
    }

    tag {
      name = "foo"
      values = [
        "bar",
        "baz",
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

    contact_groups = [
      azurerm_monitor_action_group.test.id,
    ]

    contact_roles = [
      "Owner",
    ]
  }

  notification {
    enabled        = false
    threshold      = 100.0
    operator       = "GreaterThan"
    threshold_type = "Forecasted"

    contact_emails = [
      "foo@example.com",
      "bar@example.com",
    ]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, consumptionBudgetTestStartDate().Format(time.RFC3339), consumptionBudgetTestStartDate().AddDate(1, 1, 0).Format(time.RFC3339))
}

func (ConsumptionBudgetSubscriptionResource) completeUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_subscription" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_monitor_action_group" "test" {
  name                = "acctestAG-%d"
  resource_group_name = azurerm_resource_group.test.name
  short_name          = "acctestAG"
}

resource "azurerm_consumption_budget_subscription" "test" {
  name            = "acctestconsumptionbudgetsubscription-%d"
  subscription_id = data.azurerm_subscription.current.id

  // Changed the amount from 1000 to 2000
  amount     = 2000
  time_grain = "Monthly"

  // Removed end_date
  time_period {
    start_date = "%s"
  }

  filter {
    dimension {
      name = "ResourceGroupName"
      values = [
        azurerm_resource_group.test.name,
      ]
    }

    tag {
      name = "foo"
      values = [
        "bar",
        "baz",
      ]
    }

    // Added tag: zip
    tag {
      name = "zip"
      values = [
        "zap",
        "zop",
      ]
    }

    // Removed not block 
  }

  notification {
    enabled        = true
    threshold      = 90.0
    operator       = "EqualTo"
    threshold_type = "Actual"

    contact_emails = [
      // Added baz@example.com
      "baz@example.com",
      "foo@example.com",
      "bar@example.com",
    ]

    contact_groups = [
      azurerm_monitor_action_group.test.id,
    ]
    // Removed contact_roles
  }

  notification {
    // Set enabled to true
    enabled        = true
    threshold      = 100.0
    threshold_type = "Forecasted"
    // Changed from EqualTo to GreaterThanOrEqualTo 
    operator = "GreaterThanOrEqualTo"

    contact_emails = [
      "foo@example.com",
      "bar@example.com",
    ]

    // Added contact_groups
    contact_groups = [
      azurerm_monitor_action_group.test.id,
    ]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, consumptionBudgetTestStartDate().Format(time.RFC3339))
}
