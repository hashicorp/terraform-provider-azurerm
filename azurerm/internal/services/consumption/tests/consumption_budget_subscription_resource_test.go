package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMConsumptionBudgetSubscription_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_consumption_budget_subscription", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMConsumptionBudgetSubscriptionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMConsumptionBudgetSubscription_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMConsumptionBudgetSubscriptionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMConsumptionBudgetSubscription_basicUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_consumption_budget_subscription", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMConsumptionBudgetSubscriptionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMConsumptionBudgetSubscription_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMConsumptionBudgetSubscriptionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMConsumptionBudgetSubscription_basicUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMConsumptionBudgetSubscriptionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMConsumptionBudgetSubscription_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_consumption_budget_subscription", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMConsumptionBudgetSubscriptionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMConsumptionBudgetSubscription_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMConsumptionBudgetSubscriptionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}
func TestAccAzureRMConsumptionBudgetSubscription_completeUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_consumption_budget_subscription", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMConsumptionBudgetSubscriptionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMConsumptionBudgetSubscription_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMConsumptionBudgetSubscriptionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMConsumptionBudgetSubscription_completeUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMConsumptionBudgetSubscriptionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMConsumptionBudgetSubscription_usageCategory(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_consumption_budget_subscription", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMConsumptionBudgetSubscriptionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMConsumptionBudgetSubscription_usageCategory(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMConsumptionBudgetSubscriptionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMConsumptionBudgetSubscriptionExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acceptance.AzureProvider.Meta().(*clients.Client).Consumption.BudgetsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		consumptionBudgetName := rs.Primary.Attributes["name"]
		subscriptionId, hasSubscriptionId := rs.Primary.Attributes["subscription_id"]
		scope := fmt.Sprintf("/subscriptions/%s", subscriptionId)
		if !hasSubscriptionId {
			return fmt.Errorf("bad: no subscription id found in state for Consumption Budget: %s", consumptionBudgetName)
		}

		resp, err := conn.Get(ctx, scope, consumptionBudgetName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: Consumption Budget %q for scope %q does not exist", consumptionBudgetName, scope)
			}

			return fmt.Errorf("bad: Get on consumptionBudgetClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMConsumptionBudgetSubscriptionDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).Consumption.BudgetsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_consumption_budget_subscription" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		subscriptionId := rs.Primary.Attributes["subscription_id"]
		scope := fmt.Sprintf("/subscriptions/%s", subscriptionId)

		resp, err := conn.Get(ctx, scope, name)

		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return err
			}
		}
	}

	return nil
}

func testAccAzureRMConsumptionBudgetSubscription_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_subscription" "current" {}

resource "azurerm_consumption_budget_subscription" "test" {
  name            = "acctestconsumptionbudgetsubscription-%d"
  subscription_id = data.azurerm_subscription.current.subscription_id

  amount     = 1000
  category   = "Cost"
  time_grain = "Monthly"

  time_period {
    start_date = "2020-11-01T00:00:00Z"
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
`, data.RandomInteger)
}

func testAccAzureRMConsumptionBudgetSubscription_basicUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_subscription" "current" {}

resource "azurerm_consumption_budget_subscription" "test" {
  name            = "acctestconsumptionbudgetsubscription-%d"
  subscription_id = data.azurerm_subscription.current.subscription_id

  // Changed the amount from 1000 to 2000
  amount     = 3000
  category   = "Cost"
  time_grain = "Monthly"

  // Add end_date
  time_period {
    start_date = "2020-11-01T00:00:00Z"
    end_date   = "2020-12-01T00:00:00Z"
  }

  // Changed threshold and operator
  notification {
    enabled   = true
    threshold = 95.0
    operator  = "GreaterThan"

    contact_emails = [
      "foo@example.com",
      "bar@example.com",
    ]
  }
}
`, data.RandomInteger)
}

func testAccAzureRMConsumptionBudgetSubscription_complete(data acceptance.TestData) string {
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
  subscription_id = data.azurerm_subscription.current.subscription_id

  amount     = 1000
  category   = "Cost"
  time_grain = "Monthly"

  time_period {
    start_date = "2020-11-01T00:00:00Z"
    end_date   = "2020-12-01T00:00:00Z"
  }

  filter {
    resource_groups = [
      azurerm_resource_group.test.name,
    ]
    resources = [
      azurerm_monitor_action_group.test.id,
    ]
    meters = [
      "00000000-0000-0000-0000-000000000000",
    ]
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
    enabled   = false
    threshold = 100.0
    operator  = "GreaterThan"

    contact_emails = [
      "foo@example.com",
      "bar@example.com",
    ]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMConsumptionBudgetSubscription_completeUpdate(data acceptance.TestData) string {
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
  subscription_id = data.azurerm_subscription.current.subscription_id

  // Changed the amount from 1000 to 2000
  amount     = 2000
  category   = "Cost"
  time_grain = "Monthly"

  // Removed end_date
  time_period {
    start_date = "2020-11-01T00:00:00Z"
  }

  filter {
    resource_groups = [
      azurerm_resource_group.test.name,
    ]
    // Removed resources
    meters = [
      "00000000-0000-0000-0000-000000000000",
    ]
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
  }

  notification {
    enabled   = true
    threshold = 90.0
    operator  = "EqualTo"

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
    enabled   = true
    threshold = 100.0
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMConsumptionBudgetSubscription_usageCategory(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_subscription" "current" {}

resource "azurerm_consumption_budget_subscription" "test" {
  name            = "acctestconsumptionbudgetsubscription-%d"
  subscription_id = data.azurerm_subscription.current.subscription_id

  amount     = 1000
  category   = "Usage"
  time_grain = "Monthly"

  time_period {
    start_date = "2020-11-01T00:00:00Z"
  }

  filter {
    meters = [
      "00000000-0000-0000-0000-000000000000",
    ]
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
`, data.RandomInteger)
}
